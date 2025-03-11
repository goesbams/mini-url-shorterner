package routes

import (
	"mini-url-shortener/internal/handlers"
	"net/http"
)

func SetupRoutes(urlHandler *handlers.URLHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	mux.HandleFunc("POST /shorten", urlHandler.ShortenURL)

	return mux
}
