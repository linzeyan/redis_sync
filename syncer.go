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
		typ, err := Local.Type(ctx, key).Result()
		if err != nil {
			fmt.Println("Get type err.", err)
		}
		fmt.Println(typ)
		switch typ {
		case "string":
			stringType(key)
		case "hash":
			hashType(key)
		}
	}
}

func stringType(key string) {
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

func hashType(key string) {
	value, err := Local.HGetAll(ctx, key).Result()
	if err != nil {
		fmt.Println("Get value err.", err)
	}
	if _, err := Remote.HMSet(ctx, key, value).Result(); err != nil {
		fmt.Println("Set key and value err.", map[string]interface{}{key: value}, err)
	}
}
