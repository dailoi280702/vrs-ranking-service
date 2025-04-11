package apperror

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	Raw     error
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(err error, code int, message string) *AppError {
	return &AppError{
		Raw:     err,
		Code:    code,
		Message: message,
	}
}

func (e *AppError) WithMessage(message string) *AppError {
	e.Message = message

	return e
}

func (e *AppError) WithMessagef(format string, a ...any) *AppError {
	e.Message = fmt.Sprintf(format, a...)

	return e
}

func (e *AppError) WithError(err error) *AppError {
	e.Raw = err

	return e
}

func ErrBadRequest() *AppError {
	return &AppError{
		Raw:     nil,
		Code:    http.StatusBadRequest,
		Message: "Bad request",
	}
}

func ErrInternal() *AppError {
	return &AppError{
		Raw:     nil,
		Code:    http.StatusInternalServerError,
		Message: "Interal server error",
	}
}

func ErrNotFound() *AppError {
	return &AppError{
		Raw:     nil,
		Code:    http.StatusNotFound,
		Message: "Not found",
	}
}

func ErrUnauthenticated() *AppError {
	return &AppError{
		Raw:     nil,
		Code:    http.StatusUnauthorized,
		Message: "Unauthenticated",
	}
}

func ErrForbidden() *AppError {
	return &AppError{
		Raw:     nil,
		Code:    http.StatusForbidden,
		Message: "Formbiden",
	}
}

func ErrConflicted() *AppError {
	return &AppError{
		Raw:     nil,
		Code:    http.StatusConflict,
		Message: "Conflicted",
	}
}

func Err() *AppError {
	return NewError(
		nil,
		http.StatusUnprocessableEntity,
		"Unprocessable",
	)
}

func ErrGRPC(err error) *AppError {
	st := status.Convert(err)
	msg := st.Message()

	appErr := new(AppError)

	switch st.Code() {
	case codes.Unknown, codes.Internal:
		appErr = ErrInternal()
	case codes.InvalidArgument:
		appErr = ErrBadRequest()
	case codes.Unauthenticated:
		appErr = ErrUnauthenticated()
	case codes.PermissionDenied:
		appErr = ErrForbidden()
	case codes.NotFound:
		appErr = ErrNotFound()
	default:
		appErr = ErrInternal()
	}

	return appErr.WithError(err).WithMessage(msg)
}
