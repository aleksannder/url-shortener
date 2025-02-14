package handlers

import (
	"errors"
	"github.com/aleksannder/url-shortener/services"
	"github.com/aleksannder/url-shortener/util"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
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

func (h *UrlHandler) Insert(w http.ResponseWriter, r *http.Request) {
	checkMediaType, err := util.CheckContentType(r.Header.Get("Content-Type"), "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !checkMediaType {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	url, err := util.DecodeBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Go to service

	url, err = h.service.Insert(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.RenderJSON(w, url, http.StatusCreated)
}

func (h *UrlHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := mux.Vars(r)["shortCode"]
	url, err := h.service.Redirect(shortCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if url == nil || url.URL == "" {
		http.Error(w, "short code not found", http.StatusNotFound)
		return
	}

	if !strings.HasPrefix(url.URL, "http://") && !strings.HasPrefix(url.URL, "https://") {
		url.URL = "https://" + url.URL
	}
	http.Redirect(w, r, url.URL, http.StatusTemporaryRedirect)
	return
}
