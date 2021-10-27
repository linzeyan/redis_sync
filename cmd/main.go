package main

import (
	"flag"
	"fmt"
	"os"

	redisync "github.com/linzeyan/redis_sync"
)

const (
	usage = `Sync redis keys and values

Usage: redisync [option...]

Options:
`
)

var (
	destHost   = flag.String("dh", "0.0.0.0", "Destination Redis host")
	destPort   = flag.String("dp", "6379", "Destination Redis port")
	destAuth   = flag.String("da", "", "Password for Destination Redis authentication ")
	destDb     = flag.Uint("dc", 0, "Destination Redis db")
	sourceHost = flag.String("sh", "127.0.0.1", "Source Redis host")
	sourcePort = flag.String("sp", "6379", "Source Redis port")
	sourceAuth = flag.String("sa", "", "Password for Source Redis authentication ")
	sourceDb   = flag.Uint("sc", 0, "Source Redis db")
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	redisync.Destination = redisync.NewRdb(*destHost, *destPort, *destAuth, *destDb)
	redisync.Source = redisync.NewRdb(*sourceHost, *sourcePort, *sourceAuth, *sourceDb)
	redisync.Sync()
}
