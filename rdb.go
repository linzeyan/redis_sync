package redisync

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewRdb(addr, auth string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: auth,
		DB:       db,
	})
}
