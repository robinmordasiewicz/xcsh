/**
 * F5 XC API Client
 * HTTP client for F5 Distributed Cloud API with authentication support
 */

import type {
	APIClientConfig,
	APIRequestOptions,
	APIResponse,
	APIErrorResponse,
	HTTPMethod,
	RetryConfig,
} from "./types.js";
import { APIError } from "./types.js";

/**
 * Default retry configuration
 */
const DEFAULT_RETRY_CONFIG: Required<RetryConfig> = {
	maxRetries: 3,
	initialDelayMs: 1000,
	maxDelayMs: 10000,
	backoffMultiplier: 2,
	jitter: true,
};

/**
 * Status codes that should trigger a retry
 */
const RETRYABLE_STATUS_CODES = new Set([
	408, // Request Timeout
	429, // Too Many Requests
	500, // Internal Server Error
	502, // Bad Gateway
	503, // Service Unavailable
	504, // Gateway Timeout
]);

/**
 * Calculate delay with exponential backoff and optional jitter
 */
function calculateBackoffDelay(
	attempt: number,
	config: Required<RetryConfig>,
): number {
	const exponentialDelay =
		config.initialDelayMs * Math.pow(config.backoffMultiplier, attempt);
	const cappedDelay = Math.min(exponentialDelay, config.maxDelayMs);

	if (config.jitter) {
		// Add random jitter between 0-25% of the delay
		const jitterFactor = 1 + Math.random() * 0.25;
		return Math.floor(cappedDelay * jitterFactor);
	}

	return cappedDelay;
}

/**
 * Sleep for specified milliseconds
 */
function sleep(ms: number): Promise<void> {
	return new Promise((resolve) => setTimeout(resolve, ms));
}

/**
 * API Client for F5 Distributed Cloud
 */
export class APIClient {
	private readonly serverUrl: string;
	private readonly apiToken: string;
	private readonly timeout: number;
	private readonly debug: boolean;
	private readonly retryConfig: Required<RetryConfig>;

	// Token validation state
	private _isValidated: boolean = false;
	private _validationError: string | null = null;

	constructor(config: APIClientConfig) {
		// Normalize server URL (remove trailing slash)
		this.serverUrl = config.serverUrl.replace(/\/+$/, "");
		this.apiToken = config.apiToken ?? "";
		this.timeout = config.timeout ?? 30000;
		this.debug = config.debug ?? false;
		this.retryConfig = {
			...DEFAULT_RETRY_CONFIG,
			...config.retry,
		};
	}

	/**
	 * Check if client has authentication configured (token exists).
	 * Note: This does NOT verify the token is valid. Use isValidated() for that.
	 */
	isAuthenticated(): boolean {
		return this.apiToken !== "";
	}

	/**
	 * Validate the API token by making a lightweight API call.
	 * Returns true if token is valid, false otherwise.
	 */
	async validateToken(): Promise<{ valid: boolean; error?: string }> {
		// Skip if no token configured
		if (!this.apiToken) {
			this._isValidated = false;
			this._validationError = "No API token configured";
			return { valid: false, error: this._validationError };
		}

		if (this.debug) {
			console.error(`DEBUG: Validating token against ${this.serverUrl}`);
		}

		try {
			// Use namespaces endpoint for token validation (lightweight, universal)
			await this.get<{ items?: unknown[] }>("/api/web/namespaces");
			this._isValidated = true;
			this._validationError = null;
			return { valid: true };
		} catch (error) {
			if (error instanceof APIError) {
				if (error.statusCode === 401) {
					// Definitive: token is invalid or expired
					this._isValidated = false;
					this._validationError = "Invalid or expired API token";
					return { valid: false, error: this._validationError };
				} else if (error.statusCode === 403) {
					// Definitive: token lacks permissions
					this._isValidated = false;
					this._validationError = "Token lacks required permissions";
					return { valid: false, error: this._validationError };
				} else {
					// Non-auth error (404, 500, network) - assume token is OK
					// Don't show warning for server-side issues
					if (this.debug) {
						console.error(
							`DEBUG: Validation endpoint returned ${error.statusCode}, assuming token is valid`,
						);
					}
					this._isValidated = true;
					this._validationError = null;
					return { valid: true };
				}
			} else {
				// Unknown error - assume token is OK, don't block user
				if (this.debug) {
					console.error(
						`DEBUG: Validation error: ${error instanceof Error ? error.message : "Unknown"}, assuming token is valid`,
					);
				}
				this._isValidated = true;
				this._validationError = null;
				return { valid: true };
			}
		}
	}

	/**
	 * Check if client has a validated (verified working) token
	 */
	isValidated(): boolean {
		return this._isValidated;
	}

	/**
	 * Get the validation error message, if any
	 */
	getValidationError(): string | null {
		return this._validationError;
	}

	/**
	 * Clear validation state (called on profile switch)
	 */
	clearValidationCache(): void {
		this._isValidated = false;
		this._validationError = null;
	}

	/**
	 * Get the server URL
	 */
	getServerUrl(): string {
		return this.serverUrl;
	}

	/**
	 * Build full URL from path and query parameters
	 */
	private buildUrl(path: string, query?: Record<string, string>): string {
		let baseUrl = this.serverUrl;

		// Handle case where base URL ends with /api and path starts with /api
		if (baseUrl.endsWith("/api") && path.startsWith("/api")) {
			baseUrl = baseUrl.slice(0, -4);
		}

		// Ensure path starts with /
		const normalizedPath = path.startsWith("/") ? path : `/${path}`;
		let url = `${baseUrl}${normalizedPath}`;

		// Add query parameters
		if (query && Object.keys(query).length > 0) {
			const params = new URLSearchParams(query);
			url = `${url}?${params.toString()}`;
		}

		return url;
	}

	/**
	 * Check if an error is retryable
	 */
	private isRetryableError(error: unknown): boolean {
		if (error instanceof APIError) {
			return RETRYABLE_STATUS_CODES.has(error.statusCode);
		}
		// Network errors and timeouts are retryable
		if (error instanceof Error) {
			return (
				error.name === "AbortError" ||
				error.message.includes("fetch failed") ||
				error.message.includes("network") ||
				error.message.includes("ECONNREFUSED") ||
				error.message.includes("ENOTFOUND") ||
				error.message.includes("ETIMEDOUT")
			);
		}
		return false;
	}

	/**
	 * Execute a single HTTP request attempt
	 */
	private async executeRequest<T = unknown>(
		options: APIRequestOptions,
		url: string,
		headers: Record<string, string>,
		body: string | null,
	): Promise<APIResponse<T>> {
		const controller = new AbortController();
		const timeoutId = setTimeout(() => controller.abort(), this.timeout);

		try {
			const response = await fetch(url, {
				method: options.method,
				headers,
				body,
				signal: controller.signal,
			});

			clearTimeout(timeoutId);

			// Read response body
			const responseText = await response.text();
			let data: T;

			try {
				data = responseText ? JSON.parse(responseText) : ({} as T);
			} catch {
				// If not valid JSON, wrap as string
				data = responseText as unknown as T;
			}

			// Debug response
			if (this.debug) {
				console.error(`DEBUG: Response status: ${response.status}`);
			}

			// Convert headers to record
			const responseHeaders: Record<string, string> = {};
			response.headers.forEach((value, key) => {
				responseHeaders[key] = value;
			});

			const result: APIResponse<T> = {
				statusCode: response.status,
				data,
				headers: responseHeaders,
				ok: response.ok,
			};

			// Throw error for non-2xx responses
			if (!response.ok) {
				const errorResponse = data as unknown as APIErrorResponse;
				throw new APIError(
					errorResponse.message ?? `HTTP ${response.status}`,
					response.status,
					errorResponse,
					`${options.method} ${options.path}`,
				);
			}

			return result;
		} catch (error) {
			clearTimeout(timeoutId);

			// Re-throw APIError as-is
			if (error instanceof APIError) {
				throw error;
			}

			// Handle abort/timeout
			if (error instanceof Error && error.name === "AbortError") {
				throw new APIError(
					`Request timed out after ${this.timeout}ms`,
					408,
					undefined,
					`${options.method} ${options.path}`,
				);
			}

			// Handle network errors
			if (error instanceof Error) {
				throw new APIError(
					`Network error: ${error.message}`,
					0,
					undefined,
					`${options.method} ${options.path}`,
				);
			}

			throw error;
		}
	}

	/**
	 * Execute an HTTP request with retry logic
	 */
	async request<T = unknown>(
		options: APIRequestOptions,
	): Promise<APIResponse<T>> {
		const url = this.buildUrl(options.path, options.query);

		// Prepare headers
		const headers: Record<string, string> = {
			"Content-Type": "application/json",
			Accept: "application/json",
			...options.headers,
		};

		// Add API token authorization
		if (this.apiToken) {
			headers["Authorization"] = `APIToken ${this.apiToken}`;
		}

		// Prepare body
		const body: string | null = options.body
			? JSON.stringify(options.body)
			: null;

		// Debug logging
		if (this.debug) {
			console.error(`DEBUG: ${options.method} ${url}`);
			if (body) {
				console.error(`DEBUG: Request body: ${body}`);
			}
		}

		let lastError: Error | undefined;

		// Retry loop
		for (
			let attempt = 0;
			attempt <= this.retryConfig.maxRetries;
			attempt++
		) {
			try {
				return await this.executeRequest<T>(
					options,
					url,
					headers,
					body,
				);
			} catch (error) {
				lastError =
					error instanceof Error ? error : new Error(String(error));

				// Check if we should retry
				const isRetryable = this.isRetryableError(error);
				const hasRetriesLeft = attempt < this.retryConfig.maxRetries;

				if (isRetryable && hasRetriesLeft) {
					const delay = calculateBackoffDelay(
						attempt,
						this.retryConfig,
					);

					if (this.debug) {
						const statusInfo =
							error instanceof APIError
								? ` (${error.statusCode})`
								: "";
						console.error(
							`DEBUG: Request failed${statusInfo}, retrying in ${delay}ms (attempt ${attempt + 1}/${this.retryConfig.maxRetries})`,
						);
					}

					await sleep(delay);
					continue;
				}

				// Not retryable or no retries left - throw the error
				throw error;
			}
		}

		// Should not reach here, but just in case
		throw lastError ?? new Error("Request failed after all retries");
	}

	/**
	 * GET request
	 */
	async get<T = unknown>(
		path: string,
		query?: Record<string, string>,
	): Promise<APIResponse<T>> {
		const options: APIRequestOptions = {
			method: "GET",
			path,
		};
		if (query) {
			options.query = query;
		}
		return this.request<T>(options);
	}

	/**
	 * POST request
	 */
	async post<T = unknown>(
		path: string,
		body?: Record<string, unknown>,
	): Promise<APIResponse<T>> {
		const options: APIRequestOptions = {
			method: "POST",
			path,
		};
		if (body) {
			options.body = body;
		}
		return this.request<T>(options);
	}

	/**
	 * PUT request
	 */
	async put<T = unknown>(
		path: string,
		body?: Record<string, unknown>,
	): Promise<APIResponse<T>> {
		const options: APIRequestOptions = {
			method: "PUT",
			path,
		};
		if (body) {
			options.body = body;
		}
		return this.request<T>(options);
	}

	/**
	 * DELETE request
	 */
	async delete<T = unknown>(path: string): Promise<APIResponse<T>> {
		return this.request<T>({
			method: "DELETE",
			path,
		});
	}

	/**
	 * PATCH request
	 */
	async patch<T = unknown>(
		path: string,
		body?: Record<string, unknown>,
	): Promise<APIResponse<T>> {
		const options: APIRequestOptions = {
			method: "PATCH",
			path,
		};
		if (body) {
			options.body = body;
		}
		return this.request<T>(options);
	}
}

/**
 * Create an API client from environment variables
 */
export function createClientFromEnv(
	envPrefix: string = "F5XC",
): APIClient | null {
	const serverUrl = process.env[`${envPrefix}_API_URL`];
	const apiToken = process.env[`${envPrefix}_API_TOKEN`] ?? "";

	if (!serverUrl) {
		return null;
	}

	const config: APIClientConfig = {
		serverUrl,
		debug: process.env[`${envPrefix}_DEBUG`] === "true",
	};

	if (apiToken) {
		config.apiToken = apiToken;
	}

	return new APIClient(config);
}

/**
 * Build API path for a domain resource
 */
export function buildResourcePath(
	domain: string,
	resource: string,
	action: string,
	namespace?: string,
	name?: string,
): string {
	// Standard F5 XC API path pattern:
	// /api/web/namespaces/{namespace}/{resource}
	// /api/web/namespaces/{namespace}/{resource}/{name}
	// /api/config/namespaces/{namespace}/{resource}
	// etc.

	let path = `/api/${domain}`;

	if (namespace) {
		path += `/namespaces/${namespace}`;
	}

	path += `/${resource}`;

	if (name) {
		path += `/${name}`;
	}

	// Handle specific actions that modify the path
	if (action && action !== "list" && action !== "get") {
		// Some actions like "create" are POST to base path
		// Others like "delete" or specific actions may need different handling
	}

	return path;
}

// Re-export types
export type { APIClientConfig, APIRequestOptions, APIResponse, HTTPMethod };
export { APIError };
