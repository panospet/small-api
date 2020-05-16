package main

import (
	"log"

	"github.com/panospet/small-api/internal/config"
	"github.com/panospet/small-api/pkg/api"
	"github.com/panospet/small-api/pkg/cache"
	"github.com/panospet/small-api/pkg/services"
)

func main() {
	conf := config.NewConfig()
	db, err := services.NewDb(conf.MysqlPath)
	if err != nil {
		log.Fatal("Error while initializing db", err)
	}
	redis, err := cache.NewRedisCache(conf.RedisPath)
	bpApi := api.NewApi(db, redis)
	bpApi.Run()
}
