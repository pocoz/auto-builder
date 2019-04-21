package httpserver

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/pocoz/auto-builder/types"
)

func makeCreateBuildEndpoint(svc service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createBuildRequest)
		err := svc.createBuild(ctx, req.Payload)
		return createBuildResponse{Err: err}, nil
	}
}

type createBuildRequest struct {
	Payload *types.Payload
}

type createBuildResponse struct {
	Err error
}
