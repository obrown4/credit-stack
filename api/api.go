package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	auth "github.com/obrown4/credit-stack/internal/auth"
	"github.com/obrown4/credit-stack/internal/db"
)

type Server struct {
	server *http.Server
	addr   string
	client *db.Client
}

func NewServer(ctx context.Context, addr string, client *db.Client) *Server {
	pubMux := http.NewServeMux()
	privMux := http.NewServeMux()
	setRoutes(ctx, pubMux, privMux, client)

	server := &http.Server{
		Addr:    addr,
		Handler: pubMux,
	}

	log.Printf("Server started on %s\n", addr)

	return &Server{
		server: server,
		addr:   addr,
		client: client,
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

func setRoutes(ctx context.Context, pubMux *http.ServeMux, privMux *http.ServeMux, client *db.Client) {
	// auth API routes
	handleLogin(ctx, pubMux, client)
	handleRegister(ctx, pubMux, client)

	// protected routes
	handleLogout(ctx, privMux, client)
	handlePrintMsg(ctx, privMux)

	pubMux.Handle("/", validateSession(ctx, client, privMux))
}

func validateSession(ctx context.Context, client *db.Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "No session found", http.StatusUnauthorized)
			return
		}

		csrf, err := r.Cookie("csrf_token")
		if err != nil || csrf.Value == "" {
			http.Error(w, "CSRF token is required", http.StatusBadRequest)
			return
		}

		err = auth.AuthorizeUser(ctx, client, sessionCookie.Value, csrf.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		log.Printf("Session validated\n")

		// validated, now proceed
		next.ServeHTTP(w, r)
	})
}

// HTTP handlers - handle all HTTP concerns
func handlePrintMsg(ctx context.Context, r *http.ServeMux) {
	r.HandleFunc("POST /print", func(w http.ResponseWriter, r *http.Request) {
		msg := r.FormValue("msg")
		auth.PrintMessage(ctx, msg)

		response := map[string]string{
			"status":  "success",
			"message": "Message printed successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

func handleRegister(ctx context.Context, r *http.ServeMux, client *db.Client) {
	r.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Call business logic (validation happens inside)
		err := auth.RegisterUser(ctx, client, username, password)
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
	})
}

func handleLogin(ctx context.Context, r *http.ServeMux, client *db.Client) {
	r.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
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
		result, err := auth.LoginUser(ctx, client, username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set cookies (HTTP concern handled in API layer)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    result.SessionToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    result.CSRFToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false,
		})

		response := map[string]string{
			"status":  "success",
			"message": "Login successful",
			"user":    result.Username,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

func handleLogout(ctx context.Context, r *http.ServeMux, client *db.Client) {
	r.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get session token from cookie
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "No session found", http.StatusUnauthorized)
			return
		}

		// Get username from request (you might want to get this from the session instead)
		username := r.FormValue("username")
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}

		// Call business logic
		err = auth.LogoutUser(ctx, client, username, sessionCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Clear cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HttpOnly: false,
		})

		response := map[string]string{
			"status":  "success",
			"message": "Logged out successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

func handleProtectedEndpoint(ctx context.Context, r *http.ServeMux, client *db.Client) {
	r.HandleFunc("POST /auth", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "No session found", http.StatusUnauthorized)
			return
		}

		csrf := r.Header.Get("X-CSRF-Token")
		if csrf == "" {
			http.Error(w, "CSRF token is required", http.StatusBadRequest)
			return
		}

		err = auth.AuthorizeUser(ctx, client, sessionCookie.Value, csrf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		response := map[string]string{
			"status":  "success",
			"message": "Protected endpoint accessed successfully",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}
