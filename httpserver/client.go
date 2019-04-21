package httpserver

import (
	"context"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/pocoz/auto-builder/types"

)

// Client is a client for service.
type Client struct {
	createBuild endpoint.Endpoint
}

// NewClient creates a new service client.
func NewClient(serviceURL string) (*Client, error) {
	baseURL, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		createBuild: kithttp.NewClient(
			"POST",
			baseURL,
			encodeCreateBuildRequest,
			decodeCreateBuildResponse,
		).Endpoint(),
	}

	return c, nil
}

// CreateBuild
func (c *Client) CreateBuild(ctx context.Context, payload *types.HookPayload) error {
	request := createBuildRequest{Payload: payload}
	response, err := c.createBuild(ctx, request)
	if err != nil {
		return err
	}
	res := response.(createBuildResponse)
	return res.Err
}
