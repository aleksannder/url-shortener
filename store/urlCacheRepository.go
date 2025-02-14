package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/aleksannder/url-shortener/common"
	"github.com/aleksannder/url-shortener/domain"
	"github.com/go-redis/redis"
	"log"
	"os"
	"time"
)

var ctx = context.Background()

type UrlCacheRepository struct {
	Cli *redis.Client
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

	return &UrlCacheRepository{Cli: cli}, nil

}

func (ur *UrlCacheRepository) Ping() {
	val, _ := ur.Cli.Ping().Result()
	log.Printf("Redis URL db ping info: %x", val)
}

func (ur *UrlCacheRepository) Insert(url *domain.URL) (*domain.URL, error) {
	err := ur.Cli.Set(url.ShortCode, url.URL, time.Hour*24).Err()
	if err != nil {
		return nil, err
	}

	// Push to REDIS stream for sync

	_, err = ur.Cli.XAdd(&redis.XAddArgs{
		Stream: common.GetConfig().SyncStream,
		Values: map[string]interface{}{
			"url":       url.URL,
			"shortCode": url.ShortCode,
		},
	}).Result()

	if err != nil {
		return nil, err
	}
	return url, nil
}

func (ur *UrlCacheRepository) Redirect(shortLink string) (*domain.URL, error) {
	val, err := ur.Cli.Get(shortLink).Result()
	if err != nil {
		return nil, err
	}

	return &domain.URL{URL: val}, nil
}

func (ur *UrlCacheRepository) GetAll() ([]domain.URL, error) {

	var cursor uint64
	var n int
	var resultingKeys []string
	for {
		var keys []string
		var err error
		keys, cursor, err := ur.Cli.Scan(cursor, "*", 10).Result()
		if err != nil {
			return nil, err
		}
		resultingKeys = append(resultingKeys, keys...)
		n += len(keys)
		if cursor == 0 {
			break
		}
	}

	var urls []domain.URL
	for _, key := range resultingKeys {
		val, err := ur.Cli.Get(key).Result()
		if err != nil {
			log.Printf("Redis URL db get error: %v", err)
		}
		urls = append(urls, domain.URL{ShortCode: key, URL: val, ShortLink: fmt.Sprintf("localhost:%s/%s", "8000", key)})
	}

	return urls, nil
}
