package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	// 请求参数常量
	AppIDKey     = "app_id"
	TransIDKey   = "trans_id"
	TimestampKey = "timestamp"
	TokenKey     = "token"
	AppSecretKey = "app_secrect"
)

// BuildAppParams 构建应用参数
func BuildAppParams(params map[string]interface{}) error {
	appID, ok := params[AppIDKey].(string)
	if !ok {
		return fmt.Errorf("app_id 参数类型不正确")
	}

	appSecret, ok := params[AppSecretKey].(string)
	if !ok {
		return fmt.Errorf("app_secrect 参数类型不正确")
	}

	// 使用SM3生成token和其他参数
	result, err := MakeToken(appID, appSecret)
	if err != nil {
		return err
	}

	// 设置参数
	params["token"] = result["token"]
	params["trans_id"] = result["trans_id"]
	params["timestamp"] = result["timestamp"]

	// 移除密钥参数
	delete(params, AppSecretKey)

	return nil
}

// CreateSign 创建签名
func CreateSign(characterEncoding string, parameters map[string]interface{}) (string, error) {
	// 按照键排序参数
	var keys []string
	for k := range parameters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建签名字符串
	var sb strings.Builder
	for _, k := range keys {
		v := parameters[k]
		if v == nil || (AppSecretKey == k) {
			continue
		}

		// 处理字符串数组
		if strArray, ok := v.([]string); ok && len(strArray) > 0 {
			v = strArray[0]
		}

		// 将值转换为字符串并追加到签名字符串
		if str, ok := v.(string); ok && str != "" {
			sb.WriteString(k)
			sb.WriteString(str)
		}
	}

	// 添加密钥
	if appSecret, ok := parameters[AppSecretKey].(string); ok {
		sb.WriteString(appSecret)
	}

	// 计算MD5并进行Base64编码
	return getMD5Encode(characterEncoding, sb.String())
}

// getMD5Encode MD5摘要计算并进行Base64编码
func getMD5Encode(characterEncoding string, md5Code string) (string, error) {
	h := md5.New()
	h.Write([]byte(md5Code))
	sum := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum), nil
}

// GetTimestampFormat 返回时间戳格式
func GetTimestampFormat() string {
	return "2006-01-02 15:04:05.000"
}

// GetTimestamp 获取当前时间戳字符串
func GetTimestamp() string {
	return time.Now().Format(GetTimestampFormat())
}
