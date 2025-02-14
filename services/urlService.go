package services

import (
	"errors"
	"fmt"
	"github.com/aleksannder/url-shortener/common"
	"github.com/aleksannder/url-shortener/domain"
	"github.com/aleksannder/url-shortener/store"
	"github.com/aleksannder/url-shortener/util"
)

const baseUrl = "localhost:%s"

type UrlService struct {
	repo *store.UrlCacheRepository
}

func NewUrlService(repo *store.UrlCacheRepository) (*UrlService, error) {
	if repo == nil {
		return nil, errors.New("repository must be initiated")
	}
	return &UrlService{repo: repo}, nil
}

func (s *UrlService) Insert(url *domain.URL) (*domain.URL, error) {
	validationError := url.ValidateOnCreate()
	if validationError != nil {
		return nil, validationError
	}

	// Generate shortlink
	url, err := s.generateShortLinkFromUrl(url)
	if err != nil {
		return nil, err
	}

	url, err = s.repo.Insert(url)
	if err != nil {
		return nil, err
	}

	return url, nil
}

func (s *UrlService) Redirect(shortLink string) (*domain.URL, error) {
	url := &domain.URL{ShortLink: shortLink}
	if shortLink == "" {
		return nil, domain.ErrShortLinkInvalid
	}
	result, err := s.repo.Redirect(url.ShortLink)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *UrlService) generateShortLinkFromUrl(url *domain.URL) (*domain.URL, error) {
	baseUrl := fmt.Sprintf(baseUrl, common.GetConfig().ServerPort)

	originalLink := url.URL
	shortLink := util.Encode(util.Hash(originalLink))

	result := &domain.URL{
		URL:       url.URL,
		ShortLink: fmt.Sprintf("%s/%s", baseUrl, shortLink),
		ShortCode: shortLink,
	}

	return result, nil
}
