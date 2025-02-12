package domain

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
)

type UrlCacheRepository struct {
	cli *redis.Client
}

func NewUrlCacheRepository() (*UrlCacheRepository, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" || redisPort == "" {
		return nil, errors.New("database variables not correctly initiated")
	}

	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	return &UrlCacheRepository{cli: cli}, nil

}

func (ur *UrlCacheRepository) Ping() {
	val, _ := ur.cli.Ping().Result()
	log.Printf("Redis URL db ping info: %x", val)
}
