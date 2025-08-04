package request

import (
	"encoding/json"

	"github.com/zhoudm1743/unicom-gw/api"
	"github.com/zhoudm1743/unicom-gw/api/response"
)

// CommonJsonRequest 通用JSON请求
type CommonJsonRequest struct {
	api.BaseIoTGatewayRequest
	Params map[string]interface{}
}

// NewCommonJsonRequest 创建通用JSON请求
func NewCommonJsonRequest() *CommonJsonRequest {
	req := &CommonJsonRequest{}
	req.ApiType = api.API_TYPE_JSON
	return req
}

// GetParams 获取参数
func (r *CommonJsonRequest) GetParams() map[string]interface{} {
	return r.Params
}

// SetParams 设置参数
func (r *CommonJsonRequest) SetParams(params map[string]interface{}) {
	r.Params = params
}

// GetResponseClass 获取响应类
func (r *CommonJsonRequest) GetResponseClass() interface{} {
	return &response.CommonJsonResponse{}
}

// Check 参数检查
func (r *CommonJsonRequest) Check() error {
	// TODO: 实现参数检查逻辑
	return nil
}

// ExecProcessBeforeReqSend 请求发送前处理
func (r *CommonJsonRequest) ExecProcessBeforeReqSend(params []interface{}) {
	if len(params) > 0 {
		if mapParams, ok := params[0].(map[string]interface{}); ok {
			if transID, exists := mapParams["trans_id"]; exists {
				if transIDStr, ok := transID.(string); ok {
					r.SetTransId(transIDStr)
				}
			}

			// 将参数转换为JSON
			jsonData, _ := json.Marshal(mapParams)
			r.SetReqText(string(jsonData))
		}
	}
}
