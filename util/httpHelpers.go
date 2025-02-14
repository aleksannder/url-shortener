package util

import (
	"encoding/json"
	"github.com/aleksannder/url-shortener/domain"
	"io"
	"mime"
	"net/http"
)

func CheckContentType(header, allowedContentType string) (bool, error) {
	mediaType, _, err := mime.ParseMediaType(header)
	if err != nil {
		return false, err
	}

	if mediaType != allowedContentType {
		return false, nil
	}

	return true, nil
}

func DecodeBody(r io.Reader) (*domain.URL, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var url domain.URL
	if err := dec.Decode(&url); err != nil {
		return nil, err
	}
	return &url, nil
}

func RenderJSON(w http.ResponseWriter, v interface{}, statusCode int) {
	marshal, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(marshal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
