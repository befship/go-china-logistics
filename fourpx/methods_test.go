package fourpx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/befship/go-china-logistics/methodrpc"
)

// Official request/response samples from 4PX open-platform docs
// (https://open.4px.com/apiInfo/detail) and the XMS Yuque guide
// (https://4pxgroup.yuque.com/uhl5y1/nx79y1/ktwyru).

func testClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	cfg := Config{
		AppKey:      "test-app-key",
		AppSecret:   "test-app-secret",
		Environment: EnvTest,
		HTTPClient:  srv.Client(),
		MaxRetries:  -1, // no retries in unit tests
	}
	return &Client{
		cfg: cfg,
		rpc: methodrpc.New(methodrpc.Config{
			Endpoint:   srv.URL,
			HTTPClient: cfg.HTTPClient,
			MaxRetries: cfg.MaxRetries,
			Platform:   "4px",
		}, Signer{AppSecret: cfg.AppSecret}),
	}
}

func TestCreateOrder_Success(t *testing.T) {
	const respBody = `{
  "data": {
    "ds_consignment_no": "DS800900010807",
    "4px_tracking_no": "800900010807",
    "label_barcode": "",
    "logistics_channel_name": "DHL",
    "ref_no": "TEST-1234567",
    "logistics_channel_no": "900900010237"
  },
  "msg": "System processing succeeded",
  "result": "1"
}`

	var gotMethod, gotBody string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.URL.Query().Get("method")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	})

	// Yuque create sample (simplified to fields supported by the SDK).
	req := yuqueCreateOrderRequest("TEST-1234567")
	data, err := c.CreateOrder(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateOrder: %v", err)
	}
	if gotMethod != MethodOrderCreate {
		t.Fatalf("method=%q", gotMethod)
	}
	if !strings.Contains(gotBody, `"ref_no":"TEST-1234567"`) {
		t.Fatalf("body missing ref_no: %s", gotBody)
	}
	if !strings.Contains(gotBody, `"logistics_product_code":"PY"`) {
		t.Fatalf("body missing product PY: %s", gotBody)
	}
	if data.TrackingNo4PX != "800900010807" || data.DSConsignmentNo != "DS800900010807" {
		t.Fatalf("data=%+v", data)
	}
	if data.LogisticsChannelNo != "900900010237" {
		t.Fatalf("channel=%q", data.LogisticsChannelNo)
	}
}

func TestCreateOrder_BusinessError(t *testing.T) {
	const respBody = `{
	"errors": [{
		"error_code": "DS000007",
		"error_msg": "Ref_no: 736017796559863808 in processing",
		"reference_code": "736017796559863808"
	}],
	"msg": "System processing failed",
	"result": "0"
}`
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(respBody))
	})
	_, err := c.CreateOrder(context.Background(), yuqueCreateOrderRequest("dup"))
	if err == nil {
		t.Fatal("expected business error")
	}
}

func TestGetOrder_Success(t *testing.T) {
	// Official ds.xms.order.get response example + flowchart fields.
	const respBody = `{
	"data": [{
		"consignment_info": {
			"4px_inbound_date": 1611810365000,
			"4px_outbound_date": 1611885462000,
			"4px_tracking_no": "305646756838",
			"consignment_create_date": 1611728635000,
			"consignment_status": "C",
			"ds_consignment_no": "DS303367756838",
			"get_no_mode": "C",
			"get_no_exmsg": "",
			"has_check_oda": "N",
			"insure_status": "N",
			"insure_type": "",
			"is_hold_sign": "N",
			"logistics_channel_no": "420766679365559903506222909294",
			"logistics_product_code": "VN",
			"logistics_product_name": "合伙人产品4经济-带电",
			"oda_result_sign": "N",
			"ref_no": "AQ710216678CN"
		},
		"parcel_confirm_info": {
			"confirm_parcel_qty": "1",
			"confirm_parcel_weight": 80,
			"parcel_list_confirm_info": [{
				"confirm_battery_type": "",
				"confirm_charge_weight": 0,
				"confirm_high": 0,
				"confirm_include_battery": "N",
				"confirm_length": 0,
				"confirm_volume_weight": 0,
				"confirm_weight": 0,
				"confirm_width": 0,
				"currency_code": "USD",
				"parcel_total_value_confirm": 5.01
			}]
		}
	}],
	"msg": "System processing succeeded",
	"result": "1"
}`

	var gotMethod, gotBody string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.URL.Query().Get("method")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		_, _ = w.Write([]byte(respBody))
	})

	items, err := c.GetOrder(context.Background(), &OrderGetRequest{
		RequestNo: "DJSH00718010700002", // official request example
	})
	if err != nil {
		t.Fatalf("GetOrder: %v", err)
	}
	if gotMethod != MethodOrderGet {
		t.Fatalf("method=%q", gotMethod)
	}
	if !strings.Contains(gotBody, `"request_no":"DJSH00718010700002"`) {
		t.Fatalf("body=%s", gotBody)
	}
	if len(items) != 1 || items[0].ConsignmentInfo == nil {
		t.Fatalf("items=%+v", items)
	}
	ci := items[0].ConsignmentInfo
	if ci.GetNoMode != GetNoModeCreate {
		t.Fatalf("get_no_mode=%q", ci.GetNoMode)
	}
	if ci.LogisticsChannelNo == "" {
		t.Fatal("expected logistics_channel_no for label path")
	}
	if ci.TrackingNo4PX != "305646756838" {
		t.Fatalf("tn=%q", ci.TrackingNo4PX)
	}
}

func TestGetOrder_GetNoModeBranches(t *testing.T) {
	// Flowchart: get_no_mode C branches on logistics_channel_no / get_no_exmsg.
	cases := []struct {
		name       string
		mode       string
		channelNo  string
		exMsg      string
		wantLabel  bool
		wantRetry  bool
		wantCancel bool
	}{
		{"mode U -> label", GetNoModeWarehouse, "", "", true, false, false},
		{"mode C channel ready -> label", GetNoModeCreate, "420766679365559903506222909294", "", true, false, false},
		{"mode C pending -> retry get", GetNoModeCreate, "", "", false, true, false},
		{"mode C exmsg -> cancel", GetNoModeCreate, "", "address invalid", false, false, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(map[string]any{
				"result": "1",
				"msg":    "ok",
				"data": []map[string]any{{
					"consignment_info": map[string]any{
						"get_no_mode":          tc.mode,
						"logistics_channel_no": tc.channelNo,
						"get_no_exmsg":         tc.exMsg,
						"4px_tracking_no":      "305646756838",
						"ds_consignment_no":    "DS303367756838",
					},
				}},
			})
			c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(body)
			})
			items, err := c.GetOrder(context.Background(), &OrderGetRequest{RequestNo: "305646756838"})
			if err != nil {
				t.Fatal(err)
			}
			ci := items[0].ConsignmentInfo
			gotLabel, gotRetry, gotCancel := decideLabelPath(ci)
			if gotLabel != tc.wantLabel || gotRetry != tc.wantRetry || gotCancel != tc.wantCancel {
				t.Fatalf("got label=%v retry=%v cancel=%v", gotLabel, gotRetry, gotCancel)
			}
		})
	}
}

// decideLabelPath mirrors the Yuque flowchart after ds.xms.order.get.
func decideLabelPath(ci *ConsignmentInfo) (label, retry, cancel bool) {
	if ci == nil {
		return false, false, false
	}
	if ci.GetNoMode == GetNoModeWarehouse {
		return true, false, false
	}
	if ci.GetNoExMsg != "" {
		return false, false, true
	}
	if ci.LogisticsChannelNo != "" {
		return true, false, false
	}
	return false, true, false
}

func TestCancelOrder_Success(t *testing.T) {
	const respBody = `{
	"result": "1",
	"msg": "System processing succeeded",
	"errors": null,
	"data": null
}`
	var gotMethod, gotBody string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.URL.Query().Get("method")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		_, _ = w.Write([]byte(respBody))
	})
	err := c.CancelOrder(context.Background(), &OrderCancelRequest{
		RequestNo:    "3000210210",
		CancelReason: "test",
	})
	if err != nil {
		t.Fatalf("CancelOrder: %v", err)
	}
	if gotMethod != MethodOrderCancel {
		t.Fatalf("method=%q", gotMethod)
	}
	if !strings.Contains(gotBody, `"request_no":"3000210210"`) || !strings.Contains(gotBody, `"cancel_reason":"test"`) {
		t.Fatalf("body=%s", gotBody)
	}
}

func TestCancelOrder_NotFound(t *testing.T) {
	const respBody = `{
	"result": "0",
	"msg": "System processing failed",
	"errors": [{
		"error_code": "0",
		"error_msg": "Can not find this order",
		"error_data": null,
		"reference_code": "DJSH00717122700004"
	}],
	"data": null
}`
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(respBody))
	})
	err := c.CancelOrder(context.Background(), &OrderCancelRequest{
		RequestNo:    "DJSH00717122700004",
		CancelReason: "test",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetLabel_Success(t *testing.T) {
	const respBody = `{
	"data": {
		"label_barcode": "JL789409098GB",
		"label_url_info": {
			"logistics_label": "http://fpx-xms-oss.oss-cn-shenzhen.aliyuncs.com/fpx-autonomyPrint-lablePdf-f7d66412-3d99-4108-a3bc-ca8874be6f94.pdf",
			"custom_label": "http://fpx-xms-oss.oss-cn-shenzhen.aliyuncs.com/fpx-autonomyPrint-lablePdf-5b53b9d9-8906-4b9e-a868-088e18eb3c94.pdf",
			"package_label": "http://fpx-xms-oss.oss-cn-shenzhen.aliyuncs.com/fpx-autonomyPrint-lablePdf-98e20150-f932-426e-b3ad-be324576896f.pdf"
		}
	},
	"msg": "System processing succeeded",
	"result": "1"
}`
	var gotMethod string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.URL.Query().Get("method")
		_, _ = w.Write([]byte(respBody))
	})
	data, err := c.GetLabel(context.Background(), &LabelGetRequest{
		RequestNo: "YIN201803280000012",
	})
	if err != nil {
		t.Fatalf("GetLabel: %v", err)
	}
	if gotMethod != MethodLabelGet {
		t.Fatalf("method=%q", gotMethod)
	}
	if data.LabelBarcode != "JL789409098GB" {
		t.Fatalf("barcode=%q", data.LabelBarcode)
	}
	if data.LabelURLInfo == nil || data.LabelURLInfo.LogisticsLabel == "" {
		t.Fatalf("label urls=%+v", data.LabelURLInfo)
	}
}

func TestEstimatedCost_Success(t *testing.T) {
	// Live test env returns data as a JSON string; lump_sum_fee as string.
	const respBody = `{
	"data":"[{\"4px_tracking_no\":\"\",\"charge_weight\":\"0.500\",\"estimated_time\":\"5-20\",\"is_show_track\":\"Y\",\"is_volume_cargo\":\"N\",\"logistics_channel_no\":\"\",\"logistics_product_code\":\"PW\",\"lump_sum_fee\":\"54.40\",\"ref_no\":\"\",\"remarks\":\"\"}]",
	"msg": "System processing succeeded",
	"result": "1"
}`
	var gotMethod, gotBody string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.URL.Query().Get("method")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		_, _ = w.Write([]byte(respBody))
	})
	items, err := c.EstimatedCost(context.Background(), &EstimatedCostRequest{
		CountryCode: "US",
		Weight:      "1",
		Length:      "100",
		Width:       "999",
		Height:      "999",
		CargoCode:   CargoCodeParcel,
	})
	if err != nil {
		t.Fatalf("EstimatedCost: %v", err)
	}
	if gotMethod != MethodEstimatedCostGet {
		t.Fatalf("method=%q", gotMethod)
	}
	if !strings.Contains(gotBody, `"country_code":"US"`) || !strings.Contains(gotBody, `"cargocode":"P"`) {
		t.Fatalf("body=%s", gotBody)
	}
	if len(items) != 1 || items[0].LogisticsProductCode != "PW" || items[0].LumpSumFee != 54.40 {
		t.Fatalf("items=%+v", items)
	}
}

func TestListLogisticsProducts_Success(t *testing.T) {
	const respBody = `{
	"data": [{
		"billing_volume_weight": "Y",
		"logistics_product_code": "A2",
		"logistics_product_name_cn": "4PX标准专线",
		"logistics_product_name_en": "4PX  Express",
		"order_track": "Y",
		"surface_mail": "N",
		"transport_mode": "2",
		"with_battery": "N"
	}, {
		"billing_volume_weight": "Y",
		"logistics_product_code": "A5",
		"logistics_product_name_cn": "DHL华东超值大货",
		"logistics_product_name_en": "DHL EastChina Heavy",
		"order_track": "Y",
		"surface_mail": "N"
	}],
	"msg": "System processing succeeded",
	"result": "1"
}`
	var gotMethod, gotBody string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.URL.Query().Get("method")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		_, _ = w.Write([]byte(respBody))
	})
	items, err := c.ListLogisticsProducts(context.Background(), &LogisticsProductGetListRequest{
		TransportMode: TransportAll, // official example: "1"
	})
	if err != nil {
		t.Fatalf("ListLogisticsProducts: %v", err)
	}
	if gotMethod != MethodLogisticsProductGetList {
		t.Fatalf("method=%q", gotMethod)
	}
	if !strings.Contains(gotBody, `"transport_mode":"1"`) {
		t.Fatalf("body=%s", gotBody)
	}
	if len(items) != 2 || items[0].LogisticsProductCode != "A2" {
		t.Fatalf("items=%+v", items)
	}
}

func TestTrack_Success(t *testing.T) {
	const respBody = `{
	"result": "1",
	"msg": "操作成功",
	"errors": [],
	"data": {
		"deliveryOrderNo": "1Z8E26Y00366094077",
		"destinationCountry": "CN",
		"serverName": "",
		"serverNum": "1Z8E26Y00366094077",
		"trackingList": [{
			"businessLinkCode": "AA",
			"occurDatetime": "2015-06-30 11:25:52",
			"occurLocation": "Shenzhen,China",
			"country": "China",
			"city": "Shenzhen",
			"timeZone": "UTC-05:00",
			"trackingContent": "Shipment operation completed"
		},
		{
			"businessLinkCode": "AA",
			"occurDatetime": "2015-06-30 11:02:22",
			"occurLocation": "Shenzhen,China",
			"country": "China",
			"city": "Shenzhen",
			"timeZone": "UTC-05:00",
			"trackingContent": "Shipment arrived at facility"
		}]
	}
}`
	var gotMethod, gotBody string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.URL.Query().Get("method")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		_, _ = w.Write([]byte(respBody))
	})
	data, err := c.Track(context.Background(), "1Z8E26Y00366094077")
	if err != nil {
		t.Fatalf("Track: %v", err)
	}
	if gotMethod != MethodTrackingGet {
		t.Fatalf("method=%q", gotMethod)
	}
	if !strings.Contains(gotBody, `"deliveryOrderNo":"1Z8E26Y00366094077"`) {
		t.Fatalf("body=%s", gotBody)
	}
	if data.DeliveryOrderNo != "1Z8E26Y00366094077" || len(data.TrackingList) != 2 {
		t.Fatalf("data=%+v", data)
	}
}

// yuqueCreateOrderRequest builds the create payload from the Yuque XMS guide sample
// (ref_no / BDS / duty P / product PY / 169g / USD), with required sender/recipient.
func yuqueCreateOrderRequest(refNo string) *OrderCreateRequest {
	return &OrderCreateRequest{
		RefNo:        refNo,
		BusinessType: "BDS",
		DutyType:     DutyDDP, // Yuque sample: "P"
		IOSSNo:       "IM0123456789",
		LogisticsServiceInfo: LogisticsServiceInfo{
			LogisticsProductCode: "PY",
		},
		ReturnInfo: ReturnInfo{
			IsReturnOnDomestic: Yes,
			IsReturnOnOversea:  Yes,
		},
		IsInsure: No,
		ParcelList: []Parcel{{
			Weight:         169,
			ParcelValue:    9,
			Currency:       "USD",
			IncludeBattery: No,
			DeclareProductInfo: []DeclareProductInfo{{
				DeclareProductNameCN:   "手机贴膜",
				DeclareProductNameEN:   "phone film",
				DeclareProductCodeQty:  1,
				DeclareUnitPriceExport: 9,
				CurrencyExport:         "USD",
				DeclareUnitPriceImport: 9,
				CurrencyImport:         "USD",
				BrandExport:            "N/A",
				BrandImport:            "N/A",
			}},
		}},
		Sender: Person{
			FirstName: "ZHANG",
			LastName:  "SAN",
			Phone:     "13800138000",
			PostCode:  "518000",
			Country:   "CN",
			State:     "Guangdong",
			City:      "Shenzhen",
			Street:    "Baoan District",
		},
		RecipientInfo: Person{
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			PostCode:  "90001",
			Country:   "US",
			State:     "CA",
			City:      "Los Angeles",
			Street:    "123 Main St",
		},
	}
}
