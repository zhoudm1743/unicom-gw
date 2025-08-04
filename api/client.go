package api

import (
	"time"

	"github.com/zhoudm1743/unicom-gw/api/internal/utils"
)

// IoTGatewayClient 接口定义
type IoTGatewayClient interface {
	Execute(request IoTGatewayRequest) (IoTGatewayResponse, error)
}

// DefaultIoTGatewayClient 实现IoTGatewayClient接口
type DefaultIoTGatewayClient struct {
	ServerURL      string
	AppID          string
	AppSecret      string
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	RetryCount     int
}

// NewIoTGatewayClient 创建一个新的IoTGatewayClient实例
func NewIoTGatewayClient(serverURL, appID, appSecret string) *DefaultIoTGatewayClient {
	return &DefaultIoTGatewayClient{
		ServerURL:      serverURL,
		AppID:          appID,
		AppSecret:      appSecret,
		ConnectTimeout: 2 * time.Second,
		ReadTimeout:    30 * time.Second,
		RetryCount:     0, // 连接超时后，重试次数 0：不重试 1：重试一次 以此类推。默认不重试
	}
}

// Execute 执行API请求
func (c *DefaultIoTGatewayClient) Execute(request IoTGatewayRequest) (IoTGatewayResponse, error) {
	respMsg, err := c.doPost(request)
	if err != nil {
		return nil, err
	}
	if respMsg == "" {
		return nil, nil
	}

	contentType := request.GetContentType()
	respClass := request.GetResponseClass()

	var parser IoTGatewayParser
	if contentType != "" && contentType == "text/xml" {
		parser = &XMLParser{respClass}
	} else {
		parser = &JSONParser{respClass}
	}

	resp, err := parser.Parse(respMsg)
	if err != nil {
		resp.SetSuccess(false)
	}

	return resp, nil
}

func (c *DefaultIoTGatewayClient) doPost(request IoTGatewayRequest) (string, error) {
	params := make(map[string]interface{})
	params["app_id"] = c.AppID
	params["app_secrect"] = c.AppSecret

	// 构建应用参数
	tokenInfo, err := utils.BuildAppParams(params)
	if err != nil {
		return "", err
	}

	// 确保token信息被正确应用到params中
	params["token"] = tokenInfo["token"]
	params["timestamp"] = tokenInfo["timestamp"]
	params["trans_id"] = tokenInfo["trans_id"]

	reqParams := request.GetParams()
	params["data"] = reqParams

	// 执行请求前的处理
	request.ExecProcessBeforeReqSend([]interface{}{params})

	// 构建请求URL
	apiName := request.GetApiName()
	apiVer := request.GetApiVer()
	reqText := request.GetReqText()

	// 发送HTTP请求
	resp, err := utils.DoPost(c.ServerURL, apiName, apiVer, reqText, c.ConnectTimeout, c.ReadTimeout, c.RetryCount)
	if err != nil {
		return "", err
	}

	return resp, nil
}

// GetConnectTimeout 获取连接超时时间
func (c *DefaultIoTGatewayClient) GetConnectTimeout() time.Duration {
	return c.ConnectTimeout
}

// GetReadTimeout 获取读取超时时间
func (c *DefaultIoTGatewayClient) GetReadTimeout() time.Duration {
	return c.ReadTimeout
}

// SetConnectTimeout 设置连接超时时间
func (c *DefaultIoTGatewayClient) SetConnectTimeout(timeout time.Duration) {
	c.ConnectTimeout = timeout
}

// SetReadTimeout 设置读取超时时间
func (c *DefaultIoTGatewayClient) SetReadTimeout(timeout time.Duration) {
	c.ReadTimeout = timeout
}

// GetRetryCount 获取重试次数
func (c *DefaultIoTGatewayClient) GetRetryCount() int {
	return c.RetryCount
}

// SetRetryCount 设置重试次数
func (c *DefaultIoTGatewayClient) SetRetryCount(count int) {
	c.RetryCount = count
}

// GetServerUrl 获取服务器URL
func (c *DefaultIoTGatewayClient) GetServerUrl() string {
	return c.ServerURL
}

// GetAppId 获取应用ID
func (c *DefaultIoTGatewayClient) GetAppId() string {
	return c.AppID
}

// GetAppSecret 获取应用密钥
func (c *DefaultIoTGatewayClient) GetAppSecret() string {
	return c.AppSecret
}
