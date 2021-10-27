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
	keys, err := Local.Keys(ctx, "*").Result()
	if err != nil {
		fmt.Println(err)
	}
	for _, key := range keys {
		typ, err := Local.Type(ctx, key).Result()
		if err != nil {
			fmt.Println("Get type err.", err)
		}
		switch typ {
		case "hash":
			hashType(key)
		case "list":
			listType(key)
		case "set":
			setType(key)
		case "string":
			stringType(key)
		case "zset":
			zsetType(key)
		}
	}
}

func listType(key string) {
	value, err := Local.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println("Get value err.", err)
	}
	Remote.Del(ctx, key)
	if _, err := Remote.RPush(ctx, key, value).Result(); err != nil {
		fmt.Println("Set key and value err.", map[string]interface{}{key: value}, err)
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

func setType(key string) {
	value, err := Local.SMembers(ctx, key).Result()
	if err != nil {
		fmt.Println("Get value err.", err)
	}
	if _, err := Remote.SAdd(ctx, key, value).Result(); err != nil {
		fmt.Println("Set key and value err.", map[string]interface{}{key: value}, err)
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

func zsetType(key string) {
	value, err := Local.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println("Get value err.", err)
	}
	for _, v := range value {
		score, err := Local.ZScore(ctx, key, v).Result()
		if err != nil {
			fmt.Println("Get score error.", err)
		}
		member := &redis.Z{
			Score:  score,
			Member: v,
		}
		if _, err := Remote.ZAdd(ctx, key, member).Result(); err != nil {
			fmt.Println("Set key and value err.", map[string]interface{}{key: value}, err)
		}
	}
}
