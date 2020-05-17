package main

import (
	"flag"
	"fmt"
	"github.com/panospet/small-api/internal/config"
	"github.com/panospet/small-api/pkg/cache"
	"github.com/panospet/small-api/pkg/common"
	"github.com/panospet/small-api/pkg/services"
)

func main() {
	var workers int
	var amount int
	var redisOnly bool

	flag.BoolVar(&redisOnly, "redis-only", false, "populate only redis")
	flag.IntVar(&workers, "workers", 10, "number of workers")
	flag.IntVar(&amount, "amount", 1000, "total amount of products to add")
	flag.Parse()

	conf := config.NewConfig()

	db, err := services.NewDb(conf.MysqlPath)
	if err != nil {
		panic(err)
	}
	if !redisOnly {
		fmt.Println("populating mysql...")
		common.PopulateDb(db, workers, amount)
	}

	fmt.Println("populating redis...")
	redis, err := cache.NewRedisCache(conf.RedisPath)
	if err != nil {
		panic(err)
	}
	common.PopulateRedis(db, redis, workers)
}
