package api

// IoTGatewayResponse 定义API响应接口
type IoTGatewayResponse interface {
	// GetStatus 获取状态
	GetStatus() string

	// SetStatus 设置状态
	SetStatus(status string)

	// GetMessage 获取消息
	GetMessage() string

	// SetMessage 设置消息
	SetMessage(message string)

	// GetBody 获取响应体
	GetBody() string

	// SetBody 设置响应体
	SetBody(body string)

	// IsSuccess 是否成功
	IsSuccess() bool

	// SetSuccess 设置是否成功
	SetSuccess(success bool)
}

// BaseIoTGatewayResponse 实现IoTGatewayResponse的基础抽象类
type BaseIoTGatewayResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Body      string `json:"-"`
	IsSucceed *bool  `json:"-"`
}

// GetStatus 获取状态
func (r *BaseIoTGatewayResponse) GetStatus() string {
	return r.Status
}

// SetStatus 设置状态
func (r *BaseIoTGatewayResponse) SetStatus(status string) {
	r.Status = status
}

// GetMessage 获取消息
func (r *BaseIoTGatewayResponse) GetMessage() string {
	return r.Message
}

// SetMessage 设置消息
func (r *BaseIoTGatewayResponse) SetMessage(message string) {
	r.Message = message
}

// GetBody 获取响应体
func (r *BaseIoTGatewayResponse) GetBody() string {
	return r.Body
}

// SetBody 设置响应体
func (r *BaseIoTGatewayResponse) SetBody(body string) {
	r.Body = body
}

// IsSuccess 是否成功
func (r *BaseIoTGatewayResponse) IsSuccess() bool {
	if r.IsSucceed != nil {
		return *r.IsSucceed
	}
	return false
}

// SetSuccess 设置是否成功
func (r *BaseIoTGatewayResponse) SetSuccess(success bool) {
	r.IsSucceed = &success
}
