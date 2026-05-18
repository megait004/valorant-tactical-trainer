package main

import (
	"embed"

	wailsiface "valorant-tactical-trainer/desktop/internal/interface/wails"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	services := wailsiface.NewServices()

	err := wails.Run(&options.App{
		Title:  "Valorant Tactical Trainer",
		Width:  1280,
		Height: 820,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 9, G: 10, B: 15, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			services.Settings,
			services.Analysis,
			services.Practice,
			services.Assistant,
			services.Tactical,
			services.Auth,
			services.Chat,
		},
	})
	if err != nil {
		println("err starting app:", err.Error())
	}
}
