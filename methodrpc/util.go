package methodrpc

import (
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"
)

// MD5Hex returns lowercase MD5 hex of s.
func MD5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

// SortedKeyValueConcat concatenates key+value for keys in params in lexicographic
// order, skipping keys listed in skip (e.g. "sign", "access_token", "language").
// This is the common domestic-forwarder signing building block (4PX-style).
func SortedKeyValueConcat(params map[string]string, skip ...string) string {
	skipSet := make(map[string]struct{}, len(skip))
	for _, k := range skip {
		skipSet[k] = struct{}{}
	}
	keys := make([]string, 0, len(params))
	for k := range params {
		if _, ok := skipSet[k]; ok {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString(params[k])
	}
	return b.String()
}
