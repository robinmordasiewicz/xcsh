/**
 * API Client Types
 * Type definitions for F5 XC API client
 */

/**
 * Retry configuration for transient errors
 */
export interface RetryConfig {
	/** Maximum number of retry attempts (default: 3) */
	maxRetries?: number;
	/** Initial delay in milliseconds (default: 1000) */
	initialDelayMs?: number;
	/** Maximum delay in milliseconds (default: 10000) */
	maxDelayMs?: number;
	/** Multiplier for exponential backoff (default: 2) */
	backoffMultiplier?: number;
	/** Add jitter to prevent thundering herd (default: true) */
	jitter?: boolean;
}

/**
 * API client configuration
 */
export interface APIClientConfig {
	/** F5 XC API server URL */
	serverUrl: string;
	/** API token for authentication */
	apiToken?: string;
	/** Request timeout in milliseconds */
	timeout?: number;
	/** Enable debug logging */
	debug?: boolean;
	/** Retry configuration for transient errors */
	retry?: RetryConfig;
}

/**
 * HTTP methods supported by the API
 */
export type HTTPMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH";

/**
 * API request options
 */
export interface APIRequestOptions {
	/** HTTP method */
	method: HTTPMethod;
	/** API path (e.g., /api/web/namespaces) */
	path: string;
	/** Request body (will be JSON serialized) */
	body?: Record<string, unknown>;
	/** Additional headers */
	headers?: Record<string, string>;
	/** Query parameters */
	query?: Record<string, string>;
}

/**
 * API response
 */
export interface APIResponse<T = unknown> {
	/** HTTP status code */
	statusCode: number;
	/** Response body (parsed JSON) */
	data: T;
	/** Response headers */
	headers: Record<string, string>;
	/** Whether the request was successful (2xx status) */
	ok: boolean;
}

/**
 * API error response from F5 XC
 */
export interface APIErrorResponse {
	/** Error message */
	message?: string;
	/** Error code */
	code?: string;
	/** Additional details */
	details?: string;
}

/**
 * API error with additional context
 */
export class APIError extends Error {
	constructor(
		message: string,
		public readonly statusCode: number,
		public readonly response?: APIErrorResponse,
		public readonly operation?: string,
	) {
		super(message);
		this.name = "APIError";
	}

	/**
	 * Get helpful hint based on status code
	 */
	getHint(): string {
		switch (this.statusCode) {
			case 401:
				return "Authentication failed. Check your API token with 'login profile show'";
			case 403:
				return "Permission denied. You may not have access to this resource.";
			case 404:
				return "Resource not found. Verify the name and namespace are correct.";
			case 409:
				return "Conflict - resource may already exist or be in a conflicting state.";
			case 429:
				return "Rate limited. Please wait and try again.";
			case 500:
			case 502:
			case 503:
				return "Server error. Please try again later or contact support.";
			default:
				return "";
		}
	}
}

/**
 * List response from F5 XC API
 */
export interface ListResponse<T = unknown> {
	items: T[];
	metadata?: {
		total_count?: number;
		page_size?: number;
		next_page_token?: string;
	};
}

/**
 * Resource metadata from F5 XC API
 */
export interface ResourceMetadata {
	name: string;
	namespace?: string;
	uid?: string;
	labels?: Record<string, string>;
	annotations?: Record<string, string>;
	creation_timestamp?: string;
	modification_timestamp?: string;
}

/**
 * Generic resource with metadata
 */
export interface Resource<T = unknown> {
	metadata: ResourceMetadata;
	spec?: T;
	status?: unknown;
}
