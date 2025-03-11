package services

import (
	"context"
	"fmt"
	"mini-url-shortener/internal/helpers"
	"mini-url-shortener/internal/models"
	"mini-url-shortener/internal/repositories"
)

type urlService struct {
	urlRepo repositories.URLRepository
}

func NewURLService(urlRepo repositories.URLRepository) URLService {
	return &urlService{urlRepo: urlRepo}
}

func (s *urlService) ShortenURL(ctx context.Context, originalURL string) (string, error) {

	if originalURL == "" {
		return "", fmt.Errorf("original url cannot be empty")
	}

	url := &models.URL{
		OriginalURL: originalURL,
		ShortCode:   helpers.GenerateShortCode(originalURL, 6),
	}

	if err := s.urlRepo.CreateShortCode(ctx, url); err != nil {
		return "", fmt.Errorf("error creating shorten url: %v", err)
	}

	return url.ShortCode, nil
}
