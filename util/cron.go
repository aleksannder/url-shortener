package util

import (
	"github.com/aleksannder/url-shortener/store"
	"log"
	"time"
)

type Sync struct {
	Cache      *store.UrlCacheRepository
	Persistent *store.UrlRepository
}

func (s *Sync) Sync() {
	go func() {
		log.Printf("In ticker go func")
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			log.Printf("tick")
			s.syncToConsul()
		}
	}()
}

func (s *Sync) syncToConsul() {
	log.Printf("in sync to consul")
	cachedUrls, err := s.Cache.GetAll()
	if err != nil {
		log.Panicln(err)
	}

	for _, cachedUrl := range cachedUrls {
		log.Printf("%s", cachedUrl)
		err := s.Persistent.Save(&cachedUrl)
		if err != nil {
			log.Panicln(err)
		}
	}
}
