package response

import (
	"strings"

	"github.com/zhoudm1743/unicom-gw/api"
)

// CommonJsonResponse 通用JSON响应实现
type CommonJsonResponse struct {
	api.BaseIoTGatewayResponse
	Data map[string]interface{} `json:"data"`
}

// NewCommonJsonResponse 创建一个新的通用JSON响应
func NewCommonJsonResponse() *CommonJsonResponse {
	return &CommonJsonResponse{}
}

// GetData 获取响应数据
func (r *CommonJsonResponse) GetData() map[string]interface{} {
	return r.Data
}

// SetData 设置响应数据
func (r *CommonJsonResponse) SetData(data map[string]interface{}) {
	r.Data = data
}

// IsSuccess 请求是否成功
func (r *CommonJsonResponse) IsSuccess() bool {
	// 如果已设置成功状态，直接返回
	if r.BaseIoTGatewayResponse.IsSuccess() {
		return r.BaseIoTGatewayResponse.IsSuccess()
	}

	// 如果数据为空，则请求失败
	if r.Data == nil {
		return false
	}

	// 检查是否包含错误码
	if _, ok := r.Data[api.ERROR_CODE]; ok {
		r.SetSuccess(false)
		return false
	}

	// 检查响应码
	for key, value := range r.Data {
		// 查找响应码字段
		lowerKey := strings.ToLower(key)
		if lowerKey == "respcode" || lowerKey == "rspcode" {
			// 如果响应码为0，则请求成功
			if strValue, ok := value.(string); ok {
				if strValue == "0" {
					r.SetSuccess(true)
					return true
				}
				r.SetSuccess(false)
				return false
			}
		}
	}

	// 默认认为请求成功
	r.SetSuccess(true)
	return true
}
