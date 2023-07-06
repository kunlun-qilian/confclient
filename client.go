package confclient

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Client struct {
	// IP or domain
	Host string `env:""`
	// Port
	Port int `env:""`
	// HTTP HTTPS
	Protocol string `env:""`
	// Second
	Timeout time.Duration
	client  *http.Client
}

func (c *Client) ApiServer() string {
	return fmt.Sprintf("%s://%s:%d", c.Protocol, c.Host, c.Port)
}

func (c *Client) SetDefaults() {
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}

	if c.Protocol == "" {
		c.Protocol = "http"
	}

	if c.Port == 0 {
		c.Port = 80
	}

	if c.Timeout == 0 {
		c.Timeout = 5
	}
}

func (c *Client) Init() {
	c.SetDefaults()
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
