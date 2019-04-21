package httpserver

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"
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
	logger       log.Logger
	dockerCli    *client.Client
	registryAuth string
}

// createBuild
func (s *basicService) createBuild(ctx context.Context, payload *types.Payload) error {
	// Read config
	dat, err := ioutil.ReadFile("/srv/auto-builder/config.json")
	if err != nil {
		level.Info(s.logger).Log("method", "read config", "err", err)
		return nil
	}

	// Create config
	var config *types.Configs
	err = json.Unmarshal(dat, &config)
	if err != nil {
		level.Info(s.logger).Log("method", "unmarshal config", "err", err)
		return nil
	}

	for _, e := range payload.Events {
		if e.Action == types.ActionPush {
			// Validation
			if strings.TrimSpace(e.Target.Repository) == "" {
				level.Info(s.logger).Log("err", "empty repository")
				continue
			}
			if strings.TrimSpace(e.Request.Host) == "" {
				level.Info(s.logger).Log("err", "empty request host")
				continue
			}

			// Create image name
			imageName := e.Request.Host + "/" + e.Target.Repository

			// Check for repository content in config
			if !imageInConfig(imageName, config) {
				continue
			}

			// docker ps
			containers, err := s.dockerCli.ContainerList(context.Background(), dtypes.ContainerListOptions{})
			if err != nil {
				level.Info(s.logger).Log("method", "get containers", "err", err)
				continue
			}
			for _, c := range containers {
				// Pull last c
				out, err := s.dockerCli.ImagePull(ctx, imageName, dtypes.ImagePullOptions{RegistryAuth: s.registryAuth})
				if err != nil {
					level.Info(s.logger).Log("method", "[image pull]container pull", "err", err)
					continue
				}
				_, err = io.Copy(os.Stdout, out)
				if err != nil {
					level.Info(s.logger).Log("method", "[io copy]container pull", "err", err)
					continue
				}
				err = out.Close()
				if err != nil {
					level.Info(s.logger).Log("method", "[out close]container pull", "err", err)
					continue
				}

				// Find running containers
				if c.Image == imageName {
					// Stop old container
					duration := 15 * time.Second
					err = s.dockerCli.ContainerStop(ctx, c.ID, &duration)
					if err != nil {
						level.Info(s.logger).Log("method", "container stop", "err", err)
						continue
					}

					// Remove old container
					removeOptions := dtypes.ContainerRemoveOptions{
						RemoveVolumes: false,
						RemoveLinks:   false,
						Force:         false,
					}
					err = s.dockerCli.ContainerRemove(ctx, c.ID, removeOptions)
					if err != nil {
						level.Info(s.logger).Log("method", "container remove", "err", err)
						continue
					}
				}

				// Create new container
				newContainer, err := s.dockerCli.ContainerCreate(
					ctx,
					&dcontainer.Config{
						Image: imageName,
					},
					&dcontainer.HostConfig{
						NetworkMode: "host",
					}, nil, e.Target.Repository)
				if err != nil {
					level.Info(s.logger).Log("method", "container create", "err", err)
					continue
				}

				// Start new container
				err = s.dockerCli.ContainerStart(ctx, newContainer.ID, dtypes.ContainerStartOptions{})
				if err != nil {
					level.Info(s.logger).Log("method", "container start", "err", err)
					continue
				}

			}

		}
	}

	return nil
}

func imageInConfig(image string, conf *types.Configs) bool {
	for _, c := range conf.ConfigList {
		if c.Image == image {
			return true
		}
	}
	return false
}
