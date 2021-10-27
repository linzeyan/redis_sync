package main

import redisync "github.com/linzeyan/redis_sync"

func main() {
	redisync.Local = redisync.NewRdb("localhost:6379", "", 0)
	redisync.Remote = redisync.NewRdb("192.168.185.9:6379", "", 1)
	redisync.Sync()
}
