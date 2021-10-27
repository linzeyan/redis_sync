package redisync

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewRdb(addr, port, auth string, db uint) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: auth,
		DB:       int(db),
	})
}
