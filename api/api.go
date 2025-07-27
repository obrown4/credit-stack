package api

import (
	"context"
	"log"
	"net/http"
	"time"

	auth "github.com/obrown4/credit-stack/internal/auth"
)

type Server struct {
	server *http.Server
	addr   string
}

func NewServer(ctx context.Context, addr string) *Server {
	router := http.NewServeMux()
	setRoutes(router)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &Server{
		server: server,
		addr:   addr,
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func setRoutes(s *http.ServeMux) {

	// auth API routes
	s.HandleFunc("POST /print", auth.PrintMsg)
	s.HandleFunc("POST /login", auth.Login)
	s.HandleFunc("POST /logout", auth.Logout)
	s.HandleFunc("POST /register", auth.Register)

	log.Printf("Routes Set\n")
	// service routes

}
