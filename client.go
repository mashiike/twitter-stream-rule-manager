package twstrulemgr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Endpoint    string
	BearerToken string
	HTTPClient  *http.Client
}

func newClient(config *Config) *Client {
	u := url.URL{
		Scheme: "https",
		Host:   config.Endpoint,
		Path:   "/2/tweets/search/stream/rules",
	}
	return &Client{
		Endpoint:    u.String(),
		BearerToken: config.BearerToken,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

type rulesResponse struct {
	Data Rules `json:"data,omitempty"`
	Meta struct {
		Sent    time.Time       `json:"sent,omitempty"`
		Summary json.RawMessage `json:"summary,omitempty"`
	} `json:"meta,omitempty"`
}

func (c *Client) GetRules(ctx context.Context) (Rules, time.Time, error) {
	req, err := c.newRequestWithContext(ctx, http.MethodGet, c.Endpoint, nil)
	if err != nil {
		return nil, time.Time{}, err
	}
	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, time.Time{}, err
	}
	return resp.Data, resp.Meta.Sent, nil
}

func (c *Client) PostRules(ctx context.Context, addRules Rules, deleteIDs []string, dryRun bool) error {
	endpoint := c.Endpoint
	if dryRun {
		endpoint += "?dry_run=true"
	}
	bodyMap := make(map[string]interface{}, 2)
	if len(addRules) > 0 && len(deleteIDs) == 0 {
		bodyMap["add"] = addRules
	}
	if len(addRules) == 0 && len(deleteIDs) > 0 {
		bodyMap["delete"] = map[string]interface{}{
			"ids": deleteIDs,
		}
	}
	if len(bodyMap) == 0 {
		return nil
	}
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return err
	}
	req, err := c.newRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("%s\n", string(resp.Meta.Summary))
	return nil

}

func (c *Client) newRequestWithContext(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.BearerToken))
	return req, nil
}

func (c *Client) doRequest(ctx context.Context, req *http.Request) (*rulesResponse, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: can not do: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("request failed: can not read error response: %w", err)
		}
		return nil, fmt.Errorf("request failed: HTTP Status=%d, body=%s", resp.StatusCode, string(body))
	}
	decoder := json.NewDecoder(resp.Body)
	var jsonResp rulesResponse
	if err := decoder.Decode(&jsonResp); err != nil {
		return nil, fmt.Errorf("status success, but resp decode error: %w", err)
	}
	return &jsonResp, nil
}
