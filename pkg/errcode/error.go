package errcode

import (
	"errors"
	"fmt"
)

// Error 是贯穿各层的业务错误类型。
// service 层返回它，middleware/response 层据此生成响应。
type Error struct {
	Code    Code
	Message string
	cause   error
}

// New 用错误码创建业务错误，使用错误码默认文案。
func New(code Code) *Error {
	return &Error{Code: code, Message: code.Message()}
}

// Newf 用自定义文案创建业务错误。
func Newf(code Code, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Wrap 把底层错误包装为业务错误，底层错误只用于日志，不透出给调用方。
func Wrap(code Code, cause error) *Error {
	return &Error{Code: code, Message: code.Message(), cause: cause}
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持 errors.Is/As 向下追溯。
func (e *Error) Unwrap() error { return e.cause }

// Is 按错误码判等，便于 errors.Is(err, errcode.New(errcode.UserNotFound))。
func (e *Error) Is(target error) bool {
	var t *Error
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
}

// WithCause 返回携带底层错误的副本，不修改原对象。
func (e *Error) WithCause(cause error) *Error {
	cloned := *e
	cloned.cause = cause
	return &cloned
}

// WithMessage 返回替换文案的副本，不修改原对象。
func (e *Error) WithMessage(format string, args ...any) *Error {
	cloned := *e
	cloned.Message = fmt.Sprintf(format, args...)
	return &cloned
}

// HTTPStatus 返回该错误对应的 HTTP 状态码。
func (e *Error) HTTPStatus() int { return e.Code.HTTPStatus() }

// From 从任意错误中提取业务错误，非业务错误统一归为 ServerError。
func From(err error) *Error {
	if err == nil {
		return nil
	}
	var be *Error
	if errors.As(err, &be) {
		return be
	}
	return Wrap(ServerError, err)
}

// IsCode 判断错误是否为指定业务错误码。
func IsCode(err error, code Code) bool {
	var be *Error
	if errors.As(err, &be) {
		return be.Code == code
	}
	return false
}
