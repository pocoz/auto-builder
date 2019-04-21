package httpserver

import (

	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

type handlerConfig struct {
	svc         service
	keyJWT      []byte
	logger      log.Logger
	rateLimiter *rate.Limiter
}

// newHandler creates a new HTTP handler serving service endpoints.
func newHandler(cfg *handlerConfig) http.Handler {
	svc := &loggingMiddleware{next: cfg.svc, logger: cfg.logger}

	createBuildEndpoint := makeCreateBuildEndpoint(svc)
	createBuildEndpoint = applyMiddleware(createBuildEndpoint, "CreateBuild", cfg)

	router := mux.NewRouter()

	router.Path("/api/v1/build").Methods("POST").Handler(kithttp.NewServer(
		createBuildEndpoint,
		decodeCreateBuildRequest,
		encodeCreateBuildResponse,
	))

	return router
}

func applyMiddleware(e endpoint.Endpoint, method string, cfg *handlerConfig) endpoint.Endpoint {
	return ratelimit.NewErroringLimiter(cfg.rateLimiter)(e)
}
