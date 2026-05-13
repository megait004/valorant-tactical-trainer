package main

import (
	"context"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// AppInfo returns the initial desktop shell status for the frontend.
func (a *App) AppInfo() AppInfo {
	return AppInfo{
		Name:   "Valorant Tactical Trainer",
		Status: "Go core online",
		Stack:  []string{"Wails v2", "React", "Go", "SQLite"},
	}
}

type AppInfo struct {
	Name   string   `json:"name"`
	Status string   `json:"status"`
	Stack  []string `json:"stack"`
}
