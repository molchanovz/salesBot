package app

import (
	"apisrv/pkg/gigaChat"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	somePattern = "/some"
)

func (a *App) registerBotHandlers() {
	a.b.RegisterHandler(bot.HandlerTypeMessageText, somePattern, bot.MatchTypePrefix, a.someHandler)
}

func (a *App) someHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := int(update.Message.From.ID)

	req := update.Message.Text[len(somePattern)+1:]
	if req == "" {
		fmt.Println("Запрос пустой")
		return
	}

	str, err := gigaChat.SendRequest(req)
	if err != nil {
		fmt.Println(err)
	}

	contentString := "Начальный запрос\n\n" +
		req +
		"\n\nСгенерированный ответ\n\n"

	for _, content := range str.Choices {
		contentString = fmt.Sprint(contentString, content.Message.Content, " ")
	}

	contentString = fmt.Sprint(contentString, "\n\n", "Отправить ответ?")

	var buttons [][]models.InlineKeyboardButton
	var agreementButtons []models.InlineKeyboardButton

	agreementButtons = append(agreementButtons, models.InlineKeyboardButton{
		Text: "Да", CallbackData: "someData"},
	)

	agreementButtons = append(agreementButtons, models.InlineKeyboardButton{
		Text: "Нет", CallbackData: "someData"},
	)

	buttons = append(buttons, agreementButtons)

	markup := models.InlineKeyboardMarkup{InlineKeyboard: buttons}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatId, Text: contentString, ReplyMarkup: markup})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func pointer[T any](in T) *T {
	return &in
}
