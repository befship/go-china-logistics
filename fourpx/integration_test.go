//go:build integration

package fourpx

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

// Live tests against open-test.4px.com using the Yuque XMS guide credentials
// (https://4pxgroup.yuque.com/uhl5y1/nx79y1/ktwyru):
//
// Override with FPX_APP_KEY / FPX_APP_SECRET. Run:
//
//	go test -tags=integration ./fourpx/ -v -count=1
func integrationClient(t *testing.T) *Client {
	t.Helper()
	key := envOr("FPX_APP_KEY", "test")
	secret := envOr("FPX_APP_SECRET", "test3")
	return New(Config{
		AppKey:      key,
		AppSecret:   secret,
		Environment: envOr("FPX_ENV", EnvTest),
		Timeout:     45 * time.Second,
	})
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func TestIntegration_ListLogisticsProducts(t *testing.T) {
	c := integrationClient(t)
	items, err := c.ListLogisticsProducts(context.Background(), &LogisticsProductGetListRequest{
		TransportMode: TransportAll,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("products=%d", len(items))
	if len(items) == 0 {
		t.Fatal("expected at least one logistics product")
	}
	for i, it := range items {
		if i >= 5 {
			break
		}
		t.Logf("  %s %s / %s", it.LogisticsProductCode, it.LogisticsProductNameCN, it.LogisticsProductNameEN)
	}
}

func TestIntegration_EstimatedCost(t *testing.T) {
	c := integrationClient(t)
	// Official estimated_cost example (US / 1g dims / parcel).
	items, err := c.EstimatedCost(context.Background(), &EstimatedCostRequest{
		CountryCode: "US",
		Weight:      "500",
		Length:      "10",
		Width:       "10",
		Height:      "10",
		CargoCode:   CargoCodeParcel,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("quotes=%d", len(items))
	for i, it := range items {
		if i >= 5 {
			break
		}
		t.Logf("  %s fee=%.2f eta=%s", it.LogisticsProductCode, it.LumpSumFee, it.EstimatedTime)
	}
}

func TestIntegration_OrderFlow_CreateGetLabelCancel(t *testing.T) {
	// Yuque flowchart: create → wait ~1s → get → (get_no_mode) → label / cancel.
	if os.Getenv("FPX_ALLOW_CREATE") != "1" {
		t.Skip("set FPX_ALLOW_CREATE=1 to run create/get/label/cancel against test env")
	}
	c := integrationClient(t)
	ctx := context.Background()

	product := envOr("FPX_LOGISTICS_PRODUCT_CODE", "PY") // Yuque sample product
	refNo := fmt.Sprintf("TEST-%d", time.Now().UnixNano())
	req := yuqueCreateOrderRequest(refNo)
	req.LogisticsServiceInfo.LogisticsProductCode = product

	created, err := c.CreateOrder(ctx, req)
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	t.Logf("created 4px=%s ds=%s channel=%s", created.TrackingNo4PX, created.DSConsignmentNo, created.LogisticsChannelNo)

	requestNo := created.TrackingNo4PX
	if requestNo == "" {
		requestNo = created.DSConsignmentNo
	}
	if requestNo == "" {
		requestNo = refNo
	}

	// Doc: wait ~1s before first get.
	time.Sleep(time.Second)

	var channelReady bool
	deadline := time.Now().Add(60 * time.Second)
	for {
		items, err := c.GetOrder(ctx, &OrderGetRequest{RequestNo: requestNo})
		if err != nil {
			t.Fatalf("GetOrder: %v", err)
		}
		if len(items) == 0 || items[0].ConsignmentInfo == nil {
			t.Fatalf("empty GetOrder data: %+v", items)
		}
		ci := items[0].ConsignmentInfo
		t.Logf("get mode=%s channel=%q exmsg=%q status=%s",
			ci.GetNoMode, ci.LogisticsChannelNo, ci.GetNoExMsg, ci.ConsignmentStatus)

		label, retry, cancel := decideLabelPath(ci)
		switch {
		case label:
			channelReady = true
		case cancel:
			if err := c.CancelOrder(ctx, &OrderCancelRequest{
				RequestNo:    requestNo,
				CancelReason: "integration test: get_no_exmsg",
			}); err != nil {
				t.Fatalf("CancelOrder after exmsg: %v", err)
			}
			t.Log("cancelled due to get_no_exmsg; recreate needed per flowchart")
			return
		case retry:
			if time.Now().After(deadline) {
				t.Fatal("logistics_channel_no still empty after 1 minute (flowchart escalation)")
			}
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	if channelReady {
		label, err := c.GetLabel(ctx, &LabelGetRequest{
			RequestNo:           requestNo,
			ResponseLabelFormat: LabelPDF,
		})
		if err != nil {
			t.Logf("GetLabel (best-effort): %v", err)
		} else {
			t.Logf("label barcode=%s urls=%+v", label.LabelBarcode, label.LabelURLInfo)
		}
	}

	if err := c.CancelOrder(ctx, &OrderCancelRequest{
		RequestNo:    requestNo,
		CancelReason: "integration test cleanup",
	}); err != nil {
		t.Logf("CancelOrder cleanup: %v", err)
	}
}

func TestIntegration_Track(t *testing.T) {
	tn := os.Getenv("FPX_TRACKING_NUMBER")
	if tn == "" {
		t.Skip("set FPX_TRACKING_NUMBER to exercise Track")
	}
	c := integrationClient(t)
	data, err := c.Track(context.Background(), tn)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("track=%s events=%d status=%q", data.DeliveryOrderNo, len(data.TrackingList), data.Status)
}
