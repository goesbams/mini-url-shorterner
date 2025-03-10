package handlers

import "mini-url-shortener/internal/services"

type URLHandler struct {
	urlService services.URLService
}

func NewURLHandler(urlService services.URLService) *URLHandler {
	return &URLHandler{urlService: urlService}
}
