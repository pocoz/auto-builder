package httpserver

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"

	"github.com/pocoz/auto-builder/types"
)

type service interface {
	createBuild(ctx context.Context, payload *types.HookPayload) error
}

type basicService struct {
	logger  log.Logger
}

// createBuild
func (s *basicService) createBuild(ctx context.Context, payload *types.HookPayload) error {
	fmt.Println(payload)
	return nil
}
