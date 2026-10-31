# go-china-logistics

English | [中文](./README.zh-CN.md)

Unofficial Go SDKs for domestic China cross-border forwarders that use a **Method-RPC** style open API (fixed URL + `method` query routing + MD5 signing) — currently **4PX** and **Yanwen**.

> **Not affiliated with 4PX or Yanwen.** You need your own developer credentials. See [DISCLAIMER.md](./DISCLAIMER.md).

## Why this layout

International carriers (UPS / DHL / FedEx) ship OAS-generated SDKs. Domestic forwarders usually have **no official OpenAPI**, so this repo uses a shared Method-RPC client plus hand-written models:

| Path | Purpose |
|------|---------|
| `methodrpc/` | Shared client: sign hook, retries, POST + query routing, errors |
| `fourpx/` | 4PX signer, models, typed methods (create / cancel / label / rate / track) |
| `yanwen/` | Yanwen signer, models, typed methods (+ separate track host) |
| `examples/` | Minimal runnable samples |

## Install

```bash
go get github.com/befship/go-china-logistics@latest
```

## Quick start — 4PX

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/befship/go-china-logistics/fourpx"
)

func main() {
	client := fourpx.New(fourpx.Config{
		AppKey:      os.Getenv("FPX_APP_KEY"),
		AppSecret:   os.Getenv("FPX_APP_SECRET"),
		Environment: fourpx.EnvTest, // or fourpx.EnvProduction
	})

	items, err := client.EstimatedCost(context.Background(), &fourpx.EstimatedCostRequest{
		CountryCode: "US",
		Weight:      "500",
		Length:      "10",
		Width:       "10",
		Height:      "10",
		CargoCode:   fourpx.CargoCodeParcel,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("quotes=%d\n", len(items))
}
```

### 4PX surface

| Method | API |
|--------|-----|
| `CreateOrder` | `ds.xms.order.create` |
| `CancelOrder` | `ds.xms.order.cancel` |
| `GetLabel` | `ds.xms.label.get` |
| `EstimatedCost` | `ds.xms.estimated_cost.get` |
| `ListLogisticsProducts` | `ds.xms.logistics_product.getlist` |
| `Track` | `tr.order.tracking.get` |

Sign: `MD5(sortedQueryKeyValue + bodyJSON + appSecret)` (excludes `sign` / `access_token` / `language`).

## Quick start — Yanwen

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/befship/go-china-logistics/yanwen"
)

func main() {
	client := yanwen.New(yanwen.Config{
		UserID:      os.Getenv("YANWEN_USER_ID"),
		APIToken:    os.Getenv("YANWEN_API_TOKEN"),
		Environment: yanwen.EnvSandbox, // or yanwen.EnvProduction
	})

	items, err := client.CalcList(context.Background(), &yanwen.CalcListRequest{
		CityID:    "01",
		CountryID: "US",
		Weight:    500,
		PostCode:  "90001",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("quotes=%d\n", len(items))
}
```

### Yanwen surface

| Method | API |
|--------|-----|
| `CreateOrder` | `express.order.create` |
| `CancelOrder` | `express.order.cancel` |
| `GetLabel` | `express.order.label.get` |
| `GetOrders` | `express.order.getlist` |
| `CalcList` | `calc.list` |
| `VerifyUSAddress` | `common.verify.us.address` |
| `Track` | `GET api.track.yw56.com.cn/api/tracking` |

Sign: `MD5(apiToken + user_id + data + format + method + timestamp + version + apiToken)`.

## Architecture

```
Carrier adapters (BefShip / your app)
        │
        ▼
 fourpx.Client  /  yanwen.Client   ← typed models + platform Signer
        │
        ▼
   methodrpc.Client                ← HTTP, retries, query+body POST
```

Upper-layer carrier adapters (Rate / Ship / Track / Cancel) stay outside this repo so the SDK can be reused independently.

## License

MIT — see [LICENSE](./LICENSE).
