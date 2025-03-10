package services

import "mini-url-shortener/internal/repositories"

type urlService struct {
	urlRepo repositories.URLRepository
}

func NewURLService(urlRepo repositories.URLRepository) URLService {
	return &urlService{urlRepo: urlRepo}
}
