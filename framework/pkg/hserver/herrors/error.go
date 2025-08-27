package herrors

import (
	"errors"
	"fmt"
	"net/http"
)

// 错误处理最佳实践：
// 1. 当方法返回 error 接口时，如果内部返回的是 nil Herr，应该直接返回 nil
// 2. 判断错误是否为 nil 时：
//    - 如果变量类型是 Herr，直接使用 err == nil 判断
//    - 如果变量类型是 error，使用 IsNilError(err) 判断
// 3. 避免直接将 nil Herr 赋值给 error 接口变量

const (
	ReasonStatusInternalHError = "STATUS_INTERNAL_SERVER_ERROR"
	ReasonParameterError       = "PARAMETER_ERROR"
	ReqParameterError          = "ReqParameterError"
)

// HError 服务端错误定义
type HError struct {
	Code          int
	DefMessage    string
	Reason        string
	BusinessError error
}

// Herr 错误类型别名,允许nil值
// 注意：由于 Herr 是 *HError 的别名，当返回 nil 时，使用 err != nil 判断可能会得到错误的结果
// 建议使用 IsNil() 方法来判断是否为 nil
type Herr = *HError

func New(code int, reason string, msg string) Herr {
	return &HError{Code: code, Reason: reason, DefMessage: msg}
}

// DefaultError 默认错误
func DefaultError() Herr {
	return New(http.StatusInternalServerError, ReasonStatusInternalHError, "Server Internal Error")
}

// NewErr 根据error创建错误
func NewErr(err error) Herr {
	if err == nil {
		return nil
	}
	return &HError{Code: http.StatusInternalServerError, Reason: ReasonStatusInternalHError, DefMessage: err.Error(), BusinessError: err}
}

func (r *HError) WithCode(code int) Herr {
	r.Code = code
	return r
}

func (r *HError) WithDefMsg(msg string) Herr {
	r.DefMessage = msg
	return r
}

func (r *HError) WithReason(reason string) Herr {
	r.Reason = reason
	return r
}

func (r *HError) WithBusinessError(err error) Herr {
	r.BusinessError = err
	return r
}

func (r *HError) Error() string {
	return fmt.Sprintf("code:%d,reason:%s,message:%s", r.Code, r.Reason, r.DefMessage)
}

// IsNil 判断是否为 nil
func (r *HError) IsNil() bool {
	return r == nil
}

func IsHError(err error) bool {
	var e *HError
	return errors.As(err, &e)
}

func TohError(err error) Herr {
	if err == nil {
		return nil
	}
	var e *HError
	if errors.As(err, &e) {
		return e
	}
	return NewErr(err)
}

func HaveError(err error) bool {
	return TohError(err) != nil
}

// IsNilError 判断错误是否为 nil，正确处理 error 接口和 Herr 类型的转换
func IsNilError(err error) bool {
	if err == nil {
		return true
	}
	if herr, ok := err.(Herr); ok {
		return herr == nil
	}
	return false
}

func Is(s Herr, t Herr) bool {
	return s.Code == t.Code && s.Reason == t.Reason
}
