package httphelper

import "net/http"

var (
	ErrInternalError           = NewAPIError(http.StatusInternalServerError, "Internal error")
	ErrInvalidRequestError     = NewAPIError(http.StatusBadRequest, "Invalid request")
	ErrAuthenticationError     = NewAPIError(http.StatusUnauthorized, "Authentication error")
	ErrRateLimitError          = NewAPIError(http.StatusTooManyRequests, "Rate limit exceeded")
	ErrInternalServerError     = NewAPIError(http.StatusInternalServerError, "Internal server error")
	ErrServiceUnavailableError = NewAPIError(http.StatusServiceUnavailable, "Service unavailable")
	ErrForbiddenError          = NewAPIError(http.StatusForbidden, "Forbidden")
	ErrNotFoundError           = NewAPIError(http.StatusNotFound, "Not found")
	ErrConflictError           = NewAPIError(http.StatusConflict, "Conflict")
	ErrValidationError         = NewAPIError(http.StatusBadRequest, "Validation error")
	ErrJSONMarshalError        = NewAPIError(http.StatusInternalServerError, "Json marshal error")
	ErrHTTPRequestError        = NewAPIError(http.StatusInternalServerError, "Http request error")
	ErrTooManyRequestsError    = NewAPIError(http.StatusTooManyRequests, "Too many requests")
)

// Error is a custom error type for API errors.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func NewAPIError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
