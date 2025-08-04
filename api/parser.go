package api

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"reflect"
)

// IoTGatewayParser 定义解析器接口
type IoTGatewayParser interface {
	// Parse 解析响应字符串为对应的响应对象
	Parse(respStr string) (IoTGatewayResponse, error)
}

// JSONParser JSON解析器
type JSONParser struct {
	ResponseClass interface{}
}

// Parse 解析JSON响应为对应的响应对象
func (p *JSONParser) Parse(respStr string) (IoTGatewayResponse, error) {
	// 使用反射创建响应对象实例
	respType := reflect.TypeOf(p.ResponseClass)
	if respType.Kind() == reflect.Ptr {
		respType = respType.Elem()
	}

	respInstance := reflect.New(respType).Interface()

	// 解析JSON字符串到响应对象
	err := json.Unmarshal([]byte(respStr), respInstance)
	if err != nil {
		return nil, errors.New("解析JSON响应失败: " + err.Error())
	}

	// 转换为IoTGatewayResponse接口类型
	resp, ok := respInstance.(IoTGatewayResponse)
	if !ok {
		return nil, errors.New("响应对象不实现IoTGatewayResponse接口")
	}

	// 设置原始响应体
	resp.SetBody(respStr)

	return resp, nil
}

// XMLParser XML解析器
type XMLParser struct {
	ResponseClass interface{}
}

// Parse 解析XML响应为对应的响应对象
func (p *XMLParser) Parse(respStr string) (IoTGatewayResponse, error) {
	// 使用反射创建响应对象实例
	respType := reflect.TypeOf(p.ResponseClass)
	if respType.Kind() == reflect.Ptr {
		respType = respType.Elem()
	}

	respInstance := reflect.New(respType).Interface()

	// 解析XML字符串到响应对象
	err := xml.Unmarshal([]byte(respStr), respInstance)
	if err != nil {
		return nil, errors.New("解析XML响应失败: " + err.Error())
	}

	// 转换为IoTGatewayResponse接口类型
	resp, ok := respInstance.(IoTGatewayResponse)
	if !ok {
		return nil, errors.New("响应对象不实现IoTGatewayResponse接口")
	}

	// 设置原始响应体
	resp.SetBody(respStr)

	return resp, nil
}
