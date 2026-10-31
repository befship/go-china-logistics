// Package yanwen is an unofficial Go SDK for the Yanwen open platform (Method-RPC).
//
// Order APIs: POST {host}/api/order?user_id=&method=&format=&timestamp=&sign=&version=
// Body: business JSON. Sign = MD5(token+user_id+data+format+method+timestamp+version+token).
//
// Tracking uses a separate host (see Track).
//
// Not affiliated with Yanwen. See DISCLAIMER.md at the repository root.
package yanwen

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/befship/go-china-logistics/methodrpc"
)

// Config holds Yanwen credentials and runtime options.
type Config struct {
	UserID      string
	APIToken    string
	Environment string // sandbox|production; default sandbox

	HTTPClient *http.Client
	Timeout    time.Duration
	MaxRetries int
	Logger     methodrpc.Logger
}

// Client talks to Yanwen order + tracking APIs.
type Client struct {
	cfg Config
	rpc *methodrpc.Client
}

// New builds a Yanwen Client. UserID and APIToken are required for Call.
func New(cfg Config) *Client {
	env := normalizeEnv(cfg.Environment)
	cfg.Environment = env
	rpc := methodrpc.New(methodrpc.Config{
		Endpoint:   OrderEndpoint(env),
		HTTPClient: cfg.HTTPClient,
		Timeout:    cfg.Timeout,
		MaxRetries: cfg.MaxRetries,
		Logger:     cfg.Logger,
		Platform:   "yanwen",
	}, Signer{APIToken: cfg.APIToken, UserID: cfg.UserID})
	return &Client{cfg: cfg, rpc: rpc}
}

// Environment returns the normalized environment name.
func (c *Client) Environment() string { return c.cfg.Environment }

// UserID returns the configured user id.
func (c *Client) UserID() string { return c.cfg.UserID }

// OrderEndpoint returns the full order router URL for env.
func OrderEndpoint(env string) string {
	return OrderBaseURL(env) + "/api/order"
}

// OrderBaseURL returns the order API host for env.
func OrderBaseURL(env string) string {
	if normalizeEnv(env) == EnvProduction {
		return "https://open.yw56.com.cn"
	}
	return "https://open-fat.yw56.com.cn"
}

// TrackBaseURL is the public tracking API host.
func TrackBaseURL() string {
	return "http://api.track.yw56.com.cn"
}

func normalizeEnv(env string) string {
	switch strings.ToLower(strings.TrimSpace(env)) {
	case EnvProduction, "prod":
		return EnvProduction
	default:
		return EnvSandbox
	}
}

// Call invokes a Yanwen method with a typed request body and returns the envelope.
func (c *Client) Call(ctx context.Context, method string, req any) (*Response, error) {
	params := map[string]string{
		"user_id":   c.cfg.UserID,
		"method":    method,
		"format":    "json",
		"timestamp": strconv.FormatInt(time.Now().UnixMilli(), 10),
		"version":   APIVersion,
	}

	var envelope Response
	raw, err := c.rpc.Call(ctx, params, req, &envelope)
	if err != nil {
		return nil, err
	}
	if err := checkSuccess(method, raw, &envelope); err != nil {
		return &envelope, err
	}
	return &envelope, nil
}

// checkSuccess mirrors the official envelope rules used in production adapters:
// success=true with code "0"/"200"/empty is OK; explicit success=false is failure;
// when success is omitted, non-ok codes fail.
func checkSuccess(method string, raw []byte, env *Response) error {
	var probe struct {
		Success *bool  `json:"success"`
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	_ = json.Unmarshal(raw, &probe)

	codeOK := probe.Code == "" || probe.Code == "0" || probe.Code == "200"
	if probe.Success != nil && !*probe.Success {
		return businessError(method, env, raw)
	}
	if probe.Success == nil && !codeOK {
		return businessError(method, env, raw)
	}
	return nil
}

func businessError(method string, env *Response, raw []byte) error {
	return &methodrpc.Error{
		Platform: "yanwen",
		Op:       method,
		Code:     env.Code,
		Message:  env.Message,
		Body:     methodrpc.Truncate(string(raw), 800),
	}
}

func decodeData(env *Response, dest any) error {
	if dest == nil || len(env.Data) == 0 || string(env.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(env.Data, dest); err != nil {
		return &methodrpc.Error{Platform: "yanwen", Message: "decode data", Err: err, Body: methodrpc.Truncate(string(env.Data), 400)}
	}
	return nil
}
