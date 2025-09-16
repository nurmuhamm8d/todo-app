package service

import "todo-app/backend/internal/repository"

type Stats struct {
	Total          int `json:"total"`
	Active         int `json:"active"`
	Completed      int `json:"completed"`
	Overdue        int `json:"overdue"`
	HighPriority   int `json:"highPriority"`
	MediumPriority int `json:"mediumPriority"`
	LowPriority    int `json:"lowPriority"`
}

type StatsService struct {
	repo repository.StatsRepository
}

func NewStatsService(r repository.StatsRepository) *StatsService {
	return &StatsService{repo: r}
}

func (s *StatsService) Snapshot() (Stats, error) {
	ss, err := s.repo.Snapshot()
	if err != nil {
		return Stats{}, err
	}
	return Stats{
		Total:          ss.Total,
		Active:         ss.Active,
		Completed:      ss.Completed,
		Overdue:        ss.Overdue,
		HighPriority:   ss.HighPriority,
		MediumPriority: ss.MediumPriority,
		LowPriority:    ss.LowPriority,
	}, nil
}
