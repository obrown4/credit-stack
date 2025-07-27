package api

import (
	"context"
	"encoding/json"
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
	s.HandleFunc("POST /print", handlePrintMsg)
	s.HandleFunc("POST /login", handleLogin)
	s.HandleFunc("POST /logout", handleLogout)
	s.HandleFunc("POST /register", handleRegister)

	log.Printf("Routes Set\n")
	// service routes
}

// HTTP handlers - handle all HTTP concerns
func handlePrintMsg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	msg := r.FormValue("msg")
	auth.PrintMessage(msg)

	response := map[string]string{
		"status":  "success",
		"message": "Message printed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Call business logic
	err := auth.RegisterUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "User registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Call business logic
	err := auth.LoginUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Login successful",
		"user":    username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Call business logic
	err := auth.LogoutUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Logged out successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
