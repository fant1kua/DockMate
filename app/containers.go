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

type DockerService struct {
	ctx   	context.Context
	cli			*client.Client
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

func NewDockerService() *DockerService {
	return &DockerService{}
}

func StartupDockerService(service *DockerService, ctx context.Context, cli *client.Client) {
	service.ctx = ctx
	service.cli = cli
}

func formatList(containers []container.Summary) []ContainersGroup {
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

func (s *DockerService) ListContainers() ([]ContainersGroup, error) {
	if s.cli == nil {
		return nil, fmt.Errorf("Docker client not initialized")
	}
	containers, err := s.cli.ContainerList(s.ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	return formatList(containers), nil
}

func (s *DockerService) StartWatching() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Don't start multiple listeners
	if s.cancel != nil {
		return
	}
	
	s.sendContainersUpdate()

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
					s.sendContainersUpdate()
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
 
func (ds *DockerService) StopWatching() {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	
	if ds.cancel != nil {
		ds.cancel()
		ds.cancel = nil
	}
}
 
func (s *DockerService) sendContainersUpdate() {
	containers, err := s.cli.ContainerList(s.ctx, container.ListOptions{All: true})
	if err != nil {
		return
	}

	runtime.EventsEmit(s.ctx, "docker:containers", formatList(containers))
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