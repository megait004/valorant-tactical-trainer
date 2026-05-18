package tactical

import (
	"strings"
	"time"
)

type PlanMarker struct {
	ID    string  `json:"id"`
	Kind  string  `json:"kind"`
	Label string  `json:"label"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}

type PlanLine struct {
	ID    string  `json:"id"`
	Label string  `json:"label"`
	X1    float64 `json:"x1"`
	Y1    float64 `json:"y1"`
	X2    float64 `json:"x2"`
	Y2    float64 `json:"y2"`
}

type MapPlan struct {
	MapID     string       `json:"mapId"`
	Title     string       `json:"title"`
	Side      string       `json:"side"`
	Notes     string       `json:"notes"`
	Markers   []PlanMarker `json:"markers"`
	Lines     []PlanLine   `json:"lines"`
	UpdatedAt string       `json:"updatedAt"`
}

func DefaultPlan(mapID string) MapPlan {
	return MapPlan{
		MapID:   mapID,
		Title:   "Kế hoạch trước trận",
		Side:    "attack",
		Markers: []PlanMarker{},
		Lines:   []PlanLine{},
	}
}

func NormalizePlan(plan MapPlan) MapPlan {
	plan.MapID = strings.TrimSpace(strings.ToLower(plan.MapID))
	plan.Title = strings.TrimSpace(plan.Title)
	if plan.Title == "" {
		plan.Title = "Kế hoạch trước trận"
	}
	plan.Side = strings.TrimSpace(strings.ToLower(plan.Side))
	if plan.Side != "defense" && plan.Side != "both" {
		plan.Side = "attack"
	}
	plan.Notes = strings.TrimSpace(plan.Notes)
	if plan.Markers == nil {
		plan.Markers = []PlanMarker{}
	}
	if plan.Lines == nil {
		plan.Lines = []PlanLine{}
	}
	plan.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return plan
}
