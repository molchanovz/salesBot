package app

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// WebhookMessage структура для хранения данных вебхука
type WebhookMessage struct {
	Description string  `json:"description"`
	Event       string  `json:"event"`
	MessageType string  `json:"messagetype"`
	Message     string  `json:"message"`
	UserTGId    int64   `json:"usertgid"`
	ChatTGId    int     `json:"chattgid"`
	SenderTgId  *int64  `json:"senderTgId"`
	Nickname    *string `json:"nickname"`
}

func (a App) webhookHandler(c echo.Context) error {
	r := c.Request()
	println("webhook gained ", r.Method)
	if r.Method != "POST" {
		return echo.ErrMethodNotAllowed
	}

	var message WebhookMessage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	defer r.Body.Close()
	if err != nil {
		a.Logger.Errorf("%v", err)
		return echo.ErrBadRequest
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
		if message.SenderTgId != nil {
			a.Logger.Printf("\tSenderTGID: %d\n", *message.SenderTgId)
		}
		if message.Nickname != nil {
			a.Logger.Printf("\tNickname: %s\n", *message.Nickname)
		}

		a.sendWebhookResult(message)
	} else {
		a.Logger.Printf("Ивент вебхука: %s", message.Event)
	}
	return nil
}

type Lead struct {
	Add []struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		StatusID          string `json:"status_id"`
		Price             string `json:"price"`
		ResponsibleUserID string `json:"responsible_user_id"`
		LastModified      string `json:"last_modified"`
		ModifiedUserID    string `json:"modified_user_id"`
		CreatedUserID     string `json:"created_user_id"`
		DateCreate        string `json:"date_create"`
		PipelineID        string `json:"pipeline_id"`
		AccountID         string `json:"account_id"`
		CreatedAt         string `json:"created_at"`
		UpdatedAt         string `json:"updated_at"`
	} `json:"add"`
}

func (a App) webhookAmoCRMHandler(c echo.Context) error {

	r := c.Request()
	var l Lead

	a.Logger.Printf("webhook gained from amoCrm %s", c.FormValue("leads"))
	if r.Method != "POST" {
		return echo.ErrMethodNotAllowed
	}

	err := json.Unmarshal([]byte(c.FormValue("leads")), &l)
	if err != nil {
		return err
	}

	a.Logger.Printf("Айди лида: %v", l.Add[0].ID)
	return nil
	//var message WebhookMessage
	//decoder := json.NewDecoder(r.Body)
	//err := decoder.Decode(&message)
	//defer r.Body.Close()
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Failed to decode request body: %v", err), http.StatusBadRequest)
	//	return
	//}
	//
	//if message.Event == "new_msg" {
	//	// Обрабатываем сообщение
	//	a.Logger.Printf("Received webhook:")
	//	a.Logger.Printf("\tDescription: %s\n", message.Description)
	//	a.Logger.Printf("\tEvent: %s\n", message.Event)
	//	a.Logger.Printf("\tMessagetype: %s\n", message.MessageType)
	//	a.Logger.Printf("\tMessage: %s\n", message.Message)
	//	a.Logger.Printf("\tUserTGID: %d\n", message.UserTGId)
	//	a.Logger.Printf("\tChatTGID: %d\n", message.ChatTGId)
	//
	//	// Отправляем ответ об успешной обработке вебхука
	//	w.WriteHeader(http.StatusOK)
	//	_, _ = w.Write([]byte(`{"status": "ok"}`))
	//
	//	a.sendWebhookResult(message)
	//} else {
	//	a.Logger.Printf("Ивент вебхука: %s", message.Event)
	//}

}

func (a *App) registerHandlers() {
	a.echo.Use(middleware.Logger(), middleware.Recover())

	a.echo.Any("/webhook", a.webhookHandler)
	a.echo.Any("/webhook/amocrm", a.webhookAmoCRMHandler)

	//a.echo.Any("/int/rpc/", echo.WrapHandler(a.internalRPC))
}

// runHTTPServer is a function that starts http listener using labstack/echo.
func (a *App) runHTTPServer(host string, port int) error {
	listenAddress := fmt.Sprintf("%s:%d", host, port)
	a.Printf("starting http listener at http://%s\n", listenAddress)

	return a.echo.Start(listenAddress)
}
