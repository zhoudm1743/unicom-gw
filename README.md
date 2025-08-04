# IoT网关Go SDK

这是联通格物CMP的IoT网关Go语言SDK，用于与IoT网关服务进行通信。

## 功能特点

- 支持IoT网关的API调用
- 支持JSON和XML格式的请求和响应
- 内置请求签名和身份验证
- 支持请求重试
- 支持文件上传
- 支持HTTPS连接

## 安装

```bash
go get github.com/zhoudm1743/unicom-gw
```

## 使用示例

### 基本用法

```go
package main

import (
	"fmt"
	"time"

	"github.com/zhoudm1743/unicom-gw/api"
	"github.com/zhoudm1743/unicom-gw/api/request"
)

func main() {
	// 创建客户端实例
	client := api.NewIoTGatewayClient(
		"https://api.example.com", // 服务器URL
		"YOUR_APP_ID",             // 应用ID
		"YOUR_APP_SECRET",         // 应用密钥
	)

	// 设置超时时间（可选）
	client.SetConnectTimeout(5 * time.Second)
	client.SetReadTimeout(30 * time.Second)

	// 创建请求
	req := request.NewCommonJsonRequest()
	req.SetApiName("device.query")
	req.SetApiVer("1.0")

	// 设置请求参数
	params := map[string]interface{}{
		"deviceId": "123456789",
		"details":  true,
	}
	req.SetParams(params)

	// 发送请求
	resp, err := client.Execute(req)
	if err != nil {
		fmt.Printf("调用API失败: %v\n", err)
		return
	}

	// 处理响应
	if resp.IsSuccess() {
		fmt.Println("调用API成功!")
		fmt.Printf("响应消息: %s\n", resp.GetMessage())
	} else {
		fmt.Println("调用API失败!")
		fmt.Printf("错误信息: %s\n", resp.GetMessage())
	}
}
```

### 更多示例

更多示例可以在 `example` 目录中找到。

## 目录结构

- `api/` - 核心API接口和实现
  - `internal/` - 内部工具类
    - `utils/` - 工具函数
  - `request/` - 请求类型定义
  - `response/` - 响应类型定义
- `example/` - 使用示例

## 错误处理

SDK中定义了以下两种错误类型：

- `ApiError`: API调用过程中发生的错误
- `ApiRuleError`: 参数验证失败等规则错误

## 许可证

Copyright (c) 2023 