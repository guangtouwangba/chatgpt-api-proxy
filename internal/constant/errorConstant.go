package constant

type ErrorCode int

const (
	// InternalError is returned when the error is unknown.
	InternalError ErrorCode = iota
	// InvalidRequestError is returned when the request is invalid.
	InvalidRequestError
	// AuthenticationError is returned when the authentication is invalid.
	AuthenticationError
	// RateLimitError is returned when the rate limit is exceeded.
	RateLimitError
	// InternalServerError is returned when the server returns an internal server error.
	InternalServerError
	// ServiceUnavailableError is returned when the server is unavailable.
	ServiceUnavailableError
	// ForbiddenError is returned when the request is forbidden.
	ForbiddenError
	// NotFoundError is returned when the resource is not found.
	NotFoundError
	// ConflictError is returned when the request conflicts with the current state of the server.
	ConflictError
	// ValidationError is returned when the request is invalid.
	ValidationError
	// JSONMarshalError is returned when the json marshal fails.
	JSONMarshalError
	// HTTPRequestError is returned when the http request fails.
	HTTPRequestError
	// TooManyRequestsError is returned when the request is too many.
	TooManyRequestsError
	// GatewayTimeoutError is returned when the gateway timeout.
	GatewayTimeoutError
	// JSONUnmarshalError is returned when the json unmarshal fails.
	JSONUnmarshalError
)

type BaseError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

var ErrorMessages = map[ErrorCode]string{
	InternalError:           "Internal error",
	InvalidRequestError:     "Invalid request",
	AuthenticationError:     "Authentication error",
	RateLimitError:          "Rate limit exceeded",
	InternalServerError:     "Internal server error",
	ServiceUnavailableError: "Service unavailable",
	ForbiddenError:          "Forbidden",
	NotFoundError:           "Not found",
	ConflictError:           "Conflict",
	ValidationError:         "Validation error",
	JSONMarshalError:        "Json marshal error",
	HTTPRequestError:        "Http request error",
	TooManyRequestsError:    "Too many requests",
	GatewayTimeoutError:     "Gateway timeout",
	JSONUnmarshalError:      "Json unmarshal error",
}

func (e *BaseError) Error() string {
	return e.Message
}

func Error(code ErrorCode) error {
	return &BaseError{
		Code:    code,
		Message: ErrorMessages[code],
	}
}
