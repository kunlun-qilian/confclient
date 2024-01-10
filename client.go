package confclient

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Client struct {
	// IP or domain
	Host string `env:""`
	// Second
	Timeout time.Duration
	client  *http.Client
}

func NewClient(host string) *Client {
	return &Client{
		Host: host,
	}
}

func (c *Client) ApiServer() string {
	return c.Host
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.Timeout = timeout
	return c
}

func (c *Client) WithHttpClient(httpClient *http.Client) *Client {
	c.client = httpClient
	return c
}

func (c *Client) SetDefaults() {
	if c.Host == "" {
		c.Host = "http://127.0.0.1"
	}

	if c.Timeout == 0 {
		c.Timeout = 5
	}
}

func (c *Client) Init() {
	c.SetDefaults()
	if c.client == nil {
		c.client = &http.Client{
			Timeout: c.Timeout * time.Second,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WithTrace Inject trace
func WithTrace(ctx context.Context, req *http.Request) error {
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	return nil
}
