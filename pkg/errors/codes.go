// Package errors provides standardized exit codes and error types for f5xcctl.
// These codes enable AI agents and automation scripts to programmatically
// determine the nature of errors and take appropriate action.
package errors

// Exit codes for f5xcctl CLI
// These codes follow a hierarchical scheme where higher codes indicate more specific errors.
const (
	// ExitSuccess indicates the command completed successfully
	ExitSuccess = 0

	// ExitGenericError indicates an unclassified error occurred
	ExitGenericError = 1

	// ExitValidationError indicates invalid arguments or validation failure
	// Examples: missing required flags, invalid flag values, malformed input
	ExitValidationError = 2

	// ExitAuthError indicates authentication or authorization failure
	// Examples: invalid credentials, expired token, permission denied
	ExitAuthError = 3

	// ExitConnectionError indicates connection or timeout to API
	// Examples: network unreachable, DNS failure, API timeout, server error (5xx)
	ExitConnectionError = 4

	// ExitNotFoundError indicates a requested resource was not found
	// Examples: namespace doesn't exist, object not found (404)
	ExitNotFoundError = 5

	// ExitConflictError indicates a resource conflict
	// Examples: resource already exists, version conflict (409)
	ExitConflictError = 6

	// ExitRateLimitError indicates rate limiting was encountered
	// Examples: too many requests (429)
	ExitRateLimitError = 7
)

// Error code strings for machine-readable error messages
const (
	// Validation errors
	ErrMissingFlag  = "ERR_MISSING_FLAG"
	ErrInvalidValue = "ERR_INVALID_VALUE"
	ErrInvalidInput = "ERR_INVALID_INPUT"
	ErrMissingInput = "ERR_MISSING_INPUT"

	// Authentication errors
	ErrAuthFailed   = "ERR_AUTH_FAILED"
	ErrForbidden    = "ERR_FORBIDDEN"
	ErrTokenExpired = "ERR_TOKEN_EXPIRED"
	ErrCredsMissing = "ERR_CREDS_MISSING"

	// Connection errors
	ErrConnectionFailed = "ERR_CONNECTION_FAILED"
	ErrTimeout          = "ERR_TIMEOUT"
	ErrServerError      = "ERR_SERVER_ERROR"

	// Resource errors
	ErrNotFound  = "ERR_NOT_FOUND"
	ErrConflict  = "ERR_CONFLICT"
	ErrRateLimit = "ERR_RATE_LIMIT"

	// Operation errors
	ErrOperationFailed    = "ERR_OPERATION_FAILED"
	ErrOperationCancelled = "ERR_OPERATION_CANCELLED"
)

// HTTPStatusToExitCode maps HTTP status codes to exit codes
func HTTPStatusToExitCode(statusCode int) int {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return ExitSuccess
	case statusCode == 400:
		return ExitValidationError
	case statusCode == 401:
		return ExitAuthError
	case statusCode == 403:
		return ExitAuthError
	case statusCode == 404:
		return ExitNotFoundError
	case statusCode == 409:
		return ExitConflictError
	case statusCode == 429:
		return ExitRateLimitError
	case statusCode >= 500:
		return ExitConnectionError
	default:
		return ExitGenericError
	}
}

// HTTPStatusToErrorCode maps HTTP status codes to error code strings
func HTTPStatusToErrorCode(statusCode int) string {
	switch statusCode {
	case 400:
		return ErrInvalidInput
	case 401:
		return ErrAuthFailed
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 409:
		return ErrConflict
	case 429:
		return ErrRateLimit
	case 500, 502, 503, 504:
		return ErrServerError
	default:
		return ErrOperationFailed
	}
}

// ExitCodeDescription returns a human-readable description for an exit code
func ExitCodeDescription(code int) string {
	switch code {
	case ExitSuccess:
		return "Success"
	case ExitGenericError:
		return "Generic/unknown error"
	case ExitValidationError:
		return "Invalid arguments or validation failure"
	case ExitAuthError:
		return "Authentication or authorization failure"
	case ExitConnectionError:
		return "Connection or timeout to API"
	case ExitNotFoundError:
		return "Resource not found"
	case ExitConflictError:
		return "Resource conflict"
	case ExitRateLimitError:
		return "Rate limited"
	default:
		return "Unknown exit code"
	}
}
