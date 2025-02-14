package domain

import "errors"

type URL struct {
	URL       string `json:"url"`
	ShortLink string `json:"short_link"`
	ShortCode string `json:"short_code"`
}

var (
	ErrUrlEmpty         = errors.New("url is empty")
	ErrShortLinkInvalid = errors.New("short link is invalid")
)

type IUrl interface {
	Insert(url *URL) (*URL, error)
	Redirect(shortLink string) (*URL, error)
}

func (u *URL) ValidateOnCreate() error {
	if u.URL == "" {
		return ErrUrlEmpty
	}
	return nil
}

func (u *URL) ValidateOnRedirect() error {
	if u.ShortLink == "" {
		return ErrShortLinkInvalid
	}
	return nil
}
