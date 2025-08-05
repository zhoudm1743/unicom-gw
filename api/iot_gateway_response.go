package api

// IoTGatewayResponse 定义IoT网关响应接口
type IoTGatewayResponse interface {
	// IsSuccess 请求是否成功
	IsSuccess() bool

	// SetSuccess 设置请求是否成功
	SetSuccess(success bool)

	// GetStatus 获取状态码
	GetStatus() string

	// SetStatus 设置状态码
	SetStatus(status string)

	// GetMessage 获取消息
	GetMessage() string

	// SetMessage 设置消息
	SetMessage(message string)

	// GetBody 获取响应体
	GetBody() string

	// SetBody 设置响应体
	SetBody(body string)
}

// BaseIoTGatewayResponse IoT网关响应基础实现
type BaseIoTGatewayResponse struct {
	Status  string
	Message string
	Body    string
	IsSucc  bool
}

// IsSuccess 请求是否成功
func (r *BaseIoTGatewayResponse) IsSuccess() bool {
	return r.IsSucc
}

// SetSuccess 设置请求是否成功
func (r *BaseIoTGatewayResponse) SetSuccess(success bool) {
	r.IsSucc = success
}

// GetStatus 获取状态码
func (r *BaseIoTGatewayResponse) GetStatus() string {
	return r.Status
}

// SetStatus 设置状态码
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
