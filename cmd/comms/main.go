package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Swayamjimmy/RescueNet/internal/p2p"
)

type IncomingMsg struct {
	Message string `json:"message"`
}

var MessageArr []string

var messageMu sync.Mutex

var cr *p2p.ChatRoom

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Allow-Control-Allow-Headers", "Content-Type")
}

func StoreMessage(msg string) {
	messageMu.Lock()
	defer messageMu.Unlock()
	MessageArr = append(MessageArr, msg)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Only Get Method supported", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(MessageArr)
	if err != nil {

		http.Error(w, "failed to encode messages", http.StatusInternalServerError)
		return
	}
}
