package api

import (
	"encoding/json"
	"fmt"

	"github.com/zhoudm1743/unicom-gw/api/internal/utils"
)

// IoTGatewayClient 定义IoT网关客户端接口
type IoTGatewayClient interface {
	// Execute 执行API请求
	Execute(request IoTGatewayRequest) (IoTGatewayResponse, error)

	// GetServerURL 获取服务器URL
	GetServerURL() string

	// GetAppID 获取应用ID
	GetAppID() string

	// GetAppSecret 获取应用密钥
	GetAppSecret() string

	// GetConnectTimeout 获取连接超时时间
	GetConnectTimeout() int

	// GetReadTimeout 获取读取超时时间
	GetReadTimeout() int

	// SetConnectTimeout 设置连接超时时间
	SetConnectTimeout(timeout int)

	// SetReadTimeout 设置读取超时时间
	SetReadTimeout(timeout int)

	// GetRetryCount 获取重试次数
	GetRetryCount() int

	// SetRetryCount 设置重试次数
	SetRetryCount(retryCount int)

	// GetOpenID 获取openID
	GetOpenID() string

	// SetOpenID 设置openID
	SetOpenID(openID string)
}

// DefaultIoTGatewayClient 默认IoT网关客户端实现
type DefaultIoTGatewayClient struct {
	ServerURL      string
	AppID          string
	AppSecret      string
	OpenID         string
	ConnectTimeout int
	ReadTimeout    int
	RetryCount     int
}

// NewIoTGatewayClient 创建一个新的IoT网关客户端
func NewIoTGatewayClient(serverURL, appID, appSecret, openID string) *DefaultIoTGatewayClient {
	return &DefaultIoTGatewayClient{
		ServerURL:      serverURL,
		AppID:          appID,
		AppSecret:      appSecret,
		OpenID:         openID,
		ConnectTimeout: 2000,  // 默认连接超时2秒
		ReadTimeout:    30000, // 默认读取超时30秒
		RetryCount:     0,     // 默认不重试
	}
}

// Execute 执行API请求
func (c *DefaultIoTGatewayClient) Execute(request IoTGatewayRequest) (IoTGatewayResponse, error) {
	// 执行POST请求
	respMsg, err := c.doPost(request)
	if err != nil {
		return nil, err
	}

	if respMsg == "" {
		return nil, nil
	}

	// 根据内容类型选择解析器
	contentType := request.GetContentType()
	responseClass := request.GetResponseClass()

	// 解析响应
	var response IoTGatewayResponse
	if contentType != "" && contentType == "text/xml" {
		// XML解析（实际项目中可能需要实现）
		return nil, fmt.Errorf("XML解析尚未实现")
	} else {
		// JSON解析
		err = json.Unmarshal([]byte(respMsg), responseClass)
		if err != nil {
			responseClass.SetSuccess(false)
		}
		response = responseClass
	}

	return response, nil
}

// doPost 执行POST请求
func (c *DefaultIoTGatewayClient) doPost(request IoTGatewayRequest) (string, error) {
	// 构建请求参数
	params := map[string]interface{}{
		utils.AppIDKey:     c.AppID,
		utils.AppSecretKey: c.AppSecret,
	}

	// 构建应用参数
	err := utils.BuildAppParams(params)
	if err != nil {
		return "", &ApiException{
			ErrMsg: "构建应用参数失败",
			Cause:  err,
		}
	}

	// 获取请求参数
	requestParams := request.GetParams()
	params["data"] = requestParams

	// 请求发送前的处理
	request.ExecProcessBeforeReqSend([]interface{}{params})

	// 发送请求
	return utils.DoPost(
		c.ServerURL,
		request.GetApiName(),
		request.GetApiVer(),
		request.GetReqText(),
		c.ConnectTimeout,
		c.ReadTimeout,
		c.RetryCount,
	)
}

// GetServerURL 获取服务器URL
func (c *DefaultIoTGatewayClient) GetServerURL() string {
	return c.ServerURL
}

// GetAppID 获取应用ID
func (c *DefaultIoTGatewayClient) GetAppID() string {
	return c.AppID
}

// GetAppSecret 获取应用密钥
func (c *DefaultIoTGatewayClient) GetAppSecret() string {
	return c.AppSecret
}

// GetConnectTimeout 获取连接超时时间
func (c *DefaultIoTGatewayClient) GetConnectTimeout() int {
	return c.ConnectTimeout
}

// GetReadTimeout 获取读取超时时间
func (c *DefaultIoTGatewayClient) GetReadTimeout() int {
	return c.ReadTimeout
}

// SetConnectTimeout 设置连接超时时间
func (c *DefaultIoTGatewayClient) SetConnectTimeout(timeout int) {
	c.ConnectTimeout = timeout
}

// SetReadTimeout 设置读取超时时间
func (c *DefaultIoTGatewayClient) SetReadTimeout(timeout int) {
	c.ReadTimeout = timeout
}

// GetRetryCount 获取重试次数
func (c *DefaultIoTGatewayClient) GetRetryCount() int {
	return c.RetryCount
}

// SetRetryCount 设置重试次数
func (c *DefaultIoTGatewayClient) SetRetryCount(retryCount int) {
	c.RetryCount = retryCount
}

// GetOpenID 获取openID
func (c *DefaultIoTGatewayClient) GetOpenID() string {
	return c.OpenID
}

// SetOpenID 设置openID
func (c *DefaultIoTGatewayClient) SetOpenID(openID string) {
	c.OpenID = openID
}
