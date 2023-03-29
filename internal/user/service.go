package user

import (
	"context"
	"fmt"
	"github.com/postikus/go-starter/internal/balance"
	"github.com/postikus/go-starter/model"
	"go.uber.org/zap"
)

type Service struct {
	log         *zap.Logger
	repo        *Repository
	balanceRepo *balance.Repository
}

func NewService(log *zap.Logger, repo *Repository, balanceRepo *balance.Repository) *Service {
	return &Service{
		log:         log,
		repo:        repo,
		balanceRepo: balanceRepo,
	}
}

func (s *Service) New(ctx context.Context, user *model.User) (out *model.User, err error) {
	tx, closer, err := s.repo.StartTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open tx: %w", err)
	}
	defer closer(err)

	user, err = s.repo.
		WithTx(tx).
		Insert(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	user.Balance, err = s.balanceRepo.
		WithTx(tx).
		Insert(ctx, &model.Balance{
			UserID: user.ID,
			Amount: 100,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to create new balance: %w", err)
	}

	return user, nil
}
