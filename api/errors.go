package api

import "fmt"

// ApiError API调用错误
type ApiError struct {
	Code    string
	Message string
	Cause   error
}

// Error 实现error接口
func (e *ApiError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("API错误 [%s]: %s, 原因: %s", e.Code, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("API错误 [%s]: %s", e.Code, e.Message)
}

// NewApiError 创建一个新的API错误
func NewApiError(message string, cause error) *ApiError {
	return &ApiError{
		Message: message,
		Cause:   cause,
	}
}

// ApiRuleError API规则错误
type ApiRuleError struct {
	Field   string
	Message string
}

// Error 实现error接口
func (e *ApiRuleError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("API规则错误: 字段 [%s], %s", e.Field, e.Message)
	}
	return fmt.Sprintf("API规则错误: %s", e.Message)
}

// NewApiRuleError 创建一个新的API规则错误
func NewApiRuleError(field, message string) *ApiRuleError {
	return &ApiRuleError{
		Field:   field,
		Message: message,
	}
}
