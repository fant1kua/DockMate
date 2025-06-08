package app

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
	cli *client.Client
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

	go func() {
		defer out.Close()
		reader := bufio.NewReader(out)
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
