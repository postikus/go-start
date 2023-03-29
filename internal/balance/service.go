package balance

import "go.uber.org/zap"

type Service struct {
	log  *zap.Logger
	repo *Repository
}

func NewService(log *zap.Logger, repo *Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}
