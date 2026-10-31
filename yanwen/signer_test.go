package yanwen

import "testing"

func TestSignExample(t *testing.T) {
	got := Sign("D6140AA383FD8515B09028C586493DDB", "100000", "{}", "json", "common.country.getlist", "1655798809414", "V1.0")
	if got != "bc5f1792e30936c72fa2a35efead8d73" {
		t.Fatalf("got %s", got)
	}
}
