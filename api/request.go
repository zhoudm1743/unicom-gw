package api

// IoTGatewayRequest 定义API请求接口
type IoTGatewayRequest interface {
	// GetContentType 获取内容类型
	GetContentType() string

	// GetApiName 获取API名称
	GetApiName() string

	// GetApiVer 获取API版本
	GetApiVer() string

	// GetReqText 获取请求文本
	GetReqText() string

	// GetParams 获取所有的Key-Value形式的文本请求参数集合
	// Key: 请求参数名
	// Value: 请求参数值
	GetParams() map[string]interface{}

	// GetResponseClass 获取具体响应实现类的定义
	GetResponseClass() interface{}

	// Check 客户端参数检查，减少服务端无效调用
	Check() error

	// ExecProcessBeforeReqSend 请求发送前的处理
	ExecProcessBeforeReqSend(params []interface{})
}

// BaseIoTGatewayRequest 实现IoTGatewayRequest的基础抽象类
type BaseIoTGatewayRequest struct {
	ApiName  string
	ApiVer   string
	ApiType  string
	ReqText  string
	TransId  string
	Response IoTGatewayResponse
}

// GetContentType 获取内容类型
func (r *BaseIoTGatewayRequest) GetContentType() string {
	if r.ApiType == "" {
		return ""
	}
	if r.ApiType == API_TYPE_JSON {
		return "application/json"
	} else {
		return "text/xml"
	}
}

// GetApiName 获取API名称
func (r *BaseIoTGatewayRequest) GetApiName() string {
	return r.ApiName
}

// SetApiName 设置API名称
func (r *BaseIoTGatewayRequest) SetApiName(apiName string) {
	r.ApiName = apiName
}

// GetApiVer 获取API版本
func (r *BaseIoTGatewayRequest) GetApiVer() string {
	return r.ApiVer
}

// SetApiVer 设置API版本
func (r *BaseIoTGatewayRequest) SetApiVer(apiVer string) {
	r.ApiVer = apiVer
}

// GetApiType 获取API类型
func (r *BaseIoTGatewayRequest) GetApiType() string {
	return r.ApiType
}

// SetApiType 设置API类型
func (r *BaseIoTGatewayRequest) SetApiType(apiType string) {
	r.ApiType = apiType
}

// GetReqText 获取请求文本
func (r *BaseIoTGatewayRequest) GetReqText() string {
	return r.ReqText
}

// SetReqText 设置请求文本
func (r *BaseIoTGatewayRequest) SetReqText(reqText string) {
	r.ReqText = reqText
}

// GetTransId 获取交易ID
func (r *BaseIoTGatewayRequest) GetTransId() string {
	return r.TransId
}

// SetTransId 设置交易ID
func (r *BaseIoTGatewayRequest) SetTransId(transId string) {
	r.TransId = transId
}
