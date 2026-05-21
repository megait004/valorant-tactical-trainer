package wailsiface

import (
	"valorant-tactical-trainer/desktop/internal/domain/practice"
	"valorant-tactical-trainer/desktop/internal/infrastructure/store"
)

// PracticeService quản lý progress + lịch sử session luyện tập của user.
type PracticeService struct {
	store        *store.PracticeProgressStore
	sessionStore *store.PracticeSessionStore
}

func (s *PracticeService) GetPracticeProgress() (store.PracticeProgressState, error) {
	return s.store.Load()
}

func (s *PracticeService) SetPracticeProgress(itemID string, done bool) (store.PracticeProgressState, error) {
	return s.store.Set(itemID, done)
}

func (s *PracticeService) ResetPracticeProgress() (store.PracticeProgressState, error) {
	return s.store.Reset()
}

func (s *PracticeService) GetPracticeSessions() (practice.SessionState, error) {
	return s.sessionStore.Load()
}

func (s *PracticeService) FinishPracticeSession(input practice.SessionInput) (practice.SessionState, error) {
	return s.sessionStore.Add(input)
}
