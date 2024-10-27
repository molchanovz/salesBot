package app

import (
	"apisrv/pkg/amoCRM"
	"apisrv/pkg/db"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
)

const (
	fieldTgId = 396043
	leadIdKey = "leads[add][0][id]"
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
	params, err := c.FormParams()
	if err != nil {
		a.Logger.Errorf("%+v", err)
		return err
	}
	a.Logger.Printf("webhook gained from amoCrm %+v\n", params)
	if r.Method != "POST" {
		return echo.ErrMethodNotAllowed
	}

	var leadId int
	leadId, _ = strconv.Atoi(params.Get(leadIdKey))
	if leadId == 0 {
		return errors.New("ERROR: leadId is 0")
	}
	a.Logger.Printf("leadId: %v\n", leadId)

	leadString := a.crm.GetLead(a.crm.Token, leadId)
	var lead amoCRM.Lead
	err = json.Unmarshal(leadString, &lead)
	if err != nil {
		return err
	}

	var contactId int
	contactId = lead.Embedded.Contacts[0].Id
	if contactId == 0 {
		return errors.New("ERROR: contactId is 0")
	}
	a.Logger.Printf("contactId: %v\n", contactId)

	var tgId int64
	var contact amoCRM.Contact
	err = json.Unmarshal(a.crm.GetContact(a.crm.Token, contactId), &contact)
	if err != nil {
		return err
	}
	//поиск поля с tgId (смотрим в сделке)
	for _, values := range contact.CustomFieldsValues {
		if values.FieldId == fieldTgId {
			tgId, _ = strconv.ParseInt(values.Values[0].Value, 10, 64)
		}
	}

	if tgId == 0 {
		return errors.New("ERROR: tgId is 0")
	}
	a.Logger.Printf("tgId: %v\n", tgId)

	ctx := context.Background()
	var messages []db.Gigachatmessage
	messages, _ = a.sr.GigachatmessagesByFilters(ctx, &db.GigachatmessageSearch{Tgid: &tgId}, db.PagerOne)

	if len(messages) < 1 {
		return echo.ErrNotFound
	}

	response, err := a.crm.EditContact(contactId, messages[0].Request)
	if err != nil {
		return err
	}

	a.Logger.Printf("Editing contact: %v\n", response)
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
