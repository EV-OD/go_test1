package errors

import "github.com/labstack/echo/v5"

type Code string

const (
	CodeNotFound       Code = "NOT_FOUND"
	CodeUnauthorized   Code = "UNAUTHORIZED"
	CodeInvalidRequest Code = "INVALID_REQUEST"
	CodeInternalServer Code = "INTERNAL_SERVER"
	CodeInvalidToken   Code = "INVALID_TOKEN"
	CodeInvalidClaims  Code = "INVALID_CLAIMS"
	CodeConfigError    Code = "CONFIG_ERROR"
)

type AppError struct {
	Code    Code
	Message string
	Err     error
}

type ApiError struct {
	Error string `json:"error" example:"error message"`
}

type ApiSuccess struct {
	Message string `json:"message" example:"success message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Response(c *echo.Context) error {
	return c.JSON(HTTPStatus(e.Code), map[string]string{"error": e.Message})
}

func ErrorResponse(c *echo.Context, err error) error {
	if appErr, ok := err.(*AppError); ok {
		return c.JSON(HTTPStatus(appErr.Code), map[string]string{"error": appErr.Message})
	}
	return c.JSON(500, map[string]string{"error": err.Error()})
}

func SuccessResponse(c *echo.Context, message string) error {
	return c.JSON(200, map[string]string{"message": message})
}

func New(code Code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func NewWithErr(code Code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func HTTPStatus(code Code) int {
	switch code {
	case CodeNotFound:
		return 404
	case CodeUnauthorized, CodeInvalidToken, CodeInvalidClaims:
		return 401
	case CodeInvalidRequest:
		return 400
	case CodeConfigError:
		return 500
	default:
		return 500
	}
}

var (
	ErrUnauthorized   = New(CodeUnauthorized, "Unauthorized")
	ErrInvalidRequest = New(CodeInvalidRequest, "Invalid request")
	ErrInvalidToken   = New(CodeInvalidToken, "Invalid token")
	ErrInvalidClaims  = New(CodeInvalidClaims, "Invalid claims")
	ErrNotFound       = New(CodeNotFound, "Not found")
)
