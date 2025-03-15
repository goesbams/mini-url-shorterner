package services

import "context"

type URLService interface {
	ShortenURL(ctx context.Context, originalURL string) (string, error)
	RedirectURL(ctx context.Context, shortCode string) (string, error)
}
