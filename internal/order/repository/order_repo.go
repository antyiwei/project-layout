package repository

import "context"

type Repo struct {
	// DB *gorm.DB
}

func (r *Repo) Insert(ctx context.Context, userID int64, amount int64) error {
	return nil
}
