package common

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort     string
	ServerAddress  string
	RedisHost      string
	RedisPort      string
	DbHost         string
	DbPort         string
	SyncStream     string
	SyncBatchCount int
}

func GetConfig() *Config {
	batchCount, err := strconv.Atoi(os.Getenv("SYNC_BATCH_COUNT"))
	if err != nil {
		log.Panicln("ERROR: Batch count must be an integer")
	}
	return &Config{
		ServerPort:     os.Getenv("SERVER_PORT"),
		ServerAddress:  os.Getenv("SERVER_ADDRESS"),
		RedisHost:      os.Getenv("REDIS_HOST"),
		RedisPort:      os.Getenv("REDIS_PORT"),
		DbHost:         os.Getenv("DB_HOST"),
		DbPort:         os.Getenv("DB_PORT"),
		SyncStream:     os.Getenv("SYNC_STREAM"),
		SyncBatchCount: batchCount,
	}
}
