package fourpx

import "testing"

func TestSignExample(t *testing.T) {
	params := map[string]string{
		"app_key":   "16081f05-e8fc-4250-b9c4-0660d1ecbb28",
		"format":    "json",
		"method":    "ds.xms.order.create",
		"timestamp": "1532592413187",
		"v":         "1.0",
	}
	got := Sign(params, `{"aa":"bb"}`, "7eebf328-8e5a-4030-904d-ec6e89174fbc")
	if got != "ff4af77c062a9b97d98aa29777621c4a" {
		t.Fatalf("got %s", got)
	}
}
