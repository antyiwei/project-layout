package service

import (
	"context"

	"project-layout/internal/user/model"
	"project-layout/internal/user/repository"
)

type Service struct {
	Repo *repository.Repo
}

func (s *Service) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.Repo.FindByID(ctx, id)
}
