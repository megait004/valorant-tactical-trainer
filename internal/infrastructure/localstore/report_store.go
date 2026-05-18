package localstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"valorant-tactical-trainer/desktop/internal/domain/analysis"
)

type StoredReport struct {
	Report    analysis.Report `json:"report"`
	Source    string          `json:"source"`
	Cached    bool            `json:"cached"`
	FetchedAt string          `json:"fetchedAt"`
	Message   string          `json:"message"`
}

type ReportStore struct {
	path string
}

func NewReportStore() (*ReportStore, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	return &ReportStore{path: filepath.Join(configDir, "Valorant Tactical Trainer", "last_report.json")}, nil
}

func NewReportStoreAt(path string) *ReportStore {
	return &ReportStore{path: path}
}

func (s *ReportStore) LoadLastReport() (StoredReport, bool, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return StoredReport{}, false, nil
	}
	if err != nil {
		return StoredReport{}, false, err
	}

	var report StoredReport
	if err := json.Unmarshal(data, &report); err != nil {
		return StoredReport{}, false, err
	}
	return report, true, nil
}

func (s *ReportStore) SaveLastReport(report StoredReport) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

// Delete xóa file last_report.json. Idempotent — không lỗi nếu file không tồn tại.
func (s *ReportStore) Delete() error {
	if err := os.Remove(s.path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
