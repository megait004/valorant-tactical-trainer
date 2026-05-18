package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
)

// SuggestRecommendations gọi LLM cho từng finding trong batch, có disk cache.
// Implement analysis.Coach interface (legacy entry, chỉ override Recommendations).
func (c *Coach) SuggestRecommendations(ctx context.Context, snapshot analysis.PlayerSnapshot, findings []analysis.Finding) ([]analysis.Recommendation, error) {
	if c == nil {
		return nil, errors.New("coach chưa được khởi tạo (thiếu LLM_API_KEY trong .env)")
	}
	if len(findings) == 0 {
		return nil, nil
	}
	if len(findings) > c.maxFindings {
		findings = findings[:c.maxFindings]
	}

	payload := buildPromptPayload(snapshot, findings)
	cacheKey := c.cacheKeyRecommendations(payload)
	if cached, ok := c.readCacheRecommendations(cacheKey); ok {
		return cached, nil
	}

	body, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("err marshal prompt: %w", err)
	}
	userPrompt := "Snapshot + findings (JSON):\n" + string(body) +
		"\n\nTrả về JSON theo schema đã mô tả. Mỗi finding cần đúng 1 recommendation với findingId tương ứng."

	raw, err := c.Complete(ctx, c.systemPrompt, []ChatMessage{{Role: "user", Content: userPrompt}}, true)
	if err != nil {
		return nil, err
	}

	recs, err := parseRecommendationsJSON(raw)
	if err != nil {
		return nil, fmt.Errorf("err parse LLM output: %w (raw=%s)", err, truncate(raw, 200))
	}
	c.writeCacheRecommendations(cacheKey, recs)
	return recs, nil
}

// SuggestFullReport gọi LLM viết lại đồng thời Findings/Recommendations/
// PracticePlan dựa trên report rule engine. Implement analysis.Coach.
func (c *Coach) SuggestFullReport(ctx context.Context, input analysis.CoachInput) (analysis.CoachOutput, error) {
	if c == nil {
		return analysis.CoachOutput{}, errors.New("coach chưa được khởi tạo (thiếu LLM_API_KEY trong .env)")
	}
	if len(input.Findings) == 0 {
		return analysis.CoachOutput{}, nil
	}
	if len(input.Findings) > c.maxFindings {
		input.Findings = input.Findings[:c.maxFindings]
	}

	cacheKey := c.cacheKeyFullReport(input)
	if cached, ok := c.readCacheFullReport(cacheKey); ok {
		return cached, nil
	}

	body, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return analysis.CoachOutput{}, fmt.Errorf("err marshal full-report input: %w", err)
	}
	userPrompt := "Input (JSON):\n" + string(body) +
		"\n\nViết lại findings/recommendations/practicePlan theo đúng schema. Giữ nguyên id của findings."

	raw, err := c.Complete(ctx, fullReportSystemPrompt, []ChatMessage{{Role: "user", Content: userPrompt}}, true)
	if err != nil {
		return analysis.CoachOutput{}, err
	}
	out, err := parseFullReportJSON(raw)
	if err != nil {
		return analysis.CoachOutput{}, fmt.Errorf("err parse full-report JSON: %w (raw=%s)", err, truncate(raw, 200))
	}
	c.writeCacheFullReport(cacheKey, out)
	return out, nil
}
