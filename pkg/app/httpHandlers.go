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

	//TODO для Илюхи: надо принимать номер телефона + ФИО
}

func (a App) webhookHandler(w http.ResponseWriter, r *http.Request) {

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

	if message.Event == "new_msg" {
		// Обрабатываем сообщение
		a.Logger.Printf("Received webhook:")
		a.Logger.Printf("\tDescription: %s\n", message.Description)
		a.Logger.Printf("\tEvent: %s\n", message.Event)
		a.Logger.Printf("\tMessagetype: %s\n", message.MessageType)
		a.Logger.Printf("\tMessage: %s\n", message.Message)
		a.Logger.Printf("\tUserTGID: %d\n", message.UserTGId)
		a.Logger.Printf("\tChatTGID: %d\n", message.ChatTGId)

		// Отправляем ответ об успешной обработке вебхука
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))

		a.sendWebhookResult(message)
	} else {
		a.Logger.Printf("Ивент вебхука: %s", message.Event)
	}

}

func (a App) registerHttpHandlers() error {
	http.HandleFunc("/webhook", a.webhookHandler)
	a.Logger.Printf("http run")
	return http.ListenAndServe(":8055", nil)
}
