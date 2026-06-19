package service

import (
	"context"
	"fmt"

	"project-layout/internal/order/repository"
)

type Service struct {
	Repo  *repository.Repo
	Users UserLookup
}

func (s *Service) Create(ctx context.Context, userID int64, amount int64) error {
	user, err := s.Users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user != nil && !user.Active {
		return fmt.Errorf("user %d is inactive", userID)
	}
	return s.Repo.Insert(ctx, userID, amount)
}
