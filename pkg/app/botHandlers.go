package app

import (
	sales "apisrv/pkg/client"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
)

const (
	somePattern              = "/some"
	CallBackPatternAgreement = "agree_"
)

func (a *App) registerBotHandlers() {
	a.b.RegisterHandler(bot.HandlerTypeMessageText, somePattern, bot.MatchTypePrefix, a.someHandler)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, CallBackPatternAgreement, bot.MatchTypePrefix, a.handleInfo)

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

func (a *App) handleInfo(ctx context.Context, b *bot.Bot, update *models.Update) {
	callBackData := update.CallbackQuery.Data[len(CallBackPatternAgreement):]
	params, err := NewCallbackDataParams(callBackData)
	if err != nil {
		return
	}

	fmt.Printf("Наш юзер %d", params.TgID)

	c := sales.NewDefaultClient("http://91.222.239.37:8080/v1/rpc/")
	info, err := c.Sales.SendTextMessageByTgChatID(ctx, params.TgID, params.Text)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Инфо: %+v", info)
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

	contentString := "Начальный запрос\n\n" +
		text +
		"\n\nСгенерированный ответ\n\n"

	var generatedText string

	for _, content := range str.Choices {
		generatedText = fmt.Sprint(generatedText, content.Message.Content, " ")
	}

	contentString = fmt.Sprint(generatedText, "\n\n", "Отправить ответ?")

	var buttons [][]models.InlineKeyboardButton
	var agreementButtons []models.InlineKeyboardButton

	jsonParams := "{  \"text\": \"" + generatedText + "\",\n  \"TgID\": " + strconv.Itoa(chatId) + "\n}"

	fmt.Println(jsonParams)

	agreementButtons = append(agreementButtons, models.InlineKeyboardButton{
		Text: "Да", CallbackData: "agree_" + jsonParams},
	)

	agreementButtons = append(agreementButtons, models.InlineKeyboardButton{
		Text: "Нет", CallbackData: "refuse"},
	)

	buttons = append(buttons, agreementButtons)

	markup := models.InlineKeyboardMarkup{InlineKeyboard: buttons}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: a.cfg.Bot.MainUserId, Text: contentString, ReplyMarkup: markup})
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}
}

func (a App) sendWebhookResult(message WebhookMessage) {

	ctx := context.Background()
	a.processGigachatAnswer(ctx, a.b, message.Message, message.ChatTGId)
}

type CallbackDataParams struct {
	Text string `json:"text"`
	TgID int    `json:"TgID"`
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
