package app

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DockerImagesService struct {
	DockerBaseService
	cancel context.CancelFunc
	mu     sync.Mutex
}

type ImageInfo struct {
	ID        string   `json:"id"`
	Size      int64    `json:"size"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
}

func NewDockerImagesService() *DockerImagesService {
	return &DockerImagesService{}
}

func StartupDockerImagesService(s *DockerImagesService, ctx context.Context, cli *client.Client) {
	s.ctx = ctx
	s.cli = cli
}

func (s *DockerImagesService) List() ([]ImageInfo, error) {
	if s.cli == nil || s.ctx == nil {
		return nil, fmt.Errorf("Docker client not initialized")
	}
	list, err := s.cli.ImageList(s.ctx, image.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %v", err)
	}

	result := s.formatList(list)
	return result, nil
}

func (s *DockerImagesService) StartWatching() error {
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
		eventFilter.Add("type", "image")
		eventsChan, errs := s.cli.Events(ctx, events.ListOptions{
			Filters: eventFilter,
		})

		for {
			select {
			case event := <-eventsChan:
				if event.Type == events.ImageEventType {
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

func (s *DockerImagesService) StopWatching() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
}

func (s *DockerImagesService) formatList(images []image.Summary) []ImageInfo {
	list := make([]ImageInfo, 0, len(images))
	for _, image := range images {
		list = append(list, ImageInfo{
			ID:        strings.TrimPrefix(image.ID, "sha256:"),
			Size:      image.Size,
			Tags:      image.RepoTags,
			CreatedAt: time.Unix(image.Created, 0).Format(time.RFC3339),
		})
	}
	return list
}

func (s *DockerImagesService) sendListUpdate() {
	if s.cli == nil || s.ctx == nil {
		return
	}

	images, err := s.cli.ImageList(s.ctx, image.ListOptions{})
	if err != nil {
		return
	}

	result := s.formatList(images)
	runtime.EventsEmit(s.ctx, "docker:images", result)
}

func (s *DockerImagesService) Remove(id string) error {
	if s.cli == nil || s.ctx == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	_, err := s.cli.ImageRemove(s.ctx, id, image.RemoveOptions{Force: true})
	return err
}

func (s *DockerImagesService) CreateAndStart(id string) error {
	if s.cli == nil || s.ctx == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	// Create container config
	config := &container.Config{
		Image: id,
		// Cmd:   []string{"/bin/sh"},
		Tty: false,
	}

	// Create container
	resp, err := s.cli.ContainerCreate(s.ctx, config, nil, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	// Start the container
	err = s.cli.ContainerStart(s.ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	return nil
}
