package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aleksannder/url-shortener/common"
	"github.com/aleksannder/url-shortener/domain"
	"github.com/hashicorp/consul/api"
	"log"
	"strings"
)

type UrlRepository struct {
	cli *api.Client
}

func NewUrlRepository() (*UrlRepository, error) {
	dbHost := common.GetConfig().DbHost
	dbPort := common.GetConfig().DbPort

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

func (ur *UrlRepository) Redirect(shortCode string) (*domain.URL, error) {
	kv := ur.cli.KV()

	data, _, err := kv.Get(shortCode, nil)
	if data == nil {
		return nil, errors.New("short code not found")
	}
	if err != nil {
		return nil, err
	}

	var url domain.URL
	url.URL = ur.getURLFromShortCodeKV(data.Value)
	url.ShortCode = shortCode
	return &url, nil

}

func (ur *UrlRepository) getURLFromShortCodeKV(dataValue []byte) string {
	var result string
	result = string(dataValue)
	result = strings.ReplaceAll(result, `"`, "")
	return result
}
