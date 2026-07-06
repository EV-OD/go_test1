package user

import apperrors "myapp/internal/errors"

var (
	ErrInvalidCredentials  = apperrors.New(apperrors.CodeUnauthorized, "Invalid email or password")
	ErrUserNotFound        = apperrors.New(apperrors.CodeNotFound, "User not found")
	ErrInsufficientBalance = apperrors.New(apperrors.CodeInvalidRequest, "insufficient balance")
	ErrContextType         = apperrors.New(apperrors.CodeInternalServer, "Context type error")
)
