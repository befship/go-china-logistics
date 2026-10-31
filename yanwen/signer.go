package yanwen

import "github.com/befship/go-china-logistics/methodrpc"

// Signer implements Yanwen open-platform MD5 signing.
//
// Algorithm (official):
//
//	MD5( apiToken + user_id + data + format + method + timestamp + version + apiToken )
//
// lowercase hex. Body JSON is the `data` segment.
type Signer struct {
	APIToken string
	UserID   string
}

// Sign implements methodrpc.Signer.
// params must contain format, method, timestamp, version (and typically user_id).
func (s Signer) Sign(params map[string]string, body string) string {
	return Sign(
		s.APIToken,
		s.UserID,
		body,
		params["format"],
		params["method"],
		params["timestamp"],
		params["version"],
	)
}

// Sign is a package-level helper matching the historical internal SDK signature.
func Sign(apiToken, userID, data, format, method, timestamp, version string) string {
	return methodrpc.MD5Hex(apiToken + userID + data + format + method + timestamp + version + apiToken)
}
