package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"
)

// SM3相关常量
var (
	ivHexStr = "7380166f4914b2b9172442d7da8a0600a96f30bc163138aae38dee4db0fb0e4e"
	iv, _    = new(big.Int).SetString(ivHexStr, 16)

	tj15 = uint32(0x79cc4519)
	tj63 = uint32(0x7a879d8a)

	firstPadding = []byte{0x80}
	zeroPadding  = []byte{0x00}
)

// Hex字符数组
var hexChars = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

// T 根据j值选择不同的常量
func T(j int) uint32 {
	if j >= 0 && j <= 15 {
		return tj15
	} else if j >= 16 && j <= 63 {
		return tj63
	} else {
		log.Fatal("数据无效")
		return 0
	}
}

// FF 布尔函数FF
func FF(x, y, z uint32, j int) uint32 {
	if j >= 0 && j <= 15 {
		return x ^ y ^ z
	} else if j >= 16 && j <= 63 {
		return (x & y) | (x & z) | (y & z)
	} else {
		log.Fatal("数据无效")
		return 0
	}
}

// GG 布尔函数GG
func GG(x, y, z uint32, j int) uint32 {
	if j >= 0 && j <= 15 {
		return x ^ y ^ z
	} else if j >= 16 && j <= 63 {
		return (x & y) | (^x & z)
	} else {
		log.Fatal("数据无效")
		return 0
	}
}

// P0 置换函数P0
func P0(x uint32) uint32 {
	return x ^ rotateLeft(x, 9) ^ rotateLeft(x, 17)
}

// P1 置换函数P1
func P1(x uint32) uint32 {
	return x ^ rotateLeft(x, 15) ^ rotateLeft(x, 23)
}

// rotateLeft 32位整数循环左移
func rotateLeft(x uint32, n int) uint32 {
	return (x << uint32(n)) | (x >> (32 - uint32(n)))
}

// padding 消息填充
func padding(source []byte) []byte {
	// 检查输入长度
	if int64(len(source)) >= 0x2000000000000000 {
		log.Fatal("输入数据长度超过限制")
	}

	// 计算填充
	l := int64(len(source)) * 8
	k := 448 - (l+1)%512
	if k < 0 {
		k += 512
	}

	// 填充数据
	var buffer bytes.Buffer
	buffer.Write(source)
	buffer.Write(firstPadding)

	// 填充0
	i := k - 7
	for i > 0 {
		buffer.Write(zeroPadding)
		i -= 8
	}

	// 添加长度信息（8字节）
	lenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(lenBytes, uint64(l))
	buffer.Write(lenBytes)

	return buffer.Bytes()
}

// hash SM3哈希计算
func hash(source []byte) ([]byte, error) {
	m1 := padding(source)
	n := len(m1) / 64 // 每块64字节 = 512位

	// 初始IV
	ivBytes := make([]byte, 32)
	ivBig := big.NewInt(0).Set(iv)
	ivBytesTemp := ivBig.Bytes()

	// 确保长度为32字节
	copy(ivBytes[32-len(ivBytesTemp):], ivBytesTemp)

	vi := ivBytes
	var vi1 []byte

	for i := 0; i < n; i++ {
		b := m1[i*64 : (i+1)*64]
		vi1, _ = cf(vi, b)
		vi = vi1
	}

	return vi1, nil
}

// cf SM3压缩函数
func cf(vi, bi []byte) ([]byte, error) {
	// 将字节数组转换为8个32位整数
	a := toInteger(vi, 0)
	b := toInteger(vi, 1)
	c := toInteger(vi, 2)
	d := toInteger(vi, 3)
	e := toInteger(vi, 4)
	f := toInteger(vi, 5)
	g := toInteger(vi, 6)
	h := toInteger(vi, 7)

	w := make([]uint32, 68)
	w1 := make([]uint32, 64)

	// 填充消息扩展
	for i := 0; i < 16; i++ {
		w[i] = toInteger(bi, i)
	}

	for j := 16; j < 68; j++ {
		w[j] = P1(w[j-16]^w[j-9]^rotateLeft(w[j-3], 15)) ^ rotateLeft(w[j-13], 7) ^ w[j-6]
	}

	for j := 0; j < 64; j++ {
		w1[j] = w[j] ^ w[j+4]
	}

	// 压缩函数
	var ss1, ss2, tt1, tt2 uint32

	for j := 0; j < 64; j++ {
		ss1 = rotateLeft(rotateLeft(a, 12)+e+rotateLeft(T(j), j), 7)
		ss2 = ss1 ^ rotateLeft(a, 12)
		tt1 = FF(a, b, c, j) + d + ss2 + w1[j]
		tt2 = GG(e, f, g, j) + h + ss1 + w[j]
		d = c
		c = rotateLeft(b, 9)
		b = a
		a = tt1
		h = g
		g = rotateLeft(f, 19)
		f = e
		e = P0(tt2)
	}

	// 异或处理
	v := toByteArray(a, b, c, d, e, f, g, h)
	for i := 0; i < len(v); i++ {
		v[i] = v[i] ^ vi[i]
	}

	return v, nil
}

// toInteger 将字节数组的一部分转换为uint32
func toInteger(source []byte, index int) uint32 {
	start := index * 4
	return binary.BigEndian.Uint32(source[start : start+4])
}

// toByteArray 将8个uint32转换为字节数组
func toByteArray(a, b, c, d, e, f, g, h uint32) []byte {
	result := make([]byte, 32)
	binary.BigEndian.PutUint32(result[0:], a)
	binary.BigEndian.PutUint32(result[4:], b)
	binary.BigEndian.PutUint32(result[8:], c)
	binary.BigEndian.PutUint32(result[12:], d)
	binary.BigEndian.PutUint32(result[16:], e)
	binary.BigEndian.PutUint32(result[20:], f)
	binary.BigEndian.PutUint32(result[24:], g)
	binary.BigEndian.PutUint32(result[28:], h)
	return result
}

// byteToHexString 将字节转换为16进制字符串
func byteToHexString(b byte) string {
	n := int(b)
	if n < 0 {
		n += 256
	}
	d1 := n / 16
	d2 := n % 16
	return string([]byte{hexChars[d1], hexChars[d2]})
}

// byteArrayToHexString 将字节数组转换为16进制字符串
func byteArrayToHexString(b []byte) string {
	var result bytes.Buffer
	for i := 0; i < len(b); i++ {
		result.WriteString(byteToHexString(b[i]))
	}
	return result.String()
}

// SM3Encode 计算字符串的SM3摘要
func SM3Encode(source, charset string) (string, error) {
	var resultString string
	var hashBytes []byte
	var err error

	if charset == "" {
		hashBytes, err = hash([]byte(source))
	} else {
		// Go语言中字符集编码处理不同于Java
		// 这里简化处理，实际使用中可能需要根据具体情况调整
		hashBytes, err = hash([]byte(source))
	}

	if err != nil {
		return "", err
	}

	resultString = byteArrayToHexString(hashBytes)
	return resultString, nil
}

// MakeToken 生成令牌
func MakeToken(appID, appSecret string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 获取时间戳格式化字符串
	timestamp := getFormatedDateString(8)

	// 生成transId
	transID := timestamp
	transID = strconv.FormatInt(time.Now().UnixNano(), 10)

	// 构建签名字符串
	sb := fmt.Sprintf("app_id%stimestamp%strans_id%s%s",
		appID, timestamp, transID, appSecret)

	// 计算SM3摘要
	token, err := SM3Encode(sb, "UTF-8")
	if err != nil {
		return nil, err
	}

	result["token"] = token
	result["timestamp"] = timestamp
	result["trans_id"] = transID

	return result, nil
}

// BuildAppParams 构建应用参数
func BuildAppParams(params map[string]interface{}) (map[string]interface{}, error) {
	appID := params["app_id"].(string)
	appSecret := params["app_secrect"].(string)

	result, err := MakeToken(appID, appSecret)
	if err != nil {
		return nil, err
	}

	params["token"] = result["token"]
	params["trans_id"] = result["trans_id"]
	params["timestamp"] = result["timestamp"]
	delete(params, "app_secrect")

	return result, nil
}

// getFormatedDateString 获取格式化的日期字符串
func getFormatedDateString(timeZoneOffset float64) string {
	if timeZoneOffset > 13 || timeZoneOffset < -12 {
		timeZoneOffset = 0
	}

	// 计算时区偏移的毫秒数
	offsetHours := int(math.Floor(timeZoneOffset))
	offsetMinutes := int((timeZoneOffset - float64(offsetHours)) * 60)

	// 获取当前时间
	now := time.Now()

	// 应用时区偏移
	location := time.FixedZone("Custom", offsetHours*3600+offsetMinutes*60)
	now = now.In(location)

	// 格式化时间
	return now.Format("2006-01-02 15:04:05.000")
}

// GenerateRandomNumber 生成指定位数的随机数字字符串
func GenerateRandomNumber(digits int) (string, error) {
	max := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// 确保数字长度正确（左侧补0）
	s := fmt.Sprintf("%0"+strconv.Itoa(digits)+"s", n.String())
	return s[len(s)-digits:], nil
}
