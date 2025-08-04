package response

import (
	"strconv"
	"strings"

	"github.com/zhoudm1743/unicom-gw/api"
)

// CommonJsonResponse 通用JSON响应
type CommonJsonResponse struct {
	api.BaseIoTGatewayResponse
	Data map[string]interface{} `json:"data"`
}

// GetData 获取响应数据
func (r *CommonJsonResponse) GetData() map[string]interface{} {
	return r.Data
}

// SetData 设置响应数据
func (r *CommonJsonResponse) SetData(data map[string]interface{}) {
	r.Data = data
}

// IsSuccess 判断响应是否成功
func (r *CommonJsonResponse) IsSuccess() bool {
	// 优先使用基类的成功标志
	if r.IsSucceed != nil {
		return *r.IsSucceed
	}

	// 检查数据是否为空
	if r.Data == nil {
		return false
	}

	// 检查是否包含错误码
	if _, exists := r.Data[api.ERROR_CODE]; exists {
		isSuccess := false
		r.IsSucceed = &isSuccess
		return false
	}

	// 检查响应码
	for key, value := range r.Data {
		lowerKey := strings.ToLower(key)
		if lowerKey == "respcode" || lowerKey == "rspcode" {
			// 处理值类型
			var codeStr string
			switch v := value.(type) {
			case string:
				codeStr = v
			case int:
				codeStr = strconv.Itoa(v)
			case float64:
				codeStr = strconv.Itoa(int(v))
			default:
				continue
			}

			// 检查响应码是否为0
			if codeStr == "0" {
				isSuccess := true
				r.IsSucceed = &isSuccess
				return true
			} else {
				isSuccess := false
				r.IsSucceed = &isSuccess
				return false
			}
		}
	}

	// 默认成功
	isSuccess := true
	r.IsSucceed = &isSuccess
	return true
}
