package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/tjfoc/gmsm/sm3"
)

// SM3Encode 对输入字符串进行SM3摘要，返回十六进制小写字符串
func SM3Encode(source, charset string) (string, error) {
	var data []byte
	var err error

	if charset == "" {
		data = []byte(source)
	} else {
		// Go中字符集处理比较简单，这里仅示例，实际可能需要更复杂的处理
		data = []byte(source)
	}

	hash := sm3.Sm3Sum(data)
	return fmt.Sprintf("%x", hash), err
}

// MakeToken 创建令牌及相关参数
func MakeToken(appID, appSecret string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 获取格式化时间戳
	timestamp := GetFormattedDateString(8)

	// 从时间戳创建事务ID
	transID := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(timestamp, "-", ""), " ", ""), ":", "")
	transID = transID + fmt.Sprintf("%d", rand.Intn(999999))

	// 构建原始字符串
	var sb strings.Builder
	sb.WriteString("app_id")
	sb.WriteString(appID)
	sb.WriteString("timestamp")
	sb.WriteString(timestamp)
	sb.WriteString("trans_id")
	sb.WriteString(transID)
	sb.WriteString(appSecret)

	// 计算SM3摘要
	token, err := SM3Encode(sb.String(), "UTF-8")
	if err != nil {
		return nil, err
	}

	result["token"] = token
	result["timestamp"] = timestamp
	result["trans_id"] = transID
	return result, nil
}

// GetFormattedDateString 获取格式化的日期字符串
// timeZoneOffset 表示时区，如中国一般使用东八区，因此timeZoneOffset就是8
func GetFormattedDateString(timeZoneOffset float32) string {
	if timeZoneOffset > 13 || timeZoneOffset < -12 {
		timeZoneOffset = 0
	}

	// 设置时区
	loc := time.FixedZone(fmt.Sprintf("UTC%+.1f", timeZoneOffset), int(timeZoneOffset*3600))

	// 格式化时间
	now := time.Now().In(loc)
	return now.Format("2006-01-02 15:04:05 000")
}
