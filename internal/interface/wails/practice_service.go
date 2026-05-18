package wailsiface

import (
	"valorant-tactical-trainer/desktop/internal/domain/practice"
	"valorant-tactical-trainer/desktop/internal/infrastructure/localstore"
)

// PracticeService quản lý progress + lịch sử session luyện tập của user.
type PracticeService struct {
	store        *localstore.PracticeProgressStore
	sessionStore *localstore.PracticeSessionStore
}

func (s *PracticeService) GetPracticeProgress() (localstore.PracticeProgressState, error) {
	return s.store.Load()
}

func (s *PracticeService) SetPracticeProgress(itemID string, done bool) (localstore.PracticeProgressState, error) {
	return s.store.Set(itemID, done)
}

func (s *PracticeService) ResetPracticeProgress() (localstore.PracticeProgressState, error) {
	return s.store.Reset()
}

func (s *PracticeService) GetPracticeSessions() (practice.SessionState, error) {
	return s.sessionStore.Load()
}

func (s *PracticeService) FinishPracticeSession(input practice.SessionInput) (practice.SessionState, error) {
	return s.sessionStore.Add(input)
}
