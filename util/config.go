package util

import "os"

type Config struct {
	ServerPort string
	RedisHost  string
	RedisPort  string
}

func GetConfig() Config {
	return Config{
		ServerPort: os.Getenv("SERVER_PORT"),
		RedisHost:  os.Getenv("REDIS_HOST"),
		RedisPort:  os.Getenv("REDIS_PORT"),
	}
}
