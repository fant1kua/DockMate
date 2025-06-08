package app

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx             context.Context
	cli             *client.Client
	logStreamCtx    context.Context
	logStreamCancel context.CancelFunc
}

func NewApp() *App {
	return &App{}
}

func Startup(a *App, ctx context.Context) {
	fmt.Println("App Startup")
	a.ctx = ctx
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}
	a.cli = cli
}

func (a *App) QuitApp() {
	if a.cli != nil {
		a.cli.Close()
	}
	runtime.Quit(a.ctx)
}

func (a *App) MaximiseApp() {
	runtime.WindowToggleMaximise(a.ctx)
}

func (a *App) MinimiseApp() {
	runtime.WindowMinimise(a.ctx)
}

type ContainerInfo struct {
	ID     string   `json:"id"`
	Names  []string `json:"names"`
	Image  string   `json:"image"`
	Status string   `json:"status"`
	State  string   `json:"state"`
}

type ImageInfo struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Size      int64    `json:"size"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
}

type VolumeInfo struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Size      int64    `json:"size"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
}

type ComposeProject struct {
	Name       string          `json:"name"`
	Containers []ContainerInfo `json:"containers"`
}

func (a *App) ListContainers() ([]ContainerInfo, error) {
	if a.cli == nil {
		return nil, fmt.Errorf("Docker client not initialized")
	}
	containers, err := a.cli.ContainerList(a.ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	var containerInfos []ContainerInfo
	for _, container := range containers {
		containerInfos = append(containerInfos, ContainerInfo{
			ID:     container.ID[:12],
			Names:  container.Names,
			Image:  container.Image,
			Status: container.Status,
			State:  container.State,
		})
	}

	return containerInfos, nil
}

func (a *App) StartContainer(containerID string) error {
	if a.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return a.cli.ContainerStart(a.ctx, containerID, container.StartOptions{})
}

func (a *App) StopContainer(containerID string) error {
	if a.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return a.cli.ContainerStop(a.ctx, containerID, container.StopOptions{})
}

func (a *App) RestartContainer(containerID string) error {
	if a.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return a.cli.ContainerRestart(a.ctx, containerID, container.StopOptions{})
}

func (a *App) RemoveContainer(containerID string) error {
	if a.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return a.cli.ContainerRemove(a.ctx, containerID, container.RemoveOptions{})
}

func (a *App) KillContainer(containerID string) error {
	if a.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	return a.cli.ContainerKill(a.ctx, containerID, "SIGKILL")
}

func (a *App) ContainerInspect(containerID string) (string, error) {
	if a.cli == nil {
		return "", fmt.Errorf("Docker client not initialized")
	}
	_, raw, err := a.cli.ContainerInspectWithRaw(a.ctx, containerID, false)
	if err != nil {
		return "", fmt.Errorf("failed to get container data: %v", err)
	}

	return string(raw), nil
}

func (a *App) GetContainerLogs(containerID string) (string, error) {
	if a.cli == nil {
		return "", fmt.Errorf("Docker client not initialized")
	}

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false, // Don't follow, we'll poll instead
		Tail:       "100", // Get last 100 lines
		Timestamps: true,
	}

	reader, err := a.cli.ContainerLogs(a.ctx, containerID, options)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %v", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read container logs: %v", err)
	}

	return string(logs), nil
}

func (a *App) StreamContainerLogs(containerID string) error {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
		Tail:       "10",
	}

	out, err := a.cli.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.logStreamCtx = ctx
	a.logStreamCancel = cancel

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
					runtime.EventsEmit(a.ctx, "logStream", "ERROR: "+err.Error())
					break
				}
				runtime.EventsEmit(a.ctx, "logStream", string(line))
			}
		}
	}()

	return nil
}

func (a *App) StopContainerLogs() {
	if a.logStreamCancel != nil {
		a.logStreamCancel()
		a.logStreamCancel = nil
	}
}

func (a *App) ListImages() ([]ImageInfo, error) {
	if a.cli == nil {
		return nil, fmt.Errorf("Docker client not initialized")
	}
	images, err := a.cli.ImageList(a.ctx, image.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %v", err)
	}

	var imageInfos []ImageInfo
	for _, container := range images {
		imageInfos = append(imageInfos, ImageInfo{
			ID:        strings.TrimPrefix(container.ID, "sha256:"),
			Size:      container.Size,
			Tags:      container.RepoTags,
			CreatedAt: time.Unix(container.Created, 0).Format(time.RFC3339),
		})
	}

	return imageInfos, nil
}

func (a *App) DeleteImage(imageID string) error {
	if a.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	_, err := a.cli.ImageRemove(a.ctx, imageID, image.RemoveOptions{Force: true})
	return err
}

func (a *App) ListVolumes() ([]VolumeInfo, error) {
	if a.cli == nil {
		return nil, fmt.Errorf("Docker client not initialized")
	}
	volumes, err := a.cli.VolumeList(a.ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %v", err)
	}

	var volumeInfos []VolumeInfo
	for _, vol := range volumes.Volumes {
		volumeInfos = append(volumeInfos, VolumeInfo{
			Name:      vol.Name,
			CreatedAt: vol.CreatedAt,
		})
	}

	return volumeInfos, nil
}

func (a *App) DeleteVolume(id string) error {
	if a.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}
	err := a.cli.VolumeRemove(a.ctx, id, true)
	return err
}

func (a *App) ExecContainer(containerID string, command string) error {
	if a.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}

	execConfig := container.ExecOptions{
		Cmd:          []string{"/bin/sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Tty:          true,
	}

	execID, err := a.cli.ContainerExecCreate(a.ctx, containerID, execConfig)
	if err != nil {
		return fmt.Errorf("failed to create exec instance: %v", err)
	}

	resp, err := a.cli.ContainerExecAttach(a.ctx, execID.ID, container.ExecAttachOptions{
		Tty: true,
	})
	if err != nil {
		return fmt.Errorf("failed to attach to exec instance: %v", err)
	}

	go func() {
		defer resp.Close() // Move defer to the start of the goroutine
		reader := bufio.NewReader(resp.Reader)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				runtime.EventsEmit(a.ctx, "logStream", "ERROR: "+err.Error())
				break
			}
			runtime.EventsEmit(a.ctx, "logStream", string(line))
		}
	}()

	return nil
}

func (a *App) CreateAndStartContainer(imageID string) error {
	if a.cli == nil {
		return fmt.Errorf("Docker client not initialized")
	}

	// Create container config
	config := &container.Config{
		Image: imageID,
		Cmd:   []string{"/bin/sh"},
		Tty:   true,
	}

	// Create container
	resp, err := a.cli.ContainerCreate(a.ctx, config, nil, nil, nil, "")
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	// Start the container
	err = a.cli.ContainerStart(a.ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	return nil
}

func (a *App) ListContainersByCompose() ([]ComposeProject, error) {
	if a.cli == nil {
		return nil, fmt.Errorf("Docker client not initialized")
	}

	containers, err := a.cli.ContainerList(a.ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	// Map to store containers by project
	projectMap := make(map[string][]ContainerInfo)
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
		projectName := container.Labels["com.docker.compose.project"]
		if projectName != "" {
			projectMap[projectName] = append(projectMap[projectName], containerInfo)
		} else {
			standaloneContainers = append(standaloneContainers, containerInfo)
		}
	}

	// Convert map to slice of ComposeProject
	var projects []ComposeProject
	for name, containers := range projectMap {
		projects = append(projects, ComposeProject{
			Name:       name,
			Containers: containers,
		})
	}

	// Add standalone containers as a special project
	if len(standaloneContainers) > 0 {
		projects = append(projects, ComposeProject{
			Name:       "Standalone",
			Containers: standaloneContainers,
		})
	}

	return projects, nil
}
