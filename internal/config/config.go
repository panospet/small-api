package config

type Config struct {
	MysqlPath string
	RedisPath string
}

// todo env variable
func NewConfig() *Config {
	return &Config{
		MysqlPath: "bestprice:bestprice@(localhost:3305)/bestprice?parseTime=true",
		RedisPath: "redis://localhost:6380/1",
	}
}
