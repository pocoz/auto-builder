package httpserver

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/pocoz/auto-builder/types"
)

// loggingMiddleware wraps Service and logs request information to the provided logger.
type loggingMiddleware struct {
	next   service
	logger log.Logger
}

func (m *loggingMiddleware) createBuild(ctx context.Context, payload *types.HookPayload) error {
	begin := time.Now()
	err := m.next.createBuild(ctx, payload)
	level.Info(m.logger).Log(
		"method", "CreatePayload",
		"err", err,
		"elapsed", time.Since(begin),
	)
	return err
}
