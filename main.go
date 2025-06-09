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

func main() {
	// Create an instance of the app structure
	appobj = app.NewApp()
	dockerContainersService = app.NewDockerCOntainersService()
	dockerImagesService = app.NewDockerImagesService()
	dockerVolumesService = app.NewDockerVolumesService()
	dockerNetworksService = app.NewDockerNetworksService()
	dockerLogsService = app.NewDockerLogsService()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "DockMate",
		Width:  1024,
		Height: 768,
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
	app.StartupDockerLogsService(dockerLogsService, ctx, cli)
	app.StartupDockerNetworksService(dockerNetworksService, ctx, cli)
}
