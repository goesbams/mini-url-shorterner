package services

import "context"

type URLService interface {
	ShortenURL(ctx context.Context, originalURL string) (string, error)
}
