# IoT Gateway SDK for Go

这是中国联通IoT网关的Go语言SDK实现，基于Java SDK进行移植。

## 功能特性

- 支持标准的HTTP请求
- 支持JSON请求和响应格式
- 支持SM3签名算法
- 支持超时设置和重试机制
- 简洁易用的API

## 安装

```bash
go get github.com/zhoudm1743/unicom-gw/gosdk
```

## 依赖

- Go 1.16+
- github.com/tjfoc/gmsm v1.4.1 (用于SM3算法)

## 快速开始

以下是一个简单的示例，展示如何使用SDK发送请求：

```go
package main

import (
    "fmt"
    "log"

    "github.com/zhoudm1743/unicom-gw/gosdk/api"
    "github.com/zhoudm1743/unicom-gw/gosdk/api/request"
    "github.com/zhoudm1743/unicom-gw/gosdk/api/response"
)

func main() {
    // 创建客户端
    client := api.NewIoTGatewayClient("https://gwtest.10646.cn/api/", "your_app_id", "your_app_secret")
    
    // 设置连接超时（可选）
    client.SetConnectTimeout(2000)
    
    // 创建请求
    req := request.NewCommonJsonRequest()
    req.SetApiName("YourApiName")
    req.SetApiVer("V1.0")
    
    // 设置业务参数
    params := map[string]interface{}{
        "param1": "value1",
        "param2": "value2",
    }
    req.SetParams(params)
    
    // 执行请求
    resp, err := client.Execute(req)
    if err != nil {
        log.Fatalf("请求执行失败: %v", err)
    }
    
    // 处理响应
    if resp.IsSuccess() {
        jsonResp := resp.(*response.CommonJsonResponse)
        data := jsonResp.GetData()
        fmt.Printf("请求成功，响应数据: %v\n", data)
    } else {
        fmt.Printf("请求失败，错误信息: %s\n", resp.GetMessage())
    }
}
```

## 示例

在`example`目录下提供了更多使用示例：

- [用户号码检查示例](example/usernumber_check_example.go) - 演示了如何检查用户号码

## API文档

### 核心组件

- `api.IoTGatewayClient` - 客户端接口，用于发送请求
- `api.IoTGatewayRequest` - 请求接口，定义请求方法
- `api.IoTGatewayResponse` - 响应接口，定义响应方法
- `request.CommonJsonRequest` - 通用JSON请求实现
- `response.CommonJsonResponse` - 通用JSON响应实现

### 工具类

- `utils.SM3` - SM3摘要算法工具类
- `utils.ApaasSignUtil` - 签名工具类
- `utils.WebUtils` - HTTP通信工具类

## 贡献

欢迎贡献代码和提出建议。请确保在提交Pull Request前先运行测试。

## 许可证

本SDK采用MIT许可证。 
