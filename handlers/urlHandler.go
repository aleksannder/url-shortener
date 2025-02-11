package handlers

import (
	"errors"
	"github.com/aleksannder/url-shortener/services"
	"net/http"
)

type UrlHandler struct {
	service *services.UrlService
}

func NewUrlHandler(service *services.UrlService) (*UrlHandler, error) {
	if service == nil {
		return nil, errors.New("service must be initiated")
	}
	return &UrlHandler{
		service: service,
	}, nil
}

func (h *UrlHandler) Insert(w http.ResponseWriter, r *http.Request) {}

func (h *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {}
