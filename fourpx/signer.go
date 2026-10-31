package fourpx

import "github.com/befship/go-china-logistics/methodrpc"

// Signer implements 4PX open-platform MD5 signing.
//
// Algorithm (official):
//
//	MD5( sortedKeyValue(params excluding sign/access_token/language) + bodyJSON + appSecret )
//
// lowercase hex.
type Signer struct {
	AppSecret string
}

// Sign implements methodrpc.Signer.
func (s Signer) Sign(params map[string]string, body string) string {
	base := methodrpc.SortedKeyValueConcat(params, "sign", "access_token", "language")
	return methodrpc.MD5Hex(base + body + s.AppSecret)
}

// Sign is a package-level helper matching the historical internal SDK signature.
func Sign(params map[string]string, body, secret string) string {
	return Signer{AppSecret: secret}.Sign(params, body)
}
