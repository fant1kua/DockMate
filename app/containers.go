package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DockerContainersService struct {
	DockerBaseService
	cancel	context.CancelFunc
	mu     	sync.Mutex
}

type ContainerInfo struct {
	ID     string   `json:"id"`
	Names  []string `json:"names"`
	Image  string   `json:"image"`
	Status string   `json:"status"`
	State  string   `json:"state"`
}

type ContainersGroup struct {
	Name       string          `json:"name"`
	Containers []ContainerInfo `json:"containers"`
}

func NewDockerCOntainersService() *DockerContainersService {
	return &DockerContainersService{}
}

func StartupDockerContainersService(s *DockerContainersService, ctx context.Context, cli *client.Client) {
	s.ctx = ctx
	s.cli = cli
}

func (s *DockerContainersService) List() ([]ContainersGroup, error) {
	if s.cli == nil {
		return nil, fmt.Errorf("Docker client not initialized")
	}
	list, err := s.cli.ContainerList(s.ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	return s.formatList(list), nil
}

func (s *DockerContainersService) StartWatching() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Don't start multiple listeners
	if s.cancel != nil {
		return
	}
	
	s.sendListUpdate()

	ctx, cancel := context.WithCancel(s.ctx)
	s.cancel = cancel
	
	go func() {
		eventFilter := filters.NewArgs()
    eventFilter.Add("type", "container")
		eventsChan, errs := s.cli.Events(ctx, events.ListOptions{
			Filters: eventFilter,
		})
		
		for {
			select {
			case event := <-eventsChan:
				if event.Type == events.ContainerEventType {
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
}
 
func (s *DockerContainersService) StopWatching() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}

func (s *DockerContainersService) formatList(containers []container.Summary) []ContainersGroup {
	list := make(map[string][]ContainerInfo)
	standaloneContainers := []ContainerInfo{}

	for _, container := range containers {
		containerInfo := ContainerInfo{
			ID:     container.ID[:12],
			Names:  container.Names,
			Image:  container.Image,
			Status: container.Status,
			State:  container.State,
		}

		// Check for compose project label
		name := container.Labels["com.docker.compose.project"]
		if name != "" {
			list[name] = append(list[name], containerInfo)
		} else {
			standaloneContainers = append(standaloneContainers, containerInfo)
		}
	}

	// Convert map to slice of ComposeProject
	var projects []ContainersGroup
	for name, containers := range list {
		projects = append(projects, ContainersGroup{
			Name:       name,
			Containers: containers,
		})
	}

	// Add standalone containers as a special project
	if len(standaloneContainers) > 0 {
		projects = append(projects, ContainersGroup{
			Name:       "Standalone",
			Containers: standaloneContainers,
		})
	}

	return projects
}
 
func (s *DockerContainersService) sendListUpdate() {
	containers, err := s.cli.ContainerList(s.ctx, container.ListOptions{All: true})
	if err != nil {
		return
	}

	result := s.formatList(containers)
	runtime.EventsEmit(s.ctx, "docker:containers", result)
}

func (s *DockerContainersService) Start(id string) error {
	if s.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return s.cli.ContainerStart(s.ctx, id, container.StartOptions{})
}

func (s *DockerContainersService) Stop(id string) error {
	if s.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return s.cli.ContainerStop(s.ctx, id, container.StopOptions{})
}

func (s *DockerContainersService) Restart(id string) error {
	if s.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return s.cli.ContainerRestart(s.ctx, id, container.StopOptions{})
}

func (s *DockerContainersService) Remove(id string) error {
	if s.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return s.cli.ContainerRemove(s.ctx, id, container.RemoveOptions{})
}

func (s *DockerContainersService) Kill(id string) error {
	if s.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return s.cli.ContainerKill(s.ctx, id, "SIGKILL")
}

func (s *DockerContainersService) Inspect(id string) (string, error) {
	if s.cli == nil {
		return "", fmt.Errorf("Docker client not initialized")
	}
	_, raw, err := s.cli.ContainerInspectWithRaw(s.ctx, id, false)
	if err != nil {
		return "", fmt.Errorf("failed to get container data: %v", err)
	}

	return string(raw), nil
}