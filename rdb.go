package redisync

import (
	"github.com/go-redis/redis/v8"
)

/* NewRdb returns a new redis client. */
func NewRdb(addr, port, auth string, db uint) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: auth,
		DB:       int(db),
	})
}
