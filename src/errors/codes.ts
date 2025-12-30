/**
 * Error Codes and Exit Codes
 *
 * Standardized exit codes and error types for xcsh.
 * These codes enable AI agents and automation scripts to programmatically
 * determine the nature of errors and take appropriate action.
 */

/**
 * Exit codes for xcsh CLI
 * These codes follow a hierarchical scheme where higher codes indicate more specific errors.
 */
export const ExitCode = {
	/** Command completed successfully */
	Success: 0,
	/** Unclassified error occurred */
	GenericError: 1,
	/** Invalid arguments or validation failure */
	ValidationError: 2,
	/** Authentication or authorization failure */
	AuthError: 3,
	/** Connection or timeout to API */
	ConnectionError: 4,
	/** Requested resource was not found */
	NotFoundError: 5,
	/** Resource conflict */
	ConflictError: 6,
	/** Rate limiting was encountered */
	RateLimitError: 7,
	/** Subscription quota would be exceeded */
	QuotaExceeded: 8,
	/** Required feature is not available */
	FeatureNotAvailable: 9,
} as const;

export type ExitCodeValue = (typeof ExitCode)[keyof typeof ExitCode];

/**
 * Error code strings for machine-readable error messages
 */
export const ErrorCode = {
	// Validation errors
	MissingFlag: "ERR_MISSING_FLAG",
	InvalidValue: "ERR_INVALID_VALUE",
	InvalidInput: "ERR_INVALID_INPUT",
	MissingInput: "ERR_MISSING_INPUT",

	// Authentication errors
	AuthFailed: "ERR_AUTH_FAILED",
	Forbidden: "ERR_FORBIDDEN",
	TokenExpired: "ERR_TOKEN_EXPIRED",
	CredentialsMissing: "ERR_CREDS_MISSING",

	// Connection errors
	ConnectionFailed: "ERR_CONNECTION_FAILED",
	Timeout: "ERR_TIMEOUT",
	ServerError: "ERR_SERVER_ERROR",

	// Resource errors
	NotFound: "ERR_NOT_FOUND",
	Conflict: "ERR_CONFLICT",
	RateLimit: "ERR_RATE_LIMIT",

	// Operation errors
	OperationFailed: "ERR_OPERATION_FAILED",
	OperationCancelled: "ERR_OPERATION_CANCELLED",

	// Subscription errors
	QuotaExceeded: "ERR_QUOTA_EXCEEDED",
	FeatureNotAvailable: "ERR_FEATURE_NOT_AVAILABLE",
	UpgradeRequired: "ERR_UPGRADE_REQUIRED",
} as const;

export type ErrorCodeValue = (typeof ErrorCode)[keyof typeof ErrorCode];

/**
 * Map HTTP status code to exit code
 */
export function httpStatusToExitCode(statusCode: number): ExitCodeValue {
	if (statusCode >= 200 && statusCode < 300) {
		return ExitCode.Success;
	}

	switch (statusCode) {
		case 400:
			return ExitCode.ValidationError;
		case 401:
		case 403:
			return ExitCode.AuthError;
		case 404:
			return ExitCode.NotFoundError;
		case 409:
			return ExitCode.ConflictError;
		case 429:
			return ExitCode.RateLimitError;
		default:
			if (statusCode >= 500) {
				return ExitCode.ConnectionError;
			}
			return ExitCode.GenericError;
	}
}

/**
 * Map HTTP status code to error code string
 */
export function httpStatusToErrorCode(statusCode: number): ErrorCodeValue {
	switch (statusCode) {
		case 400:
			return ErrorCode.InvalidInput;
		case 401:
			return ErrorCode.AuthFailed;
		case 403:
			return ErrorCode.Forbidden;
		case 404:
			return ErrorCode.NotFound;
		case 409:
			return ErrorCode.Conflict;
		case 429:
			return ErrorCode.RateLimit;
		case 500:
		case 502:
		case 503:
		case 504:
			return ErrorCode.ServerError;
		default:
			return ErrorCode.OperationFailed;
	}
}

/**
 * Get human-readable description for an exit code
 */
export function exitCodeDescription(code: ExitCodeValue): string {
	switch (code) {
		case ExitCode.Success:
			return "Success";
		case ExitCode.GenericError:
			return "Generic/unknown error";
		case ExitCode.ValidationError:
			return "Invalid arguments or validation failure";
		case ExitCode.AuthError:
			return "Authentication or authorization failure";
		case ExitCode.ConnectionError:
			return "Connection or timeout to API";
		case ExitCode.NotFoundError:
			return "Resource not found";
		case ExitCode.ConflictError:
			return "Resource conflict";
		case ExitCode.RateLimitError:
			return "Rate limited";
		case ExitCode.QuotaExceeded:
			return "Subscription quota exceeded";
		case ExitCode.FeatureNotAvailable:
			return "Feature not available in subscription";
		default:
			return "Unknown exit code";
	}
}

/**
 * Get recovery hint for an exit code
 */
export function exitCodeHint(code: ExitCodeValue): string {
	switch (code) {
		case ExitCode.Success:
			return "";
		case ExitCode.ValidationError:
			return "Check command arguments and flag values";
		case ExitCode.AuthError:
			return "Verify credentials with 'login profile show' or authenticate with 'login'";
		case ExitCode.ConnectionError:
			return "Check network connectivity and API URL configuration";
		case ExitCode.NotFoundError:
			return "Verify the resource name and namespace are correct";
		case ExitCode.ConflictError:
			return "The resource may already exist or be in a conflicting state";
		case ExitCode.RateLimitError:
			return "Wait a moment and try again";
		case ExitCode.QuotaExceeded:
			return "Check subscription quotas in the F5 XC console";
		case ExitCode.FeatureNotAvailable:
			return "This feature may require a plan upgrade";
		default:
			return "An unexpected error occurred";
	}
}

/**
 * Structured error with code for programmatic handling
 */
export interface StructuredError {
	code: ErrorCodeValue;
	exitCode: ExitCodeValue;
	message: string;
	hint?: string;
	details?: Record<string, unknown>;
}

/**
 * Create a structured error from an HTTP response
 */
export function createStructuredError(
	statusCode: number,
	message: string,
	details?: Record<string, unknown>,
): StructuredError {
	const exitCode = httpStatusToExitCode(statusCode);
	const hint = exitCodeHint(exitCode);

	const error: StructuredError = {
		code: httpStatusToErrorCode(statusCode),
		exitCode,
		message,
	};

	if (hint) {
		error.hint = hint;
	}

	if (details) {
		error.details = details;
	}

	return error;
}

/**
 * Format error for display
 */
export function formatError(error: StructuredError): string[] {
	const lines: string[] = [];
	lines.push(`Error [${error.code}]: ${error.message}`);
	if (error.hint) {
		lines.push(`Hint: ${error.hint}`);
	}
	return lines;
}
