package db

import "project-layout/internal/shared/config"

type DB struct {
	// underlying *gorm.DB
}

func Open(cfg config.DatabaseConfig) *DB {
	return &DB{}
}
