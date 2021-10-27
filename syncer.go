package redisync

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	Local, Remote *redis.Client
)

func Sync() {
	var err error
	keys, err := Local.Keys(ctx, "*").Result()
	if err != nil {
		fmt.Println(err)
	}
	for _, key := range keys {
		value, err := Local.Get(ctx, key).Result()
		if err != nil {
			fmt.Println("Get value err.", err)
		}
		ttl, err := Local.TTL(ctx, key).Result()
		if err != nil {
			fmt.Println("Get ttl err.", err)
		}
		if _, err := Remote.Set(ctx, key, value, ttl*time.Second).Result(); err != nil {
			fmt.Println("Set key and value err.", map[string]string{key: value}, err)
		}
	}
}
