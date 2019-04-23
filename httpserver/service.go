package httpserver

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	dtypes "github.com/docker/docker/api/types"
	dcontainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/pocoz/auto-builder/types"
)

type service interface {
	createBuild(ctx context.Context, payload *types.Payload) error
}

type basicService struct {
	logger        log.Logger
	dockerCli     *client.Client
	dockerConfigs []*types.Config
	registryAuth  string
	mu            sync.Mutex
}

// createBuild
func (s *basicService) createBuild(ctx context.Context, payload *types.Payload) error {
	for _, e := range payload.Events {
		if e.Action == types.ActionPush {
			// Validation
			s.mu.Lock()
			if strings.TrimSpace(e.Target.Repository) == "" {
				level.Info(s.logger).Log("err", "empty repository")
				s.mu.Unlock()
				continue
			}
			if strings.TrimSpace(e.Request.Host) == "" {
				level.Info(s.logger).Log("err", "empty request host")
				s.mu.Unlock()
				continue
			}

			// Create image name
			imageName := e.Request.Host + "/" + e.Target.Repository

			// Check for repository content in config and get him envs
			conf, ok := imageInConfig(imageName, s.dockerConfigs)
			if !ok {
				s.mu.Unlock()
				continue
			}

			// Get running containers
			containers, err := s.dockerCli.ContainerList(context.TODO(), dtypes.ContainerListOptions{})
			if err != nil {
				level.Info(s.logger).Log("method", "get containers", "err", err)
				s.mu.Unlock()
				continue
			}

			for _, c := range containers {
				// Find running containers
				if c.Image == imageName {
					// Stop old container
					duration := 5 * time.Second
					err = s.dockerCli.ContainerStop(context.TODO(), c.ID, &duration)
					if err != nil {
						level.Info(s.logger).Log("method", "container stop", "err", err)
						s.mu.Unlock()
						continue
					}

					// Remove old container
					err = s.dockerCli.ContainerRemove(context.TODO(), c.ID, dtypes.ContainerRemoveOptions{})
					if err != nil {
						level.Info(s.logger).Log("method", "container remove", "err", err)
						s.mu.Unlock()
						continue
					}
				}
			}

			// Pull last container
			out, err := s.dockerCli.ImagePull(context.TODO(), imageName, dtypes.ImagePullOptions{RegistryAuth: s.registryAuth})
			if err != nil {
				level.Info(s.logger).Log("method", "container pull", "err", err)
				s.mu.Unlock()
				continue
			}
			_, err = io.Copy(os.Stdout, out)
			if err != nil {
				level.Info(s.logger).Log("method", "[io copy]container pull", "err", err)
				s.mu.Unlock()
				continue
			}
			err = out.Close()
			if err != nil {
				level.Info(s.logger).Log("method", "[out close]container pull", "err", err)
				s.mu.Unlock()
				continue
			}

			volumes := make(map[string]struct{})
			for _, v := range conf.Volumes {
				volumes[v] = struct{}{}
			}

			// Create new container
			newContainer, err := s.dockerCli.ContainerCreate(
				context.TODO(),
				&dcontainer.Config{
					Image:   conf.Image,
					Env:     conf.Environments,
					Cmd:     conf.Cmd,
					Volumes: volumes,
				},
				&dcontainer.HostConfig{
					NetworkMode: "host",
					RestartPolicy: dcontainer.RestartPolicy{
						Name: "always",
					},
					Binds: conf.Volumes,
				}, nil, e.Target.Repository)
			if err != nil {
				// Remove old container
				level.Info(s.logger).Log("method", "container create", "err", err)
				s.mu.Unlock()
				continue
			}

			// Start new container
			err = s.dockerCli.ContainerStart(context.TODO(), newContainer.ID, dtypes.ContainerStartOptions{})
			if err != nil {
				level.Info(s.logger).Log("method", "container start", "err", err)
				s.mu.Unlock()
				continue
			}
			s.mu.Unlock()
		}
	}

	return nil
}

func imageInConfig(image string, conf []*types.Config) (*types.Config, bool) {
	for _, c := range conf {
		if c.Image == image {
			return c, true
		}
	}
	return nil, false
}
