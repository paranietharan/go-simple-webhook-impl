package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookPayload struct {
	Event     string                 `json:"event"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

var users []User

func main() {
	http.HandleFunc("/user", createUserHandler)
	http.HandleFunc("/users", listUsersHandler)

	fmt.Println("Server listening on :8081...")
	http.ListenAndServe(":8081", nil)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Email == "" {
		http.Error(w, "Bad request: username and email are required", http.StatusBadRequest)
		return
	}

	users = append(users, user)

	payload := WebhookPayload{
		Event: "user.created",
		Data: map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		},
		Timestamp: time.Now().Unix(),
	}

	go sendWebhook(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
		"status":  "success",
	})
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func sendWebhook(payload WebhookPayload) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/webhook", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending webhook:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Webhook sent successfully")
	} else {
		fmt.Printf("Webhook failed with status: %s\n", resp.Status)
	}
}
