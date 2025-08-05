package api

// 公用常量

const (
	// 请求相关参数
	APP_KEY     = "app_key"
	TIMESTAMP   = "timestamp"
	VERSION     = "v"
	SIGN        = "sign"
	SIGN_METHOD = "sign_method"

	// 日期时间格式
	DATE_TIME_FORMAT = "2006-01-02 15:04:05"
	DATE_TIMEZONE    = "GMT+8"

	// 字符集
	CHARSET_UTF8 = "UTF-8"
	CHARSET_GBK  = "GBK"

	// 签名方法
	SIGN_METHOD_MD5  = "md5"
	SIGN_METHOD_HMAC = "hmac"

	// SDK版本
	SDK_VERSION = "iot-gateway-sdk-go-20240101"

	// 响应编码
	ERROR_CODE = "status"
	ERROR_MSG  = "message"

	// HTTP头
	ACCEPT_ENCODING       = "Accept-Encoding"
	CONTENT_ENCODING      = "Content-Encoding"
	CONTENT_ENCODING_GZIP = "gzip"

	// API类型
	API_TYPE_JSON = "json"
	API_TYPE_WS   = "ws"
)
