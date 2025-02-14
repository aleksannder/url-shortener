package util

import (
	"fmt"
	"github.com/aleksannder/url-shortener/common"
	"github.com/aleksannder/url-shortener/domain"
	"github.com/aleksannder/url-shortener/store"
	"github.com/go-redis/redis"
	"log"
)

const baseUrl = "localhost:%s"

type Sync struct {
	Cache      *store.UrlCacheRepository
	Persistent *store.UrlRepository
}

func (s *Sync) Sync() {
	baseUrl := fmt.Sprintf(baseUrl, common.GetConfig().ServerPort)
	for {
		msgs, err := s.Cache.Cli.XRead(&redis.XReadArgs{
			Streams: []string{common.GetConfig().SyncStream, "0"},
			Count:   int64(common.GetConfig().SyncBatchCount),
			Block:   0,
		}).Result()

		if err != nil {
			log.Println(err)
			continue
		}

		for _, msg := range msgs[0].Messages {
			shortCode := msg.Values["shortCode"].(string)
			url := msg.Values["url"].(string)

			insert := &domain.URL{URL: url, ShortCode: shortCode, ShortLink: fmt.Sprintf("%s/%s", baseUrl, shortCode)}
			err := s.Persistent.Save(insert)
			if err == nil {
				s.Cache.Cli.XAck(common.GetConfig().SyncStream, "consumer_group", msg.ID)
			} else {
				log.Println(err)
			}
		}
	}
}
