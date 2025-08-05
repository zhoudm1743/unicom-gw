package request

import (
	"encoding/json"

	"github.com/zhoudm1743/unicom-gw/api"
	"github.com/zhoudm1743/unicom-gw/api/internal/utils"
	"github.com/zhoudm1743/unicom-gw/api/response"
)

// CommonJsonRequest 通用JSON请求实现
type CommonJsonRequest struct {
	api.BaseIoTGatewayRequest
	Params map[string]interface{}
}

// NewCommonJsonRequest 创建一个新的通用JSON请求
func NewCommonJsonRequest() *CommonJsonRequest {
	return &CommonJsonRequest{
		Params: make(map[string]interface{}),
	}
}

// GetParams 获取请求参数
func (r *CommonJsonRequest) GetParams() map[string]interface{} {
	return r.Params
}

// SetParams 设置请求参数
func (r *CommonJsonRequest) SetParams(params map[string]interface{}) {
	r.Params = params
}

// GetResponseClass 获取响应类型
func (r *CommonJsonRequest) GetResponseClass() api.IoTGatewayResponse {
	return &response.CommonJsonResponse{}
}

// Check 客户端参数检查，减少服务端无效调用
func (r *CommonJsonRequest) Check() error {
	// 可以在这里实现参数验证逻辑
	return nil
}

// ExecProcessBeforeReqSend 请求发送前的处理
func (r *CommonJsonRequest) ExecProcessBeforeReqSend(params []interface{}) {
	if len(params) > 0 {
		if mapParams, ok := params[0].(map[string]interface{}); ok {
			// 设置交易ID
			if transID, ok := mapParams[utils.TransIDKey].(string); ok {
				r.SetTransId(transID)
			}

			// 将参数转换为JSON
			jsonData, err := json.Marshal(mapParams)
			if err == nil {
				r.SetReqText(string(jsonData))
			}
		}
	}
}
