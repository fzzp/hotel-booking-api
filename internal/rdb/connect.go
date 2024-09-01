package rdb

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, pwsd string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwsd,
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
