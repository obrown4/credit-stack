package handlers

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
