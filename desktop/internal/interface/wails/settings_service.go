package wailsiface

import (
	"context"
	"fmt"

	"valorant-tactical-trainer/internal/infrastructure/storage"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SettingsService struct {
	ctx   context.Context
	store *storage.Store
}

func NewSettingsService(store *storage.Store) *SettingsService {
	return &SettingsService{store: store}
}

func (service *SettingsService) Startup(ctx context.Context) {
	service.ctx = ctx
}

type ResetResult struct {
	Message string `json:"message"`
}

func (service *SettingsService) ResetAllData() (ResetResult, error) {
	if service.ctx != nil {
		choice, err := runtime.MessageDialog(service.ctx, runtime.MessageDialogOptions{
			Type:          runtime.QuestionDialog,
			Title:         "Reset local data",
			Message:       "Reset all local Valorant Tactical Trainer data on this machine? This deletes player, consent, cache, matches, reports, and recommendations.",
			Buttons:       []string{"Reset", "Cancel"},
			DefaultButton: "Cancel",
			CancelButton:  "Cancel",
		})
		if err != nil {
			return ResetResult{}, fmt.Errorf("err showing reset dialog: %w", err)
		}
		if choice != "Reset" {
			return ResetResult{Message: "reset cancelled"}, nil
		}
	}

	if err := service.store.ResetAll(context.Background()); err != nil {
		return ResetResult{}, fmt.Errorf("err reset data: %w", err)
	}

	return ResetResult{Message: "local data reset"}, nil
}
