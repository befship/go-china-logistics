# go-china-logistics

[English](./README.md) | 中文

面向国内跨境货代开放平台（**Method-RPC**：固定入口 URL + `method` 路由 + MD5 签名）的非官方 Go SDK，当前覆盖 **4PX** 与 **燕文**。

> **与 4PX / 燕文无任何隶属或官方背书。** 请自行申请开发者账号。详见 [DISCLAIMER.md](./DISCLAIMER.md)。

## 目录结构

| 路径 | 说明 |
|------|------|
| `methodrpc/` | 通用 RPC 客户端（签名接口、重试、POST、错误） |
| `fourpx/` | 4PX 签名、手写 Model、类型化 API |
| `yanwen/` | 燕文签名、手写 Model、类型化 API（含独立轨迹域名） |
| `examples/` | 最小可运行示例 |

与国际承运商（UPS/DHL/FedEx）的 OAS 自动生成 SDK 不同，国内货代无官方 OpenAPI，因此采用 **手写结构体 + 可插拔 Signer**。

## 安装

```bash
go get github.com/befship/go-china-logistics@latest
```

## 4PX 快速开始

```go
client := fourpx.New(fourpx.Config{
    AppKey:      os.Getenv("FPX_APP_KEY"),
    AppSecret:   os.Getenv("FPX_APP_SECRET"),
    Environment: fourpx.EnvTest,
})
items, err := client.EstimatedCost(ctx, &fourpx.EstimatedCostRequest{
    CountryCode: "US",
    Weight:      "500",
})
```

签名：`MD5(字典序参数拼接 + body + appSecret)`。

## 燕文快速开始

```go
client := yanwen.New(yanwen.Config{
    UserID:      os.Getenv("YANWEN_USER_ID"),
    APIToken:    os.Getenv("YANWEN_API_TOKEN"),
    Environment: yanwen.EnvSandbox,
})
items, err := client.CalcList(ctx, &yanwen.CalcListRequest{
    CityID: "01", CountryID: "US", Weight: 500,
})
```

签名：`MD5(token + user_id + data + format + method + timestamp + version + token)`。

## 许可

MIT — 见 [LICENSE](./LICENSE)。
