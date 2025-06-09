package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DockerVolumesService struct {
	DockerBaseService
	cancel context.CancelFunc
	mu     sync.Mutex
}

type VolumeInfo struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Size      int64    `json:"size"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
}

func NewDockerVolumesService() *DockerVolumesService {
	return &DockerVolumesService{}
}

func StartupDockerVolumesService(s *DockerVolumesService, ctx context.Context, cli *client.Client) {
	s.ctx = ctx
	s.cli = cli
}

func (s *DockerVolumesService) List() ([]VolumeInfo, error) {
	if s.cli == nil || s.ctx == nil {
		return nil, fmt.Errorf("Docker client not initialized")
	}
	list, err := s.cli.VolumeList(s.ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %v", err)
	}

	result := s.formatList(list)
	return result, nil
}

func (s *DockerVolumesService) StartWatching() error {
	if s.cli == nil || s.ctx == nil {
		return fmt.Errorf("Docker client not initialized")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Don't start multiple listeners
	if s.cancel != nil {
		return nil
	}

	s.sendListUpdate()

	ctx, cancel := context.WithCancel(s.ctx)
	s.cancel = cancel

	go func() {
		eventFilter := filters.NewArgs()
		eventFilter.Add("type", "volume")
		eventsChan, errs := s.cli.Events(ctx, events.ListOptions{
			Filters: eventFilter,
		})

		for {
			select {
			case event := <-eventsChan:
				if event.Type == events.VolumeEventType {
					s.sendListUpdate()
				}
			case err := <-errs:
				if err != nil {
					time.Sleep(2 * time.Second)
					s.StopWatching() // stop on error
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (s *DockerVolumesService) StopWatching() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}

func (s *DockerVolumesService) formatList(volumes volume.ListResponse) []VolumeInfo {
	list := make([]VolumeInfo, 0, len(volumes.Volumes))
	for _, vol := range volumes.Volumes {
		list = append(list, VolumeInfo{
			Name:      vol.Name,
			CreatedAt: vol.CreatedAt,
		})
	}
	return list
}

func (s *DockerVolumesService) sendListUpdate() {
	list, err := s.cli.VolumeList(s.ctx, volume.ListOptions{})
	if err != nil {
		return
	}

	result := s.formatList(list)
	runtime.EventsEmit(s.ctx, "docker:volumes", result)
}

func (s *DockerVolumesService) Remove(id string) error {
	if s.cli == nil || s.ctx == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	err := s.cli.VolumeRemove(s.ctx, id, true)
	return err
}

func (s *DockerVolumesService) Inspect(id string) (string, error) {
	if s.cli == nil || s.ctx == nil {
		return "{}", fmt.Errorf("Docker client not initialized")
	}
	_, raw, err := s.cli.VolumeInspectWithRaw(s.ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get volume data: %v", err)
	}

	return string(raw), nil
}
