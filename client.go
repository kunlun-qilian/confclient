package confclient

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type RestClient struct {
	// IP or domain
	Host string `env:""`
	// Port
	Port int `env:""`
	// HTTP HTTPS
	Protocol string `env:""`

	host string
}

func (c *RestClient) ApiServer() string {
	return c.host
}

func (c *RestClient) SetDefaults() {
	if c.Protocol == "" {
		c.Protocol = "http"
	}

	if c.Port == 0 {
		c.Port = 80
	}

	c.host = fmt.Sprintf("%s://%s:%d", c.Protocol, c.Host, c.Port)
}

func (c *RestClient) Init() {
	c.SetDefaults()
}

// WithTrace Inject trace
func WithTrace(ctx context.Context, req *http.Request) error {
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))
	return nil
}
