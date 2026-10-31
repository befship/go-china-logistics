// Package methodrpc is a shared HTTP client for domestic China logistics
// Method-RPC platforms (fixed URL + method query routing + custom MD5 signing).
package methodrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Logger is an optional structured logger hook.
type Logger interface {
	Debugf(format string, args ...any)
	Errorf(format string, args ...any)
}

// Config configures a Method-RPC client.
type Config struct {
	// Endpoint is the full router URL (including path), e.g.
	// https://open.4px.com/router/api/service
	Endpoint string

	HTTPClient *http.Client
	Timeout    time.Duration // used when HTTPClient is nil; default 45s
	MaxRetries int           // default 2 extra retries on 5xx/timeout (3 attempts total)
	Logger     Logger
	Platform   string // label for errors ("4px", "yanwen", …)
}

// Client performs signed Method-RPC POSTs against a fixed endpoint.
type Client struct {
	cfg    Config
	http   *http.Client
	signer Signer
}

// New builds a Client. signer may be nil if CallWithParams supplies a pre-signed query.
func New(cfg Config, signer Signer) *Client {
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 0
	} else if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 2
	}
	hc := cfg.HTTPClient
	if hc == nil {
		timeout := cfg.Timeout
		if timeout <= 0 {
			timeout = 45 * time.Second
		}
		hc = &http.Client{Timeout: timeout}
	}
	attempts := cfg.MaxRetries + 1
	hc = WrapClient(hc, RetryConfig{MaxAttempts: attempts, BaseDelay: 200 * time.Millisecond})
	return &Client{cfg: cfg, http: hc, signer: signer}
}

// Endpoint returns the configured router URL.
func (c *Client) Endpoint() string { return c.cfg.Endpoint }

// HTTPClient returns the underlying HTTP client.
func (c *Client) HTTPClient() *http.Client { return c.http }

// Call marshals req to JSON, signs query params, POSTs to the endpoint, and
// unmarshals the raw response body into resp (optional).
//
// params must already contain public fields (method, app_key/user_id, timestamp, …)
// but must NOT include "sign" — Call adds it via Signer.
func (c *Client) Call(ctx context.Context, params map[string]string, req any, resp any) ([]byte, error) {
	body := []byte("{}")
	if req != nil {
		raw, err := json.Marshal(req)
		if err != nil {
			return nil, &Error{Platform: c.cfg.Platform, Op: params["method"], Err: err, Message: "marshal request"}
		}
		body = raw
	}
	return c.CallRaw(ctx, params, body, resp)
}

// CallRaw is like Call but accepts a pre-serialized JSON body.
func (c *Client) CallRaw(ctx context.Context, params map[string]string, body []byte, resp any) ([]byte, error) {
	if c.cfg.Endpoint == "" {
		return nil, &Error{Platform: c.cfg.Platform, Message: "endpoint is empty"}
	}
	if body == nil {
		body = []byte("{}")
	}

	q := cloneParams(params)
	if c.signer != nil {
		q["sign"] = c.signer.Sign(q, string(body))
	}

	u, err := url.Parse(c.cfg.Endpoint)
	if err != nil {
		return nil, &Error{Platform: c.cfg.Platform, Op: q["method"], Err: err, Message: "parse endpoint"}
	}
	values := u.Query()
	for k, v := range q {
		values.Set(k, v)
	}
	u.RawQuery = values.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, &Error{Platform: c.cfg.Platform, Op: q["method"], Err: err, Message: "build request"}
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}

	c.logDebugf("%s POST %s body=%s", c.cfg.Platform, truncateURL(u), Truncate(string(body), 400))

	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		c.logErrorf("%s %s: %v", c.cfg.Platform, q["method"], err)
		return nil, &Error{Platform: c.cfg.Platform, Op: q["method"], Err: err, Message: "http do"}
	}
	defer httpResp.Body.Close()

	raw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, &Error{Platform: c.cfg.Platform, Op: q["method"], Err: err, Message: "read body"}
	}

	if httpResp.StatusCode >= 300 {
		return raw, &Error{
			Platform:   c.cfg.Platform,
			Op:         q["method"],
			StatusCode: httpResp.StatusCode,
			Message:    fmt.Sprintf("unexpected status"),
			Body:       Truncate(string(raw), 800),
		}
	}

	if resp != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, resp); err != nil {
			return raw, &Error{
				Platform: c.cfg.Platform,
				Op:       q["method"],
				Err:      err,
				Message:  "decode response",
				Body:     Truncate(string(raw), 400),
			}
		}
	}
	return raw, nil
}

func cloneParams(in map[string]string) map[string]string {
	out := make(map[string]string, len(in)+1)
	for k, v := range in {
		out[k] = v
	}
	return out
}

func truncateURL(u *url.URL) string {
	cp := *u
	q := cp.Query()
	if q.Has("sign") {
		q.Set("sign", "***")
	}
	cp.RawQuery = q.Encode()
	return cp.String()
}

func (c *Client) logDebugf(format string, args ...any) {
	if c.cfg.Logger != nil {
		c.cfg.Logger.Debugf(format, args...)
	}
}

func (c *Client) logErrorf(format string, args ...any) {
	if c.cfg.Logger != nil {
		c.cfg.Logger.Errorf(format, args...)
	}
}
