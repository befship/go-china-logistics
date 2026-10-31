package methodrpc

// Signer produces a platform-specific signature over query params and request body.
// Each logistics platform (4PX, Yanwen, …) supplies its own implementation.
type Signer interface {
	// Sign returns the lowercase hex signature for params + body.
	// Implementations must ignore the "sign" key if present in params.
	Sign(params map[string]string, body string) string
}

// SignFunc adapts a function to the Signer interface.
type SignFunc func(params map[string]string, body string) string

// Sign implements Signer.
func (f SignFunc) Sign(params map[string]string, body string) string {
	return f(params, body)
}
