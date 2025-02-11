package services

import (
	"errors"
	"github.com/aleksannder/url-shortener/domain"
)

type UrlService struct {
	repo *domain.UrlRepository
}

func NewUrlService(repo *domain.UrlRepository) (*UrlService, error) {
	if repo == nil {
		return nil, errors.New("repository must be initiated")
	}
	return &UrlService{repo: repo}, nil
}
