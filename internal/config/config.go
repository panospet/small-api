package config

import (
	"log"
	"os"
)

type Config struct {
	MysqlPath string
	RedisPath string
}

func NewConfig() *Config {
	mysqlPath := os.Getenv("MYSQL_PATH")
	if mysqlPath == "" {
		log.Println("Mysql path from env is empty. Using default value")
		mysqlPath = "bestprice:bestprice@(localhost:3305)/bestprice?parseTime=true"
	}
	redisPath := os.Getenv("REDIS_PATH")
	if redisPath == "" {
		log.Println("Redis path from env is empty. Using default value")
		redisPath = "redis://localhost:6380/1"
	}
	return &Config{
		MysqlPath: mysqlPath,
		RedisPath: redisPath,
	}
}
