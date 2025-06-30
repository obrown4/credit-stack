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

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	passowrd := r.FormValue("password")

	if len(username) < 8 || len(passowrd) < 8 {
		http.Error(w, "Username and password must be at least 8 characters long",
			http.StatusBadRequest)
		return
	}

	// users := db.Client.Database("creditStack").Collection("users")
	// filter := bson.D{{"name", "Bagels N Buns"}}

}

func Login(w http.ResponseWriter, r *http.Request) {}

func Logout(w http.ResponseWriter, r *http.Request) {}
