package wailsiface

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type WindowService struct {
	ctx context.Context
}

type WindowModeResult struct {
	Overlay bool   `json:"overlay"`
	Message string `json:"message"`
}

func NewWindowService() *WindowService {
	return &WindowService{}
}

func (service *WindowService) Startup(ctx context.Context) {
	service.ctx = ctx
}

func (service *WindowService) SetAssistantOverlay(enabled bool) (WindowModeResult, error) {
	if service.ctx == nil {
		return WindowModeResult{}, fmt.Errorf("WindowUnavailable: desktop context not ready")
	}

	if enabled {
		runtime.WindowSetAlwaysOnTop(service.ctx, true)
		runtime.WindowSetSize(service.ctx, 520, 760)
		runtime.WindowSetTitle(service.ctx, "VTA Overlay - Valorant Tactical Trainer")
		return WindowModeResult{Overlay: true, Message: "assistant overlay enabled"}, nil
	}

	runtime.WindowSetAlwaysOnTop(service.ctx, false)
	runtime.WindowSetSize(service.ctx, 1024, 768)
	runtime.WindowCenter(service.ctx)
	runtime.WindowSetTitle(service.ctx, "valorant-tactical-trainer")
	return WindowModeResult{Overlay: false, Message: "assistant overlay disabled"}, nil
}
