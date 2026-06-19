package cache

import "project-layout/internal/shared/config"

type Redis struct {
	// client *redis.Client
}

func Open(cfg config.RedisConfig) *Redis {
	return &Redis{}
}
