package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// 默认字符集
	DEFAULT_CHARSET = "UTF-8"
	// HTTP方法
	METHOD_POST = "POST"
	METHOD_GET  = "GET"
	// 内容编码类型
	CONTENT_ENCODING_GZIP = "gzip"
)

var (
	// 忽略SSL检查
	ignoreSSLCheck = true
	// 忽略HOST检查
	ignoreHostCheck = false
)

// FileItem 文件项
type FileItem struct {
	FileName string
	Content  []byte
	MimeType string
}

// SetIgnoreSSLCheck 设置是否忽略SSL检查
func SetIgnoreSSLCheck(ignore bool) {
	ignoreSSLCheck = ignore
}

// SetIgnoreHostCheck 设置是否忽略HOST检查
func SetIgnoreHostCheck(ignore bool) {
	ignoreHostCheck = ignore
}

// DoPost 执行HTTP POST请求
func DoPost(serverURL string, apiName string, apiVersion string, reqText string, connectTimeout time.Duration, readTimeout time.Duration, retryCount int) (string, error) {
	// 构建URL路径
	if !strings.HasSuffix(serverURL, "/") {
		serverURL += "/"
	}
	if strings.HasPrefix(apiName, "cn.") {
		apiName = apiName[3:]
	}
	serverURL += strings.ReplaceAll(apiName, ".", "/")
	serverURL += "/v" + apiVersion

	// 构建请求内容类型和请求体
	ctype := "application/json;charset=" + DEFAULT_CHARSET
	reqBody := []byte(reqText)

	// 执行POST请求
	return doPost(serverURL, ctype, reqBody, connectTimeout, readTimeout, retryCount)
}

// DoPostWithParams 使用参数执行HTTP POST请求
func DoPostWithParams(serverURL string, apiName string, apiVersion string, params map[string]interface{}, connectTimeout time.Duration, readTimeout time.Duration, retryCount int) (string, error) {
	// 构建URL路径
	if !strings.HasSuffix(serverURL, "/") {
		serverURL += "/"
	}
	if strings.HasPrefix(apiName, "cn.") {
		apiName = apiName[3:]
	}
	serverURL += strings.ReplaceAll(apiName, ".", "/")
	serverURL += "/v" + apiVersion

	// 构建请求内容类型和请求体
	ctype := "application/json;charset=" + DEFAULT_CHARSET

	// 将参数转换为JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	// 执行POST请求
	return doPost(serverURL, ctype, jsonData, connectTimeout, readTimeout, retryCount)
}

// doPost 执行HTTP POST请求（内部方法）
func doPost(url string, contentType string, reqBody []byte, connectTimeout time.Duration, readTimeout time.Duration, retryCount int) (string, error) {
	var resp string
	var err error

	// 重试逻辑
	for i := 0; i <= retryCount; i++ {
		resp, err = performPost(url, contentType, reqBody, connectTimeout, readTimeout)
		if err == nil {
			return resp, nil
		}

		// 如果是最后一次尝试且失败，则返回错误
		if i == retryCount {
			return "", err
		}
	}

	return resp, err
}

// performPost 执行单次HTTP POST请求
func performPost(url string, contentType string, reqBody []byte, connectTimeout time.Duration, readTimeout time.Duration) (string, error) {
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: connectTimeout + readTimeout,
	}

	// 配置TLS
	if strings.HasPrefix(url, "https://") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ignoreSSLCheck,
			},
		}
		client.Transport = tr
	}

	// 创建请求
	req, err := http.NewRequest(METHOD_POST, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "iot-sdk-go")
	req.Header.Set("Accept", "text/xml,text/javascript,application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	// 处理响应体
	var body []byte
	if resp.Header.Get("Content-Encoding") == CONTENT_ENCODING_GZIP {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", err
		}
		defer reader.Close()
		body, err = ioutil.ReadAll(reader)
	} else {
		body, err = ioutil.ReadAll(resp.Body)
	}

	if err != nil {
		return "", err
	}

	return string(body), nil
}

// DoPostWithFiles 执行带文件上传的HTTP POST请求
func DoPostWithFiles(targetURL string, textParams map[string]string, fileParams map[string]FileItem, connectTimeout time.Duration, readTimeout time.Duration) (string, error) {
	// 如果没有文件参数，则使用普通POST请求
	if len(fileParams) == 0 {
		return DoPostForm(targetURL, textParams, connectTimeout, readTimeout)
	}

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文本参数
	for key, val := range textParams {
		err := writer.WriteField(key, val)
		if err != nil {
			return "", err
		}
	}

	// 添加文件参数
	for key, item := range fileParams {
		if item.Content == nil || len(item.Content) == 0 {
			continue
		}

		part, err := writer.CreateFormFile(key, item.FileName)
		if err != nil {
			return "", err
		}

		_, err = part.Write(item.Content)
		if err != nil {
			return "", err
		}
	}

	err := writer.Close()
	if err != nil {
		return "", err
	}

	// 设置内容类型
	contentType := writer.FormDataContentType()

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: connectTimeout + readTimeout,
	}

	// 配置TLS
	if strings.HasPrefix(targetURL, "https://") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ignoreSSLCheck,
			},
		}
		client.Transport = tr
	}

	// 创建请求
	req, err := http.NewRequest(METHOD_POST, targetURL, body)
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "iot-sdk-go")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// DoPostForm 执行HTTP POST表单请求
func DoPostForm(targetURL string, params map[string]string, connectTimeout time.Duration, readTimeout time.Duration) (string, error) {
	// 创建URL编码的表单数据
	formData := url.Values{}
	for key, val := range params {
		formData.Add(key, val)
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: connectTimeout + readTimeout,
	}

	// 配置TLS
	if strings.HasPrefix(targetURL, "https://") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ignoreSSLCheck,
			},
		}
		client.Transport = tr
	}

	// 创建请求
	req, err := http.NewRequest(METHOD_POST, targetURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "iot-sdk-go")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// DoGet 执行HTTP GET请求
func DoGet(targetURL string, params map[string]string) (string, error) {
	// 构建带参数的URL
	queryURL, err := buildGetURL(targetURL, params)
	if err != nil {
		return "", err
	}

	// 创建HTTP客户端
	client := &http.Client{}

	// 配置TLS
	if strings.HasPrefix(targetURL, "https://") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ignoreSSLCheck,
			},
		}
		client.Transport = tr
	}

	// 创建请求
	req, err := http.NewRequest(METHOD_GET, queryURL, nil)
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("User-Agent", "iot-sdk-go")
	req.Header.Set("Accept", "text/xml,text/javascript,application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// buildGetURL 构建带参数的GET URL
func buildGetURL(baseURL string, params map[string]string) (string, error) {
	// 如果没有参数，直接返回原URL
	if len(params) == 0 {
		return baseURL, nil
	}

	// 解析URL
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// 创建查询参数
	q := u.Query()
	for key, val := range params {
		q.Add(key, val)
	}

	// 更新URL的查询参数
	u.RawQuery = q.Encode()

	return u.String(), nil
}
