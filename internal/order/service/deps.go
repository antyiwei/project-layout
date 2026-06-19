package service

import (
	"context"

	usermodel "project-layout/internal/user/model"
)

// Deps 集中声明 order 域对外部域的依赖接口。
type UserLookup interface {
	GetByID(ctx context.Context, id int64) (*usermodel.User, error)
}
