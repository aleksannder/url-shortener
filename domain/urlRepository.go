package domain

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
)

type UrlRepository struct {
	cli *redis.Client
}

func NewUrlRepository() (*UrlRepository, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" || redisPort == "" {
		return nil, errors.New("database variables not correctly initiated")
	}

	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	return &UrlRepository{cli: cli}, nil

}

func (ur *UrlRepository) Ping() {
	val, _ := ur.cli.Ping().Result()
	log.Printf("Redis URL db ping info: %x", val)
}
