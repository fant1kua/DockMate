package app

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerService struct {
	ctx   context.Context
	cli		*client.Client
}

func NewDockerService() *DockerService {
	return &DockerService{}
}

func StartupDockerService(service *DockerService, ctx context.Context, cli *client.Client) {
	service.ctx = ctx
	service.cli = cli
}

func Test() {}

func (s *DockerService) StartWatchingContainers() (string, error) {
	return "Test", nil
}

func (s *DockerService) StartContainer(containerID string) error {
	if s.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return s.cli.ContainerStart(s.ctx, containerID, container.StartOptions{})
}

func (s *DockerService) StopContainer(containerID string) error {
	if s.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return s.cli.ContainerStop(s.ctx, containerID, container.StopOptions{})
}

func (s *DockerService) RestartContainer(containerID string) error {
	if s.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return s.cli.ContainerRestart(s.ctx, containerID, container.StopOptions{})
}

func (s *DockerService) RemoveContainer(containerID string) error {
	if s.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return s.cli.ContainerRemove(s.ctx, containerID, container.RemoveOptions{})
}

func (s *DockerService) KillContainer(containerID string) error {
	if s.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return s.cli.ContainerKill(s.ctx, containerID, "SIGKILL")
}

func (s *DockerService) ContainerInspect(containerID string) (string, error) {
	if s.cli == nil {
		return "", fmt.Errorf("Docker client not initialized")
	}
	_, raw, err := s.cli.ContainerInspectWithRaw(s.ctx, containerID, false)
	if err != nil {
		return "", fmt.Errorf("failed to get container data: %v", err)
	}

	return string(raw), nil
}