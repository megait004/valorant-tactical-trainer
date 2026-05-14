package main

import (
	"context"
	"embed"
	"log"

	"valorant-tactical-trainer/internal/infrastructure/storage"
	wailsiface "valorant-tactical-trainer/internal/interface/wails"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	store, err := storage.Open(context.Background())
	if err != nil {
		log.Fatalf("err opening storage: %v", err)
	}
	defer store.Close()
	playerService := wailsiface.NewPlayerService(store)
	matchService := wailsiface.NewMatchService(store)
	rankService := wailsiface.NewRankService(store)
	analysisService := wailsiface.NewAnalysisService(store)
	assistantService := wailsiface.NewAssistantService(store)
	settingsService := wailsiface.NewSettingsService(store)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "valorant-tactical-trainer",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			settingsService.Startup(ctx)
		},
		Bind: []interface{}{
			app,
			playerService,
			matchService,
			rankService,
			analysisService,
			assistantService,
			settingsService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
