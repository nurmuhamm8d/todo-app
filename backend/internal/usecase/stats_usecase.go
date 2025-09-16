package usecase

import "todo-app/backend/internal/service"

type StatsUsecase struct {
	svc *service.StatsService
}

func NewStatsUsecase(s *service.StatsService) *StatsUsecase {
	return &StatsUsecase{svc: s}
}

func (u *StatsUsecase) Snapshot() (service.Stats, error) {
	return u.svc.Snapshot()
}
