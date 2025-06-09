package main

import (
	"context"
	"dockmate/app"
	"embed"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var appobj *app.App
var dockerContainersService *app.DockerContainersService
var dockerImagesService *app.DockerImagesService
var dockerVolumesService *app.DockerVolumesService
var dockerNetworksService *app.DockerNetworksService
var dockerLogsService *app.DockerLogsService
var dockerTerminalService *app.DockerContainersTerminal

func main() {
	// Create an instance of the app structure
	appobj = app.NewApp()
	dockerContainersService = app.NewDockerContainersService()
	dockerImagesService = app.NewDockerImagesService()
	dockerVolumesService = app.NewDockerVolumesService()
	dockerNetworksService = app.NewDockerNetworksService()
	dockerLogsService = app.NewDockerLogsService()
	dockerTerminalService = app.NewDockerTerminalService()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "DockMate",
		Width:     1024,
		Height:    768,
		MinWidth:  320,
		MinHeight: 240,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		OnStartup:        startup,
		Frameless:        true,
		Bind: []interface{}{
			appobj,
			dockerContainersService,
			dockerImagesService,
			dockerVolumesService,
			dockerNetworksService,
			dockerLogsService,
			dockerTerminalService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func startup(ctx context.Context) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}
	app.Startup(appobj, ctx, cli)
	app.StartupDockerContainersService(dockerContainersService, ctx, cli)
	app.StartupDockerImagesService(dockerImagesService, ctx, cli)
	app.StartupDockerVolumesService(dockerVolumesService, ctx, cli)
	app.StartupDockerNetworksService(dockerNetworksService, ctx, cli)
	app.StartupDockerLogsService(dockerLogsService, ctx, cli)
	app.StartupDockerTerminalService(dockerTerminalService, ctx, cli)
}
