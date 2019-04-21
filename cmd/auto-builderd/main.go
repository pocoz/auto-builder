package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/pocoz/auto-builder/types"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/time/rate"

	"github.com/pocoz/auto-builder/httpserver"
)

type configuration struct {
	HTTPPort         string        `envconfig:"BUILDER_HTTP_PORT"         default:"24001"`
	RateLimitEvery   time.Duration `envconfig:"BUILDER_RATE_LIMIT_EVERY"  default:"1us"`
	RateLimitBurst   int           `envconfig:"BUILDER_RATE_LIMIT_BURST"  default:"100"`
	RegistryLogin    string        `envconfig:"BUILDER_REGISTRY_LOGIN"`
	RegistryPassword string        `envconfig:"BUILDER_REGISTRY_PASSWORD"`
}

func main() {
	const (
		exitCodeSuccess = 0
		exitCodeFailure = 1
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	var cfg configuration
	if err := envconfig.Process("", &cfg); err != nil {
		level.Error(logger).Log("msg", "failed to load configuration", "err", err)
		os.Exit(exitCodeFailure)
	}

	dockerCli, err := client.NewEnvClient()
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize docker client", "err", err)
		os.Exit(exitCodeFailure)
	}

	registryAuthConfig := dtypes.AuthConfig{
		Username: cfg.RegistryLogin,
		Password: cfg.RegistryPassword,
	}
	encodedRegistryAuthConfig, err := json.Marshal(registryAuthConfig)
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize registry auth config", "err", err)
		os.Exit(exitCodeFailure)
	}
	registryAuth := base64.URLEncoding.EncodeToString(encodedRegistryAuthConfig)

	// Read config
	dat, err := ioutil.ReadFile("/srv/auto-builder/config.json")
	if err != nil {
		level.Error(logger).Log("msg", "failed to read config", "err", err)
		os.Exit(exitCodeFailure)
	}

	// Create config
	var configs *types.Configs
	err = json.Unmarshal(dat, &configs)
	if err != nil {
		level.Error(logger).Log("msg", "failed to create config", "err", err)
		os.Exit(exitCodeFailure)
	}

	serverHTTP, err := httpserver.New(&httpserver.Config{
		Logger:        logger,
		Port:          cfg.HTTPPort,
		DockerCli:     dockerCli,
		DockerConfigs: configs,
		RegistryAuth:  registryAuth,
		RateLimiter:   rate.NewLimiter(rate.Every(cfg.RateLimitEvery), cfg.RateLimitBurst),
	})
	if err != nil {
		level.Error(logger).Log("msg", "failed to initialize HTTP server", "err", err)
		os.Exit(exitCodeFailure)
	}
	go func() {
		level.Info(logger).Log("msg", "starting HTTP server", "port", cfg.HTTPPort)
		if err := serverHTTP.Run(); err != nil {
			level.Error(logger).Log("msg", "HTTP server run failure", "err", err)
			os.Exit(exitCodeFailure)
		}
	}()

	errc := make(chan error, 1)
	donec := make(chan struct{})
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, os.Interrupt)
	defer func() {
		signal.Stop(sigc)
		cancel()
	}()

	go func() {
		select {
		case sig := <-sigc:
			level.Info(logger).Log("msg", "received signal, exiting", "signal", sig)
			serverHTTP.Shutdown() // Shutdown HTTP server
			signal.Stop(sigc)
			close(donec)
		case <-errc:
			level.Info(logger).Log("msg", "now exiting with error", "error code", exitCodeFailure)
			os.Exit(exitCodeFailure)
		}
	}()

	<-donec
	level.Info(logger).Log("msg", "goodbye")
	os.Exit(exitCodeSuccess)
}
