package yanwen

import "testing"

func TestCheckSuccess(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{"success true", `{"success":true,"code":"0","message":"ok"}`, false},
		{"success true code 200", `{"success":true,"code":"200"}`, false},
		{"success false", `{"success":false,"code":"1001","message":"fail"}`, true},
		{"no success code ok", `{"code":"0"}`, false},
		{"no success code bad", `{"code":"500","message":"x"}`, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var env Response
			_ = env
			err := checkSuccess("m", []byte(tc.raw), &Response{Code: "x", Message: "m"})
			if tc.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected: %v", err)
			}
		})
	}
}
