package apperror

import "errors"

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
	ErrDatabase   = errors.New("database error")
)

type Error struct {
	Kind    error
	Message string
	Details map[string]any
	Cause   error
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Kind != nil {
		return e.Kind.Error()
	}
	return "application error"
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(kind error, message string) *Error {
	return &Error{Kind: kind, Message: message}
}

func WithDetails(kind error, message string, details map[string]any) *Error {
	return &Error{Kind: kind, Message: message, Details: details}
}

func Wrap(kind error, message string, cause error) *Error {
	return &Error{Kind: kind, Message: message, Cause: cause}
}

func Is(err error, kind error) bool {
	var appErr *Error
	if errors.As(err, &appErr) {
		return errors.Is(appErr.Kind, kind)
	}
	return errors.Is(err, kind)
}
