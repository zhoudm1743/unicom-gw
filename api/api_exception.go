package api

import (
	"fmt"
)

// ApiException API异常类
type ApiException struct {
	ErrMsg  string
	ErrCode string
	Cause   error
}

// Error 实现error接口
func (e *ApiException) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("API错误 - 错误码: %s, 错误信息: %s, 原因: %v", e.ErrCode, e.ErrMsg, e.Cause)
	}
	return fmt.Sprintf("API错误 - 错误码: %s, 错误信息: %s", e.ErrCode, e.ErrMsg)
}

// NewApiException 创建一个新的API异常
func NewApiException(errMsg string, errCode string, cause error) *ApiException {
	return &ApiException{
		ErrMsg:  errMsg,
		ErrCode: errCode,
		Cause:   cause,
	}
}

// ApiRuleException API规则异常类
type ApiRuleException struct {
	ErrMsg  string
	ErrCode string
	Cause   error
}

// Error 实现error接口
func (e *ApiRuleException) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("API规则错误 - 错误码: %s, 错误信息: %s, 原因: %v", e.ErrCode, e.ErrMsg, e.Cause)
	}
	return fmt.Sprintf("API规则错误 - 错误码: %s, 错误信息: %s", e.ErrCode, e.ErrMsg)
}

// NewApiRuleException 创建一个新的API规则异常
func NewApiRuleException(errMsg string, errCode string, cause error) *ApiRuleException {
	return &ApiRuleException{
		ErrMsg:  errMsg,
		ErrCode: errCode,
		Cause:   cause,
	}
}
