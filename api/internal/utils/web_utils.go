package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

const (
	// HTTP方法
	MethodPost = "POST"
	MethodGet  = "GET"

	// 默认字符集
	DefaultCharset = "UTF-8"
)

// TrustAllCerts 信任所有证书
var TrustAllCerts = &tls.Config{
	InsecureSkipVerify: true,
}

// HTTPClient HTTP客户端，可根据需要进行配置
var HTTPClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: TrustAllCerts,
	},
}

// DoPost 执行HTTP POST请求
func DoPost(serverURL, apiName, apiVersion, reqText string, connectTimeout, readTimeout, retryCount int) (string, error) {
	// 确保URL以"/"结尾
	if !strings.HasSuffix(serverURL, "/") {
		serverURL += "/"
	}

	// 去除"cn."前缀（如果存在）
	if strings.HasPrefix(apiName, "cn.") {
		apiName = apiName[3:]
	}

	// 构建完整URL
	fullURL := serverURL + strings.ReplaceAll(apiName, ".", "/") + "/v" + apiVersion
	contentType := "application/json;charset=" + DefaultCharset

	// 执行POST请求
	return doPost(fullURL, contentType, []byte(reqText), connectTimeout, readTimeout, retryCount)
}

// DoPostWithParams 带参数执行HTTP POST请求
func DoPostWithParams(serverURL, apiName, apiVersion string, params map[string]interface{}, connectTimeout, readTimeout, retryCount int) (string, error) {
	// 确保URL以"/"结尾
	if !strings.HasSuffix(serverURL, "/") {
		serverURL += "/"
	}

	// 去除"cn."前缀（如果存在）
	if strings.HasPrefix(apiName, "cn.") {
		apiName = apiName[3:]
	}

	// 构建完整URL
	fullURL := serverURL + strings.ReplaceAll(apiName, ".", "/") + "/v" + apiVersion
	contentType := "application/json;charset=" + DefaultCharset

	// 将参数转换为JSON
	jsonData, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	// 执行POST请求
	return doPost(fullURL, contentType, jsonData, connectTimeout, readTimeout, retryCount)
}

// doPost 执行HTTP POST请求（内部方法）
func doPost(url, contentType string, content []byte, connectTimeout, readTimeout, retryCount int) (string, error) {
	var resp string
	var err error

	// 重试逻辑
	for i := 0; i <= retryCount; i++ {
		resp, err = executePost(url, contentType, content, connectTimeout, readTimeout)
		if err == nil {
			return resp, nil
		}

		// 最后一次重试失败，直接返回错误
		if i == retryCount {
			return "", err
		}
	}

	return resp, err
}

// executePost 执行单个HTTP POST请求
func executePost(urlStr, contentType string, content []byte, connectTimeout, readTimeout int) (string, error) {
	// 创建请求
	req, err := http.NewRequest(MethodPost, urlStr, bytes.NewBuffer(content))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "iot-gateway-sdk-go")
	req.Header.Set("Accept", "text/xml,text/javascript")
	fmt.Println(fmt.Sprintf("请求url: %s \n 请求体: %s", req.URL.String(), string(content)))
	// 设置超时
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: TrustAllCerts,
		},
		Timeout: time.Duration(connectTimeout+readTimeout) * time.Millisecond,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP错误: %d %s", resp.StatusCode, resp.Status)
	}

	// 处理gzip压缩
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return "", err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	// 读取响应内容
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// DoGet 执行HTTP GET请求
func DoGet(urlStr string, params map[string]string) (string, error) {
	return DoGetWithCharset(urlStr, params, DefaultCharset)
}

// DoGetWithCharset 带字符集执行HTTP GET请求
func DoGetWithCharset(urlStr string, params map[string]string, charset string) (string, error) {
	// 构建完整URL
	fullURL, err := BuildGetURL(urlStr, params, charset)
	if err != nil {
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest(MethodGet, fullURL, nil)
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("User-Agent", "iot-gateway-sdk-go")
	req.Header.Set("Accept", "text/xml,text/javascript")

	// 发送请求
	resp, err := HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP错误: %d %s", resp.StatusCode, resp.Status)
	}

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// BuildGetURL 构建GET请求URL
func BuildGetURL(urlStr string, params map[string]string, charset string) (string, error) {
	// 构建查询字符串
	if len(params) == 0 {
		return urlStr, nil
	}

	query := url.Values{}
	for k, v := range params {
		if k != "" && v != "" {
			query.Add(k, v)
		}
	}

	queryStr := query.Encode()
	if queryStr == "" {
		return urlStr, nil
	}

	// 根据URL是否已有查询参数，添加"?"或"&"
	if strings.Contains(urlStr, "?") {
		if strings.HasSuffix(urlStr, "?") || strings.HasSuffix(urlStr, "&") {
			return urlStr + queryStr, nil
		}
		return urlStr + "&" + queryStr, nil
	}
	return urlStr + "?" + queryStr, nil
}

// DoPostWithFile 执行带文件上传的HTTP POST请求
func DoPostWithFile(urlStr string, textParams map[string]string, fileParams map[string]*FileItem, charset string, connectTimeout, readTimeout int) (string, error) {
	// 如果没有文件参数，执行普通POST请求
	if len(fileParams) == 0 {
		// 转换参数类型
		params := make(map[string]string)
		for k, v := range textParams {
			params[k] = v
		}

		// 构建查询字符串
		values := url.Values{}
		for k, v := range params {
			if k != "" && v != "" {
				values.Add(k, v)
			}
		}

		// 创建请求
		req, err := http.NewRequest(MethodPost, urlStr, strings.NewReader(values.Encode()))
		if err != nil {
			return "", err
		}

		// 设置请求头
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset="+charset)
		req.Header.Set("User-Agent", "iot-gateway-sdk-go")
		req.Header.Set("Accept", "text/xml,text/javascript")

		// 设置客户端超时
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: TrustAllCerts,
			},
			Timeout: time.Duration(connectTimeout+readTimeout) * time.Millisecond,
		}

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		// 检查响应状态
		if resp.StatusCode >= 400 {
			return "", fmt.Errorf("HTTP错误: %d %s", resp.StatusCode, resp.Status)
		}

		// 读取响应内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(body), nil
	}

	// 创建多部分表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文本参数
	for key, val := range textParams {
		if err := writer.WriteField(key, val); err != nil {
			return "", err
		}
	}

	// 添加文件参数
	for key, fileItem := range fileParams {
		if fileItem != nil && len(fileItem.Content) > 0 {
			part, err := writer.CreateFormFile(key, filepath.Base(fileItem.FileName))
			if err != nil {
				return "", err
			}

			if _, err = part.Write(fileItem.Content); err != nil {
				return "", err
			}
		}
	}

	// 关闭writer
	if err := writer.Close(); err != nil {
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest(MethodPost, urlStr, body)
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "iot-gateway-sdk-go")
	req.Header.Set("Accept", "text/xml,text/javascript")

	// 设置客户端超时
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: TrustAllCerts,
		},
		Timeout: time.Duration(connectTimeout+readTimeout) * time.Millisecond,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP错误: %d %s", resp.StatusCode, resp.Status)
	}

	// 读取响应内容
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// FileItem 文件项目
type FileItem struct {
	FileName string
	Content  []byte
	MimeType string
}
