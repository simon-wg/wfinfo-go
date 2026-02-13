package wfm

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultBaseURL = "https://api.warframe.market"
	DefaultTimeout = 30 * time.Second
)

// Client is a Warframe Market API client.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	ctx        context.Context
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithBaseURL sets the base URL for the API.
func WithBaseURL(baseURL *url.URL) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets the HTTP client for the API.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// NewClient creates a new Warframe Market API client.
func NewClient(opts ...ClientOption) *Client {
	u, _ := url.Parse(DefaultBaseURL)
	c := &Client{
		baseURL: u,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// WithContext returns a shallow copy of the client with the provided context.
func (c *Client) WithContext(ctx context.Context) *Client {
	if ctx == nil {
		return c
	}
	c2 := *c
	c2.ctx = ctx
	return &c2
}

func (c *Client) context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (c *Client) do(req *http.Request, v any) error {
	ctx := req.Context()
	var resp *http.Response
	var err error

	failureCount := 0
	for {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusTooManyRequests && failureCount < 5 {
			resp.Body.Close()
			failureCount++
			sleepDuration := time.Duration(math.Pow(2, float64(failureCount))) * time.Second

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(sleepDuration):
				continue
			}
		}
		break
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// FetchItems fetches items from the Warframe Market API.
func (c *Client) FetchItems() ([]Item, error) {
	u := c.baseURL.JoinPath("v2", "items")
	req, err := http.NewRequestWithContext(c.context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var resp itemsResponse
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	return resp.Data, nil
}
