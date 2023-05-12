package redisync

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

/* Syncer includes Source and Destination redis, and uses Sync method to sync data. */
type Syncer struct {
	Source, Destination *redis.Client

	ctx context.Context
}

func (s *Syncer) setKeys(key string) {
	typ, err := s.Source.Type(s.ctx, key).Result()
	if err != nil {
		log.Printf("Recognize %s type error. %s", key, err)
	}
	switch typ {
	case "hash":
		s.hashType(key)
	case "list":
		s.listType(key)
	case "set":
		s.setType(key)
	case "string":
		s.stringType(key)
	case "zset":
		s.zsetType(key)
	}
}

func (s *Syncer) listType(key string) {
	value, err := s.Source.LRange(s.ctx, key, 0, -1).Result()
	if err != nil {
		log.Println("Get value error.", err)
	}
	s.Destination.Del(s.ctx, key)
	if _, err := s.Destination.RPush(s.ctx, key, value).Result(); err != nil {
		log.Println("Set key and value error.", map[string]interface{}{key: value}, err)
	}
}

func (s *Syncer) hashType(key string) {
	value, err := s.Source.HGetAll(s.ctx, key).Result()
	if err != nil {
		log.Println("Get value error.", err)
	}
	if _, err := s.Destination.HMSet(s.ctx, key, value).Result(); err != nil {
		log.Println("Set key and value error.", map[string]interface{}{key: value}, err)
	}
}

func (s *Syncer) setType(key string) {
	value, err := s.Source.SMembers(s.ctx, key).Result()
	if err != nil {
		log.Println("Get value error.", err)
	}
	if _, err := s.Destination.SAdd(s.ctx, key, value).Result(); err != nil {
		log.Println("Set key and value error.", map[string]interface{}{key: value}, err)
	}
}

func (s *Syncer) stringType(key string) {
	value, err := s.Source.Get(s.ctx, key).Result()
	if err != nil {
		log.Println("Get value error.", err)
	}
	ttl, err := s.Source.TTL(s.ctx, key).Result()
	if err != nil {
		log.Println("Get ttl error.", err)
	}
	if _, err := s.Destination.Set(s.ctx, key, value, ttl*time.Second).Result(); err != nil {
		log.Println("Set key and value error.", map[string]string{key: value}, err)
	}
}

func (s *Syncer) zsetType(key string) {
	value, err := s.Source.ZRange(s.ctx, key, 0, -1).Result()
	if err != nil {
		log.Println("Get value error.", err)
	}
	for _, v := range value {
		score, err := s.Source.ZScore(s.ctx, key, v).Result()
		if err != nil {
			log.Println("Get score error.", err)
		}
		member := &redis.Z{
			Score:  score,
			Member: v,
		}
		if _, err := s.Destination.ZAdd(s.ctx, key, member).Result(); err != nil {
			log.Println("Set key and value error.", map[string]interface{}{key: value}, err)
		}
	}
}

/* Sync lists all keys from Source server and writes to Destination server. */
func (s *Syncer) Sync() {
	keys, err := s.Source.Keys(s.ctx, "*").Result()
	if err != nil {
		log.Println(err)
	}
	var wg sync.WaitGroup
	wg.Add(len(keys))
	log.Println("Start syncer")
	for i := range keys {
		go func(k string) {
			defer wg.Done()
			s.setKeys(k)
		}(keys[i])
	}
	wg.Wait()
	log.Println("Syncer stop")
}

/* NewSyncer returns a pointer Syncer. */
func NewSyncer(src, dest *redis.Client) *Syncer {
	return &Syncer{
		Source:      src,
		Destination: dest,
		ctx:         context.Background(),
	}
}
