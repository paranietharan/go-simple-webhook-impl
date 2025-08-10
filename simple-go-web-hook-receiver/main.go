package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type WebhookPayload struct {
	Event     string                 `json:"event"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("Webhook receiver listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload WebhookPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Process the webhook
	fmt.Printf("Received webhook event: %s\n", payload.Event)
	fmt.Printf("Payload data: %+v\n", payload.Data)
	fmt.Printf("Timestamp: %d\n", payload.Timestamp)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
}
