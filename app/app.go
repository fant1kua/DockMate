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
	ctx             context.Context
	cli             *client.Client
	logStreamCtx    context.Context
	logStreamCancel context.CancelFunc
}

type DockerBaseService struct {
	ctx context.Context
	cli *client.Client
}

func NewApp() *App {
	return &App{}
}

func Startup(a *App, ctx context.Context, cli *client.Client) {
	a.ctx = ctx
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
