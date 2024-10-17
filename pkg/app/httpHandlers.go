package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WebhookMessage структура для хранения данных вебхука
type WebhookMessage struct {
	Description string `json:"description"`
	Event       string `json:"event"`
	MessageType string `json:"messagetype"`
	Message     string `json:"message"`
	UserTGId    int64  `json:"usertgid"`
	ChatTGId    int    `json:"chattgid"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {

	println("webhook gained ", r.Method)
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var message WebhookMessage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	// Обрабатываем сообщение
	fmt.Println("Received webhook:")
	fmt.Printf("\tDescription: %s\n", message.Description)
	fmt.Printf("\tEvent: %s\n", message.Event)
	fmt.Printf("\tMessagetype: %s\n", message.MessageType)
	fmt.Printf("\tMessage: %s\n", message.Message)
	fmt.Printf("\tUserTGID: %d\n", message.UserTGId)
	fmt.Printf("\tChatTGID: %d\n", message.ChatTGId)

	// Отправляем ответ об успешной обработке вебхука
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status": "ok"}`))
}

func (a App) registerHttpHandlers() error {
	http.HandleFunc("/webhook", webhookHandler)

	fmt.Println("http run")
	return http.ListenAndServe(":8055", nil)
}
