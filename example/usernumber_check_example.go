package main

import (
	"fmt"
	"log"

	"github.com/zhoudm1743/unicom-gw/api"
	"github.com/zhoudm1743/unicom-gw/api/request"
	"github.com/zhoudm1743/unicom-gw/api/response"
)

func main() {
	// 创建客户端实例
	client := api.NewIoTGatewayClient("https://gwapi.10646.cn/api/", "azjJ2DKEhQ", "c8P00UiFwiBqHHWxU1iSV1ok0iCwUY", "46027ouagaXQO3R")

	// 设置超时时间（可选）
	client.SetConnectTimeout(2000) // 连接超时时长设置（单位：毫秒），可以不设置，默认两秒
	// client.SetReadTimeout(30000) // 读取超时时长设置（单位：毫秒），可以不设置，默认三十秒
	// client.SetRetryCount(1) // 连接超时后，重试次数 0：不重试 1：重试一次 以此类推。可以不设置，默认不重试

	// 创建请求
	req := request.NewCommonJsonRequest()

	// 完整的url：https://gwtest.10646.cn/api/GetAccountIdByAcctName_V1_0Main/vV1.0
	req.SetApiName("wsGetTerminalDetails/V1/1Main")
	req.SetApiVer("V1.1")

	// 业务参数，需要发给能力提供者
	params := map[string]interface{}{
		"messageId": "1",
		"version":   "V1.1",
		"openId":    "46027ouagaXQO3R",
		"iccids":    []string{"89860625680009634556"},
	}
	req.SetParams(params)

	// 执行请求
	resp, err := client.Execute(req)
	if err != nil {
		log.Fatalf("请求执行失败: %v", err)
	}

	// 输出请求信息
	fmt.Printf("发送请求流水号: %s\n", req.GetTransId())
	fmt.Printf("请求报文: %s\n", req.GetReqText())
	fmt.Printf("请求是否被成功处理: %v\n", resp.IsSuccess())

	// 获取响应数据
	if jsonResp, ok := resp.(*response.CommonJsonResponse); ok {
		data := jsonResp.GetData()

		if resp.IsSuccess() {
			fmt.Printf("成功返回业务参数: %v\n", data)
		} else {
			fmt.Printf("处理失败返回消息: %v\n", data)
		}
	}
}
