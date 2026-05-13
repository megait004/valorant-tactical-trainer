package wailsiface

import (
	"context"
	"fmt"

	analysisdomain "valorant-tactical-trainer/internal/domain/analysis"
	"valorant-tactical-trainer/internal/infrastructure/storage"
)

type AnalysisService struct {
	store *storage.Store
}

func NewAnalysisService(store *storage.Store) *AnalysisService {
	return &AnalysisService{store: store}
}

type ReportDTO struct {
	ID              int64               `json:"id"`
	PlayerPUUID     string              `json:"playerPuuid"`
	MatchCount      int                 `json:"matchCount"`
	AverageKDA      float64             `json:"averageKda"`
	HeadshotPercent float64             `json:"headshotPercent"`
	AverageDamage   float64             `json:"averageDamage"`
	TopAgent        string              `json:"topAgent"`
	TopMap          string              `json:"topMap"`
	Summary         string              `json:"summary"`
	Findings        []FindingDTO        `json:"findings"`
	Recommendations []RecommendationDTO `json:"recommendations"`
}

type FindingDTO struct {
	Type        string   `json:"type"`
	Severity    string   `json:"severity"`
	Confidence  float64  `json:"confidence"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Evidence    []string `json:"evidence"`
}

type RecommendationDTO struct {
	Title    string   `json:"title"`
	Drill    string   `json:"drill"`
	Priority string   `json:"priority"`
	Reason   string   `json:"reason"`
	Evidence []string `json:"evidence"`
	Status   string   `json:"status"`
}

func (service *AnalysisService) GenerateReport(puuid string) (ReportDTO, error) {
	if puuid == "" {
		return ReportDTO{}, fmt.Errorf("InvalidPlayer: puuid is required")
	}

	ctx := context.Background()
	matches, err := service.store.MatchesForPlayer(ctx, puuid)
	if err != nil {
		return ReportDTO{}, err
	}

	report := analysisdomain.GenerateReport(puuid, matches)
	savedReport, err := service.store.SaveReport(ctx, report)
	if err != nil {
		return ReportDTO{}, err
	}

	return toReportDTO(savedReport), nil
}

func toReportDTO(report analysisdomain.Report) ReportDTO {
	findings := make([]FindingDTO, 0, len(report.Findings))
	for _, finding := range report.Findings {
		findings = append(findings, FindingDTO{
			Type:        finding.Type,
			Severity:    finding.Severity,
			Confidence:  finding.Confidence,
			Title:       finding.Title,
			Description: finding.Description,
			Evidence:    finding.Evidence,
		})
	}

	recommendations := make([]RecommendationDTO, 0, len(report.Recommendations))
	for _, recommendation := range report.Recommendations {
		recommendations = append(recommendations, RecommendationDTO{
			Title:    recommendation.Title,
			Drill:    recommendation.Drill,
			Priority: recommendation.Priority,
			Reason:   recommendation.Reason,
			Evidence: recommendation.Evidence,
			Status:   recommendation.Status,
		})
	}

	return ReportDTO{
		ID:              report.ID,
		PlayerPUUID:     report.PlayerPUUID,
		MatchCount:      report.MatchCount,
		AverageKDA:      report.AverageKDA,
		HeadshotPercent: report.HeadshotPercent,
		AverageDamage:   report.AverageDamage,
		TopAgent:        report.TopAgent,
		TopMap:          report.TopMap,
		Summary:         report.Summary,
		Findings:        findings,
		Recommendations: recommendations,
	}
}
