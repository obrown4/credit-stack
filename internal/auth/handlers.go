package auth

import (
	"fmt"
	"net/http"
)

func PrintMsg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Print to server console
	fmt.Println("POST request received!")

	// Respond with JSON
	msg := r.FormValue("msg")
	fmt.Println("Message:", msg)
}

func Register(w http.ResponseWriter, r *http.Request) {}

func Login(w http.ResponseWriter, r *http.Request) {}

func Logout(w http.ResponseWriter, r *http.Request) {}
