package errors

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

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
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
