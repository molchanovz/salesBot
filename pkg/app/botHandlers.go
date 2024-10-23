package app

import (
	sales "apisrv/pkg/client"
	"apisrv/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
	"strings"
)

const (
	somePattern               = "/some"
	CallBackPatternAgreement  = "agree_"
	CallBackPatternRefusement = "refuse_"
)

func (a *App) registerBotHandlers() {
	a.b.RegisterHandler(bot.HandlerTypeMessageText, somePattern, bot.MatchTypePrefix, a.someHandler)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, CallBackPatternAgreement, bot.MatchTypePrefix, a.handleAgree)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, CallBackPatternRefusement, bot.MatchTypePrefix, a.handleRefuse)

}

func (a *App) someHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := int(update.Message.From.ID)

	req := update.Message.Text[len(somePattern)+1:]
	if req == "" {
		a.Logger.Errorf("Запрос пустой")
		return
	}

	a.processGigachatAnswer(ctx, b, req, chatId)
}

func (a *App) handleAgree(ctx context.Context, b *bot.Bot, update *models.Update) {
	callBackData := update.CallbackQuery.Data[len(CallBackPatternAgreement):]
	params, err := NewCallbackDataParams(callBackData)
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	a.Logger.Printf("Наш юзер %d", params.TgID)

	c := sales.NewDefaultClient("http://91.222.239.37:8080/v1/rpc/", a.cfg.Client.Token)

	message, err := a.sr.GigachatMessageByID(ctx, params.MessageId)
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	if message == nil {
		a.Logger.Errorf("empty")
		return
	}

	info, err := c.Sales.SendTextMessageByTgChatID(ctx, params.TgID, message.Message)
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	fmt.Printf("Инфо: %+v", info)
}

func (a *App) handleRefuse(ctx context.Context, b *bot.Bot, update *models.Update) {
	callBackData := update.CallbackQuery.Data[len(CallBackPatternRefusement):]
	params, err := NewCallbackDataParams(callBackData)
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	message, err := a.sr.GigachatMessageByID(ctx, params.MessageId)
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	if message == nil {
		a.Logger.Errorf("empty")
		return
	}

	cmdForEdit := "Новое сообщение отправьте /edit" + strconv.Itoa(params.MessageId)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.CallbackQuery.From.ID, Text: cmdForEdit, ParseMode: models.ParseModeMarkdown})
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

}

/*
Функция обращения к API gigaChat
*/
func (a *App) processGigachatAnswer(ctx context.Context, b *bot.Bot, text string, chatId int) {

	fmt.Println(chatId)

	str, err := a.g.SendRequest(text)
	if err != nil {
		a.Logger.Errorf("%v", err)
	}

	res := strings.Builder{}

	res.WriteString("Начальный запрос\n\n```" +
		text +
		"```\n\nСгенерированный ответ\n\n")

	var generatedText string
	for _, content := range str.Choices {
		generatedText += content.Message.Content + " "
	}

	res.WriteString("```" + generatedText + "```\n\nОтправить ответ?")

	var buttons [][]models.InlineKeyboardButton
	var agreementButtons []models.InlineKeyboardButton

	message, err := a.sr.AddGigachatMessage(ctx, &db.GigachatMessage{Message: generatedText})
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	if message == nil {
		println("empty msg")
		return
	}

	jsonParams := "{  \"messageId\": " + strconv.Itoa(message.ID) + ",\n  \"TgID\": " + strconv.Itoa(chatId) + "\n}"

	fmt.Println(jsonParams)

	agreementButtons = append(agreementButtons, models.InlineKeyboardButton{
		Text: "Да", CallbackData: CallBackPatternAgreement + jsonParams},
	)

	agreementButtons = append(agreementButtons, models.InlineKeyboardButton{
		Text: "Нет", CallbackData: CallBackPatternRefusement + jsonParams},
	)

	fmt.Println(agreementButtons[0].CallbackData)

	buttons = append(buttons, agreementButtons)

	markup := models.InlineKeyboardMarkup{InlineKeyboard: buttons}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: a.cfg.Bot.MainUserId, Text: res.String(), ReplyMarkup: markup, ParseMode: models.ParseModeMarkdown})
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

}

func (a App) sendWebhookResult(message WebhookMessage) {
	ctx := context.Background()
	if strings.Contains(strings.ToLower(message.Message), "двер") {

		// TODO принимать информацию контакта для a.crm.AddContact()

		a.processGigachatAnswer(ctx, a.b, message.Message, message.ChatTGId)
	}
}

type CallbackDataParams struct {
	MessageId int `json:"messageId"`
	TgID      int `json:"TgID"`
}

func NewCallbackDataParams(s string) (CallbackDataParams, error) {
	var b CallbackDataParams
	err := json.Unmarshal([]byte(s), &b)
	return b, err
}

func (b CallbackDataParams) String() (string, error) {
	s, err := json.Marshal(b)
	return string(s), err
}
