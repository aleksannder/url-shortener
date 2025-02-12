package domain

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
)

type UrlRepository struct {
	cli *api.Client
}

func NewUrlRepository() (*UrlRepository, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%s", dbHost, dbPort)

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &UrlRepository{client}, nil
}
