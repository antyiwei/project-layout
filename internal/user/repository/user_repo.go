package repository

import (
	"context"

	"project-layout/internal/user/model"
)

type Repo struct {
	// DB *gorm.DB
}

func (r *Repo) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return &model.User{ID: id}, nil
}
