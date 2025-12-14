package errors

import (
	"fmt"
)

// ExitError represents an error with an associated exit code.
// This allows the CLI to return specific exit codes for different error types,
// enabling automation scripts and AI agents to handle errors appropriately.
type ExitError struct {
	Code    int    // Exit code to return
	ErrCode string // Machine-readable error code (e.g., "ERR_NOT_FOUND")
	Message string // Human-readable error message
	Details string // Additional details for troubleshooting
	Hint    string // Helpful hint for resolving the error
}

// Error implements the error interface
func (e *ExitError) Error() string {
	if e.ErrCode != "" {
		return fmt.Sprintf("%s (code: %s)", e.Message, e.ErrCode)
	}
	return e.Message
}

// Unwrap returns nil as ExitError doesn't wrap other errors
func (e *ExitError) Unwrap() error {
	return nil
}

// WithDetails adds details to the error
func (e *ExitError) WithDetails(details string) *ExitError {
	e.Details = details
	return e
}

// WithHint adds a hint to the error
func (e *ExitError) WithHint(hint string) *ExitError {
	e.Hint = hint
	return e
}

// NewExitError creates a new ExitError with the specified exit code and message
func NewExitError(code int, errCode, message string) *ExitError {
	return &ExitError{
		Code:    code,
		ErrCode: errCode,
		Message: message,
	}
}

// ValidationError creates an error for validation failures
func ValidationError(message string) *ExitError {
	return &ExitError{
		Code:    ExitValidationError,
		ErrCode: ErrInvalidValue,
		Message: message,
	}
}

// MissingFlagError creates an error for missing required flags
func MissingFlagError(flagName string) *ExitError {
	return &ExitError{
		Code:    ExitValidationError,
		ErrCode: ErrMissingFlag,
		Message: fmt.Sprintf("required flag \"--%s\" not set", flagName),
		Hint:    fmt.Sprintf("Provide the --%s flag with a valid value", flagName),
	}
}

// InvalidInputError creates an error for invalid input
func InvalidInputError(message string) *ExitError {
	return &ExitError{
		Code:    ExitValidationError,
		ErrCode: ErrInvalidInput,
		Message: message,
	}
}

// AuthError creates an error for authentication failures
func AuthError(message string) *ExitError {
	return &ExitError{
		Code:    ExitAuthError,
		ErrCode: ErrAuthFailed,
		Message: message,
		Hint:    "Check your credentials with 'f5xcctl configure show' or set VES_API_TOKEN",
	}
}

// ForbiddenError creates an error for permission denied
func ForbiddenError(message string) *ExitError {
	return &ExitError{
		Code:    ExitAuthError,
		ErrCode: ErrForbidden,
		Message: message,
		Hint:    "You may not have permission to access this resource",
	}
}

// ConnectionError creates an error for connection failures
func ConnectionError(message string) *ExitError {
	return &ExitError{
		Code:    ExitConnectionError,
		ErrCode: ErrConnectionFailed,
		Message: message,
		Hint:    "Check your network connection and server URL",
	}
}

// TimeoutError creates an error for timeout failures
func TimeoutError(message string) *ExitError {
	return &ExitError{
		Code:    ExitConnectionError,
		ErrCode: ErrTimeout,
		Message: message,
		Hint:    "Try increasing the --timeout value or check network connectivity",
	}
}

// ServerError creates an error for server-side failures
func ServerError(statusCode int, message string) *ExitError {
	return &ExitError{
		Code:    ExitConnectionError,
		ErrCode: ErrServerError,
		Message: message,
		Details: fmt.Sprintf("HTTP status code: %d", statusCode),
		Hint:    "The server encountered an error. Please try again later or contact support.",
	}
}

// NotFoundError creates an error for resource not found
func NotFoundError(resourceType, name string) *ExitError {
	return &ExitError{
		Code:    ExitNotFoundError,
		ErrCode: ErrNotFound,
		Message: fmt.Sprintf("%s '%s' not found", resourceType, name),
		Hint:    "Verify the resource name and namespace are correct",
	}
}

// ConflictError creates an error for resource conflicts
func ConflictError(message string) *ExitError {
	return &ExitError{
		Code:    ExitConflictError,
		ErrCode: ErrConflict,
		Message: message,
		Hint:    "The resource may already exist or be in a conflicting state",
	}
}

// RateLimitError creates an error for rate limiting
func RateLimitError() *ExitError {
	return &ExitError{
		Code:    ExitRateLimitError,
		ErrCode: ErrRateLimit,
		Message: "rate limit exceeded",
		Hint:    "Please wait and try again later",
	}
}

// FromHTTPStatus creates an ExitError from an HTTP status code and response body
func FromHTTPStatus(statusCode int, operation, body string) *ExitError {
	exitCode := HTTPStatusToExitCode(statusCode)
	errCode := HTTPStatusToErrorCode(statusCode)

	var message string
	var hint string

	switch statusCode {
	case 400:
		message = fmt.Sprintf("%s failed: invalid request", operation)
		hint = "Check your input data and try again"
	case 401:
		message = fmt.Sprintf("%s failed: authentication required", operation)
		hint = "Check your credentials with 'f5xcctl configure show'"
	case 403:
		message = fmt.Sprintf("%s failed: permission denied", operation)
		hint = "You may not have access to this resource"
	case 404:
		message = fmt.Sprintf("%s failed: resource not found", operation)
		hint = "Verify the name and namespace are correct"
	case 409:
		message = fmt.Sprintf("%s failed: resource conflict", operation)
		hint = "The resource may already exist or be in a conflicting state"
	case 429:
		message = fmt.Sprintf("%s failed: rate limit exceeded", operation)
		hint = "Please wait and try again later"
	default:
		if statusCode >= 500 {
			message = fmt.Sprintf("%s failed: server error (%d)", operation, statusCode)
			hint = "Please try again later or contact support"
		} else {
			message = fmt.Sprintf("%s failed (HTTP %d)", operation, statusCode)
		}
	}

	return &ExitError{
		Code:    exitCode,
		ErrCode: errCode,
		Message: message,
		Details: body,
		Hint:    hint,
	}
}

// IsExitError checks if an error is an ExitError
func IsExitError(err error) bool {
	_, ok := err.(*ExitError)
	return ok
}

// GetExitCode returns the exit code for an error.
// If the error is an ExitError, returns its code.
// Otherwise returns ExitGenericError.
func GetExitCode(err error) int {
	if exitErr, ok := err.(*ExitError); ok {
		return exitErr.Code
	}
	return ExitGenericError
}
