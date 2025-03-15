package handlers

import (
	"encoding/json"
	"log"
	"mini-url-shortener/internal/models"
	"mini-url-shortener/internal/services"
	"net/http"
)

type URLHandler struct {
	urlService services.URLService
}

func NewURLHandler(urlService services.URLService) *URLHandler {
	return &URLHandler{urlService: urlService}
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var request models.URLRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("error decoding request body", err)
		http.Error(w, "invalid url body request", http.StatusBadRequest)

		return
	}

	shortURL, err := h.urlService.ShortenURL(r.Context(), request.OriginalURL)
	if err != nil {
		log.Println("error shortening url", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	response := models.URLResponse{ShortCode: shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortcode")

	originalURL, err := h.urlService.RedirectURL(r.Context(), shortCode)
	if err != nil {
		log.Println("error redirecting to url", err)
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
