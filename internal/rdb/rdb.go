package rdb

import (
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	defaultTimeout = 3 * time.Second
	// 常用缓存时长

	minute3  = 3 * time.Minute
	minute5  = 5 * time.Minute
	minute10 = 10 * time.Minute
	minute30 = 30 * time.Hour
	minute60 = time.Hour
)

type RedisRepo struct {
	Client *redis.Client
}

func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{Client: client}
}
