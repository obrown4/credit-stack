package api

import (
	"log"
	"net/http"

	auth "github.com/obrown4/credit-stack/internal/auth"
)

func NewServer() http.ServeMux {
	return *http.NewServeMux()
}

func SetRoutes(s *http.ServeMux) {

	// auth API routes
	s.HandleFunc("POST /print", auth.PrintMsg)
	s.HandleFunc("POST /login", auth.Login)
	s.HandleFunc("POST /logout", auth.Logout)
	s.HandleFunc("POST /register", auth.Register)

	log.Printf("Routes Set\n")
	// service routes

}
