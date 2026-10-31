package methodrpc

import (
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

// RetryConfig controls RoundTripper retries for transient failures.
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
}

// DefaultRetry retries 5xx / network timeouts up to 3 times.
func DefaultRetry() RetryConfig {
	return RetryConfig{MaxAttempts: 3, BaseDelay: 200 * time.Millisecond}
}

// RetryTransport wraps a base RoundTripper with exponential backoff.
type RetryTransport struct {
	Base   http.RoundTripper
	Config RetryConfig
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	cfg := t.Config
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	if cfg.BaseDelay <= 0 {
		cfg.BaseDelay = 200 * time.Millisecond
	}

	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if attempt > 1 && req.Body != nil && req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, err
			}
			req.Body = body
		}

		resp, err := base.RoundTrip(req)
		if err != nil {
			lastErr = err
			if !isTransientNet(err) || attempt == cfg.MaxAttempts {
				return nil, err
			}
			sleep(req, cfg.BaseDelay, attempt)
			continue
		}

		if resp.StatusCode >= 500 && attempt < cfg.MaxAttempts {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			sleep(req, cfg.BaseDelay, attempt)
			continue
		}
		return resp, nil
	}
	return nil, lastErr
}

func sleep(req *http.Request, base time.Duration, attempt int) {
	d := base * time.Duration(1<<(attempt-1))
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-req.Context().Done():
	case <-timer.C:
	}
}

func isTransientNet(err error) bool {
	var ne net.Error
	if errors.As(err, &ne) {
		return ne.Timeout()
	}
	return false
}

// WrapClient installs RetryTransport on hc (or a fresh client).
func WrapClient(hc *http.Client, cfg RetryConfig) *http.Client {
	if hc == nil {
		hc = &http.Client{Timeout: 45 * time.Second}
	}
	clone := *hc
	clone.Transport = &RetryTransport{Base: clone.Transport, Config: cfg}
	return &clone
}
