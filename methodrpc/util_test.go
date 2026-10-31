package methodrpc

import "testing"

func TestSortedKeyValueConcat(t *testing.T) {
	params := map[string]string{
		"b":    "2",
		"a":    "1",
		"sign": "xxx",
	}
	got := SortedKeyValueConcat(params, "sign")
	if got != "a1b2" {
		t.Fatalf("got %q", got)
	}
}

func TestMD5Hex(t *testing.T) {
	if MD5Hex("") != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Fatal("empty md5 mismatch")
	}
}
