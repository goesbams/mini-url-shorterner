package server

import "net/http"

type Server struct {
	httpServer *http.Server
}

func NewServer(handler http.Handler, port string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
