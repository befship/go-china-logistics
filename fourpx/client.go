// Package fourpx is an unofficial Go SDK for the 4PX open platform (Method-RPC).
//
// Transport: POST {host}/router/api/service?method=&app_key=&v=&timestamp=&format=&sign=
// Body: business JSON. Sign = MD5(sortedQuery + body + secret).
//
// Not affiliated with 4PX. See DISCLAIMER.md at the repository root.
package fourpx

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/befship/go-china-logistics/methodrpc"
)

// Config holds 4PX credentials and runtime options.
type Config struct {
	AppKey      string
	AppSecret   string
	Environment string // sandbox|test|production; default sandbox

	HTTPClient *http.Client
	Timeout    time.Duration
	MaxRetries int
	Logger     methodrpc.Logger

	Language    string // optional cn|en
	AccessToken string // optional OAuth / ISV
}

// Client talks to the 4PX open-platform router.
type Client struct {
	cfg Config
	rpc *methodrpc.Client
}

// New builds a 4PX Client. AppKey and AppSecret are required for Call.
func New(cfg Config) *Client {
	env := normalizeEnv(cfg.Environment)
	cfg.Environment = env
	rpc := methodrpc.New(methodrpc.Config{
		Endpoint:   Endpoint(env),
		HTTPClient: cfg.HTTPClient,
		Timeout:    cfg.Timeout,
		MaxRetries: cfg.MaxRetries,
		Logger:     cfg.Logger,
		Platform:   "4px",
	}, Signer{AppSecret: cfg.AppSecret})
	return &Client{cfg: cfg, rpc: rpc}
}

// Environment returns the normalized environment name.
func (c *Client) Environment() string { return c.cfg.Environment }

// Endpoint returns the full router URL for env.
func Endpoint(env string) string {
	return BaseURL(env) + "/router/api/service"
}

// BaseURL returns the API host for env (without path).
func BaseURL(env string) string {
	switch normalizeEnv(env) {
	case EnvProduction:
		return "https://open.4px.com"
	case EnvSandbox:
		// Legacy host; official docs now list open-test.4px.com
		return "https://open.sandbox.4px.com"
	default:
		return "https://open-test.4px.com"
	}
}

func normalizeEnv(env string) string {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case EnvProduction, "prod":
		return EnvProduction
	case EnvSandbox:
		return EnvSandbox
	case "", EnvTest, "fat", "dev":
		return EnvTest
	default:
		return strings.ToLower(strings.TrimSpace(env))
	}
}

// Call invokes a 4PX method with a typed request body and returns the envelope.
func (c *Client) Call(ctx context.Context, method string, req any) (*Response, error) {
	params := map[string]string{
		"method":    method,
		"app_key":   c.cfg.AppKey,
		"v":         APIVersion,
		"timestamp": strconv.FormatInt(time.Now().UnixMilli(), 10),
		"format":    "json",
	}
	if c.cfg.Language != "" {
		params["language"] = c.cfg.Language
	}
	if c.cfg.AccessToken != "" {
		params["access_token"] = c.cfg.AccessToken
	}

	var envelope Response
	raw, err := c.rpc.Call(ctx, params, req, &envelope)
	if err != nil {
		return nil, err
	}
	if err := checkResult(method, &envelope, raw); err != nil {
		return &envelope, err
	}
	return &envelope, nil
}

func checkResult(method string, env *Response, raw []byte) error {
	// Official: result "1" success; "0" failure. Some APIs also treat "2" as ok.
	// Empty result is accepted for a few no-payload successes.
	switch env.Result {
	case "1", "2", "":
		return nil
	default:
		return businessError(method, env, raw)
	}
}

func businessError(method string, env *Response, raw []byte) error {
	code := env.Result
	msg := env.Msg
	if msg == "" && len(env.Errors) > 0 {
		code = env.Errors[0].ErrorCode
		msg = env.Errors[0].ErrorMsg
	}
	return &methodrpc.Error{
		Platform: "4px",
		Op:       method,
		Code:     code,
		Message:  msg,
		Body:     methodrpc.Truncate(string(raw), 800),
	}
}

// decodeData unmarshals envelope.Data into dest.
func decodeData(env *Response, dest any) error {
	if dest == nil || len(env.Data) == 0 || string(env.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(env.Data, dest); err != nil {
		return &methodrpc.Error{Platform: "4px", Message: "decode data", Err: err, Body: methodrpc.Truncate(string(env.Data), 400)}
	}
	return nil
}
