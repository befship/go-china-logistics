package methodrpc

import (
	"fmt"
	"strings"
)

// Error is a normalized Method-RPC / transport failure.
type Error struct {
	Platform   string // e.g. "4px", "yanwen"
	Op         string // method name or operation
	StatusCode int    // HTTP status when relevant
	Code       string // platform business code
	Message    string
	Body       string // truncated raw body for diagnostics
	Err        error
}

func (e *Error) Error() string {
	var b strings.Builder
	if e.Platform != "" {
		b.WriteString(e.Platform)
	} else {
		b.WriteString("methodrpc")
	}
	if e.Op != "" {
		b.WriteString(": ")
		b.WriteString(e.Op)
	}
	if e.StatusCode != 0 {
		fmt.Fprintf(&b, " HTTP %d", e.StatusCode)
	}
	if e.Code != "" {
		fmt.Fprintf(&b, " code=%s", e.Code)
	}
	if e.Message != "" {
		fmt.Fprintf(&b, ": %s", e.Message)
	}
	if e.Err != nil {
		fmt.Fprintf(&b, ": %v", e.Err)
	}
	if e.Body != "" && e.Message == "" {
		fmt.Fprintf(&b, ": %s", Truncate(e.Body, 512))
	}
	return b.String()
}

func (e *Error) Unwrap() error { return e.Err }

// Truncate shortens s for error messages.
func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
