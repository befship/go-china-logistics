package yanwen

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/befship/go-china-logistics/methodrpc"
)

func testClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	cfg := Config{
		UserID:      "100000",
		APIToken:    "test-token",
		Environment: EnvSandbox,
		HTTPClient:  srv.Client(),
		MaxRetries:  -1,
	}
	return &Client{
		cfg: cfg,
		rpc: methodrpc.New(methodrpc.Config{
			Endpoint:   srv.URL,
			HTTPClient: cfg.HTTPClient,
			MaxRetries: cfg.MaxRetries,
			Platform:   "yanwen",
		}, Signer{APIToken: cfg.APIToken, UserID: cfg.UserID}),
	}
}

func TestGetChannels_NumericErrorCode(t *testing.T) {
	// Live error bodies use numeric code (e.g. sign failure), unlike doc samples.
	const respBody = `{"code":500,"message":"验证服务-签名错误,可能秘钥已失效,请获取新秘钥对接","success":false}`
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	})
	_, err := c.GetChannels(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	var me *methodrpc.Error
	if !errors.As(err, &me) {
		t.Fatalf("want methodrpc.Error, got %T: %v", err, err)
	}
	if me.Code != "500" {
		t.Fatalf("code=%q want 500", me.Code)
	}
}

func TestGetChannels_Success(t *testing.T) {
	// Sample from https://opendocs.yw56.com.cn/webfile/6993833730687111168/
	const respBody = `{
  "success": true,
  "code": "0",
  "message": "操作成功",
  "data": [
    {
      "nameCh": "燕文专线追踪-普货",
      "nameEn": "Direct Line Tracked Packet-P",
      "id": "481"
    },
    {
      "nameCh": "燕文专线追踪-特货",
      "nameEn": "",
      "id": "484"
    }
  ]
}`

	var gotMethod, gotBody string
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.URL.Query().Get("method")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(respBody))
	})

	items, err := c.GetChannels(context.Background())
	if err != nil {
		t.Fatalf("GetChannels: %v", err)
	}
	if gotMethod != MethodChannelGetList {
		t.Fatalf("method=%q", gotMethod)
	}
	if gotBody != "{}" {
		t.Fatalf("body=%q want {}", gotBody)
	}
	if len(items) != 2 {
		t.Fatalf("len=%d", len(items))
	}
	if items[0].ID != "481" || items[0].NameCh != "燕文专线追踪-普货" {
		t.Fatalf("item0=%+v", items[0])
	}
	if items[1].ID != "484" {
		t.Fatalf("item1=%+v", items[1])
	}
}
