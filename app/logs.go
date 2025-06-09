package app

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DockerLogsService struct {
	DockerBaseService
	cancel context.CancelFunc
	mu     sync.Mutex
}

func NewDockerLogsService() *DockerLogsService {
	return &DockerLogsService{}
}

func StartupDockerLogsService(s *DockerLogsService, ctx context.Context, cli *client.Client) {
	s.ctx = ctx
	s.cli = cli
}

func (s *DockerLogsService) StartWatching(id string) error {
	if s.cli == nil || s.ctx == nil {
		return fmt.Errorf("Docker client not initialized")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Don't start multiple listeners
	if s.cancel != nil {
		return nil
	}

	ctx, cancel := context.WithCancel(s.ctx)
	s.cancel = cancel

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
		Tail:       "10",
	}

	out, err := s.cli.ContainerLogs(ctx, id, options)
	if err != nil {
		return err
	}

	go func() {
		defer out.Close()
		reader := bufio.NewReader(out)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					runtime.EventsEmit(s.ctx, "docker:logs", "ERROR: "+err.Error())
					break
				}
				runtime.EventsEmit(s.ctx, "docker:logs", string(line))
			}
		}
	}()

	return nil
}

func (s *DockerLogsService) StopWatching() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}
