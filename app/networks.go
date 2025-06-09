package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DockerNetworksService struct {
	DockerBaseService
	cancel context.CancelFunc
	mu     sync.Mutex
}

type NetworkInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewDockerNetworksService() *DockerNetworksService {
	return &DockerNetworksService{}
}

func StartupDockerNetworksService(s *DockerNetworksService, ctx context.Context, cli *client.Client) {
	s.ctx = ctx
	s.cli = cli
}

func (s *DockerNetworksService) List() ([]NetworkInfo, error) {
	if s.cli == nil || s.ctx == nil {
		return nil, fmt.Errorf("Docker client not initialized")
	}
	list, err := s.cli.NetworkList(s.ctx, network.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %v", err)
	}

	result := s.formatList(list)
	return result, nil
}

func (s *DockerNetworksService) StartWatching() error {
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
		eventFilter.Add("type", "network")
		eventsChan, errs := s.cli.Events(ctx, events.ListOptions{
			Filters: eventFilter,
		})

		for {
			select {
			case event := <-eventsChan:
				if event.Type == events.NetworkEventType {
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

func (s *DockerNetworksService) StopWatching() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}

func (s *DockerNetworksService) formatList(networks []network.Summary) []NetworkInfo {
	list := make([]NetworkInfo, 0, len(networks))
	for _, network := range networks {
		list = append(list, NetworkInfo{
			ID:   network.ID,
			Name: network.Name,
		})
	}
	return list
}

func (s *DockerNetworksService) sendListUpdate() {
	list, err := s.cli.NetworkList(s.ctx, network.ListOptions{})
	if err != nil {
		return
	}

	result := s.formatList(list)
	runtime.EventsEmit(s.ctx, "docker:networks", result)
}

func (s *DockerNetworksService) Remove(id string) error {
	if s.cli == nil || s.ctx == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	err := s.cli.NetworkRemove(s.ctx, id)
	return err
}

func (s *DockerNetworksService) Inspect(id string) (string, error) {
	if s.cli == nil || s.ctx == nil {
		return "{}", fmt.Errorf("Docker client not initialized")
	}
	_, raw, err := s.cli.NetworkInspectWithRaw(s.ctx, id, network.InspectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get network data: %v", err)
	}

	return string(raw), nil
}
