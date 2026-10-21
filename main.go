package main

import (
	"context"
	"embed"
	"log"

	"github.com/lburfeindt/trackops/internal/desktop"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := &desktop.App{}
	if err := wails.Run(&options.App{
		Title:     "TrackOps",
		Width:     640,
		Height:    480,
		MinWidth:  480,
		MinHeight: 360,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(context context.Context) { desktop.Startup(app, context) },
		Bind:      []interface{}{app},
	}); err != nil {
		log.Fatal(err)
	}
}
