package store

import (
	"encoding/json"
	"fmt"
	"github.com/aleksannder/url-shortener/domain"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
)

type UrlRepository struct {
	cli *api.Client
}

func NewUrlRepository() (*UrlRepository, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	log.Printf("DB HOST: %s, DB PORT: %s", dbHost, dbPort)
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%s", dbHost, dbPort)

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &UrlRepository{client}, nil
}

func (ur *UrlRepository) Save(url *domain.URL) error {
	kv := ur.cli.KV()

	data, err := json.Marshal(url.URL)
	if err != nil {
		return err
	}

	keyValue := &api.KVPair{Key: url.ShortCode, Value: data}
	_, err = kv.Put(keyValue, nil)
	if err != nil {
		return err
	}

	return nil
}
