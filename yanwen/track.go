package yanwen

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/befship/go-china-logistics/methodrpc"
)

// Track queries the public tracking API (separate host; Authorization = user_id).
func (c *Client) Track(ctx context.Context, trackingNumbers string) (*TrackingResponse, error) {
	u, err := url.Parse(TrackBaseURL() + "/api/tracking")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("nums", trackingNumbers)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.cfg.UserID)

	hc := c.rpc.HTTPClient()
	resp, err := hc.Do(req)
	if err != nil {
		return nil, &methodrpc.Error{Platform: "yanwen", Op: "track", Err: err, Message: "http do"}
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &methodrpc.Error{Platform: "yanwen", Op: "track", Err: err, Message: "read body"}
	}
	if resp.StatusCode >= 300 {
		return nil, &methodrpc.Error{
			Platform:   "yanwen",
			Op:         "track",
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("unexpected status"),
			Body:       methodrpc.Truncate(string(raw), 800),
		}
	}

	var out TrackingResponse
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &out); err != nil {
			return nil, &methodrpc.Error{
				Platform: "yanwen",
				Op:       "track",
				Err:      err,
				Message:  "decode response",
				Body:     methodrpc.Truncate(string(raw), 400),
			}
		}
	}
	return &out, nil
}
