package app

import (
	"apisrv/pkg/amoCRM"
	"apisrv/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
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

func (a App) webhookAmoCRMHandler(c echo.Context) error {
	r := c.Request()
	a.Logger.Printf("webhook gained from amoCrm %+v", c.Request().PostForm)
	if r.Method != "POST" {
		return echo.ErrMethodNotAllowed
	}

	leadId, _ := strconv.Atoi(c.Request().PostForm.Get("leads[add][0][id]"))

	a.Logger.Printf("leadId: %v", leadId)

	leadString := a.crm.GetLead(a.crm.Token, leadId)

	var lead amoCRM.Lead

	json.Unmarshal(leadString, &lead)

	contactId := lead.Embedded.Contacts[0].Id

	a.Logger.Printf("contactId: %v", contactId)

	var tgId int64
	var err error
	for i := range lead.Embedded.Contacts[0].CustomFieldsValues {
		if lead.Embedded.Contacts[0].CustomFieldsValues[i].FieldId == 396043 {
			value := lead.Embedded.Contacts[0].CustomFieldsValues[i].Values[0].Value
			tgId, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				panic(err)
			}
		}
	}

	a.Logger.Printf("Наш контакт айди: %v", contactId)
	a.Logger.Printf("Наш тг айди: %v", tgId)

	ctx := context.Background()
	var message []db.Gigachatmessage
	message, _ = a.sr.GigachatmessagesByFilters(ctx, &db.GigachatmessageSearch{Tgid: &tgId}, db.PagerOne)

	a.crm.EditContact(contactId, strconv.FormatInt(*message[0].Tgid, 10))

	return nil

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
