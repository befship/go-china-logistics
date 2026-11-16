//go:build integration

package yanwen

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/befship/go-china-logistics/methodrpc"
)

// Live tests against open-fat.yw56.com.cn using the open-platform sandbox
// credentials (https://opendocs.yw56.com.cn/webfile/6993833730687111168/):
//
//	go test -tags=integration ./yanwen/ -v -count=1
func integrationClient(t *testing.T) *Client {
	t.Helper()
	return New(Config{
		UserID:      envOr("YANWEN_USER_ID", "100000"),
		APIToken:    envOr("YANWEN_API_TOKEN", "test"),
		Environment: envOr("YANWEN_ENV", EnvSandbox),
		Timeout:     45 * time.Second,
	})
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func TestIntegration_GetChannels(t *testing.T) {
	// Docs: https://opendocs.yw56.com.cn/webfile/6993833730687111168/
	c := integrationClient(t)
	items, err := c.GetChannels(context.Background())
	if err != nil {
		var me *methodrpc.Error
		if errors.As(err, &me) {
			t.Logf("raw body: %s", me.Body)
		}
		t.Fatal(err)
	}
	t.Logf("channels=%d", len(items))
	if len(items) == 0 {
		t.Fatal("expected at least one activated channel/product")
	}
	for i, it := range items {
		if i >= 5 {
			break
		}
		t.Logf("  id=%s nameCh=%q nameEn=%q", it.ID, it.NameCh, it.NameEn)
	}
}

func TestIntegration_CalcList(t *testing.T) {
	// calc.list often needs sales activation; official sandbox may return product-config errors.
	c := integrationClient(t)
	items, err := c.CalcList(context.Background(), &CalcListRequest{
		CityID:    "01",
		CountryID: "US",
		Weight:    500,
		Length:    10,
		Width:     10,
		High:      10,
		PostCode:  "90001",
	})
	if err != nil {
		t.Logf("CalcList (sandbox may lack rate config): %v", err)
		return
	}
	t.Logf("quotes=%d", len(items))
	for i, it := range items {
		if i >= 5 {
			break
		}
		t.Logf("  %s %s total=%.2f %s", it.ProductNumber, it.ProductName, it.TotalMoney, it.Currency)
	}
}

func TestIntegration_VerifyUSAddress(t *testing.T) {
	c := integrationClient(t)
	data, err := c.VerifyUSAddress(context.Background(), &VerifyUSAddressRequest{
		ReceiverInfo: USAddressFields{
			Address: "1600 Amphitheatre Parkway",
			City:    "Mountain View",
			State:   "CA",
			ZipCode: "94043",
		},
	})
	if err != nil {
		// Official sandbox intermittently returns temporary service errors.
		t.Logf("VerifyUSAddress (best-effort): %v", err)
		return
	}
	t.Logf("verified address=%+v", data.ReceiverInfo)
}

func TestIntegration_OrderFlow_CreateGetLabelCancel(t *testing.T) {
	// Flow: channel.getlist → create → getlist → label → cancel.
	if os.Getenv("YANWEN_ALLOW_CREATE") != "1" {
		t.Skip("set YANWEN_ALLOW_CREATE=1 to run create/get/label/cancel against sandbox")
	}
	c := integrationClient(t)
	ctx := context.Background()

	channelID := os.Getenv("YANWEN_CHANNEL_ID")
	t.Logf("channelID=%s", channelID)
	if channelID == "" {
		channels, err := c.GetChannels(ctx)
		if err != nil {
			t.Fatalf("GetChannels: %v", err)
		}
		if len(channels) == 0 {
			t.Fatal("no activated channels to create against")
		}
		channelID = channels[0].ID
		t.Logf("using channel %s (%s)", channelID, channels[0].NameCh)
	}

	orderNo := fmt.Sprintf("TEST-%d", time.Now().UnixNano())
	req := sandboxCreateOrderRequest(orderNo, channelID)
	created, err := c.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	t.Logf("created waybill=%s order=%s yanwen=%s",
		created.WaybillNumber, created.OrderNumber, created.YanwenOrderNumber)

	lookup := created.WaybillNumber
	if lookup == "" {
		lookup = created.OrderNumber
	}
	if lookup == "" {
		lookup = orderNo
	}

	time.Sleep(time.Second)

	orders, err := c.GetOrders(ctx, []string{lookup})
	if err != nil {
		t.Fatalf("GetOrders: %v", err)
	}
	t.Logf("orders=%d", len(orders))
	for _, o := range orders {
		t.Logf("  waybill=%s order=%s channel=%s status=%d",
			o.WaybillNumber, o.OrderNumber, o.ChannelID, o.Status)
	}

	waybill := created.WaybillNumber
	if waybill == "" && len(orders) > 0 {
		waybill = orders[0].WaybillNumber
	}
	if waybill != "" {
		label, err := c.GetLabel(ctx, &LabelGetRequest{WaybillNumber: waybill})
		if err != nil {
			t.Logf("GetLabel (best-effort): %v", err)
		} else {
			t.Logf("label waybill=%s success=%v b64len=%d err=%q",
				label.WaybillNumber, label.IsSuccess, len(label.Base64String), label.ErrorMsg)
		}

		if err := c.CancelOrder(ctx, &OrderCancelRequest{
			WaybillNumber: waybill,
			Note:          "integration test cleanup",
		}); err != nil {
			t.Logf("CancelOrder cleanup: %v", err)
		}
	}
}

func TestIntegration_Track(t *testing.T) {
	tn := os.Getenv("YANWEN_TRACKING_NUMBER")
	if tn == "" {
		t.Skip("set YANWEN_TRACKING_NUMBER to exercise Track")
	}
	c := integrationClient(t)
	data, err := c.Track(context.Background(), tn)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("track rows=%d", len(data.Data))
	for i, sh := range data.Data {
		if i >= 3 {
			break
		}
		t.Logf("  status=%q exchange=%q pod=%q", sh.Status, sh.ExchangeNumber, sh.POD)
	}
}

// sandboxCreateOrderRequest builds a minimal US create payload using doc fields
// (https://opendocs.yw56.com.cn/webfile/6993833835662151680/).
func sandboxCreateOrderRequest(orderNumber, channelID string) *OrderCreateRequest {
	return &OrderCreateRequest{
		ChannelID:   channelID,
		OrderSource: "api",
		OrderNumber: orderNumber,
		CompanyCode: envOr("YANWEN_COMPANY_CODE", CompanyLAX01),
		ReceiverInfo: AddressInfo{
			Name:    "Test Receiver",
			Phone:   "6265550100",
			Country: "US",
			State:   "CA",
			City:    "Los Angeles",
			ZipCode: "90001",
			Address: "123 Test Street",
		},
		ParcelInfo: ParcelInfo{
			HasBattery:    0,
			Currency:      "USD",
			TotalQuantity: 1,
			TotalWeight:   200,
			Length:        10,
			Width:         10,
			Height:        10,
			ProductList: []ProductItem{{
				GoodsNameCh: "杯子",
				GoodsNameEn: "Cup",
				Price:       12.5,
				Quantity:    1,
				Weight:      200,
			}},
		},
	}
}
