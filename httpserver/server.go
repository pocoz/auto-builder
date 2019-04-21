package httpserver

import (
	"net/http"
	"time"

	"github.com/docker/docker/client"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/time/rate"
)

// ServerHTTP is a service http server.
type ServerHTTP struct {
	logger log.Logger
	srv    *http.Server
}

// Config is a http server configuration.
type Config struct {
	Logger       log.Logger
	Port         string
	DockerCli    *client.Client
	RegistryAuth string
	RateLimiter  *rate.Limiter
}

// New creates a new http server.
func New(cfg *Config) (*ServerHTTP, error) {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	server := &ServerHTTP{
		logger: cfg.Logger,
		srv:    srv,
	}

	svc := &basicService{
		logger:       cfg.Logger,
		dockerCli:    cfg.DockerCli,
		registryAuth: cfg.RegistryAuth,
	}

	handler := newHandler(&handlerConfig{
		svc:         svc,
		logger:      cfg.Logger,
		rateLimiter: cfg.RateLimiter,
	})

	mux.Handle("/api/v1/", accessControl(handler))
	return server, nil
}

// CORS headers
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// Run starts the server.
func (s *ServerHTTP) Run() error {
	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown stopped the http server.
func (s *ServerHTTP) Shutdown() {
	err := s.srv.Close()
	if err != nil {
		level.Info(s.logger).Log("msg", "HTTP server: shutdown has err", "err:", err)
		return
	}
	level.Info(s.logger).Log("msg", "HTTP server: shutdown complete")
}
