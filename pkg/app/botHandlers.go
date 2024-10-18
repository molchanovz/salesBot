package app

import (
	"context"
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
	userIdStr := update.CallbackQuery.Data[len(CallBackPatternAgreement):]
	userId, _ := strconv.Atoi(userIdStr)
	fmt.Printf("Наш юзер %d", userId)

	//c := sales.NewDefaultClient("http://91.222.239.37:8080/v1/rpc/")
	//info, err := c.Sales.SendTextMessageByTgChatID(ctx, userId, "Ответное сообщение")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Printf("Инфо: %+v", info)
}

/*
Функция обращения к API gigaChat
*/
func (a *App) processGigachatAnswer(ctx context.Context, b *bot.Bot, text string, chatId int) {

	str, err := a.g.SendRequest(text)
	if err != nil {
		a.Logger.Errorf("%v", err)
	}

	contentString := "Начальный запрос\n\n" +
		text +
		"\n\nСгенерированный ответ\n\n"

	for _, content := range str.Choices {
		contentString = fmt.Sprint(contentString, content.Message.Content, " ")
	}

	contentString = fmt.Sprint(contentString, "\n\n", "Отправить ответ?")

	var buttons [][]models.InlineKeyboardButton
	var agreementButtons []models.InlineKeyboardButton

	agreementButtons = append(agreementButtons, models.InlineKeyboardButton{
		Text: "Да", CallbackData: "agree_" + strconv.Itoa(chatId)},
	)

	agreementButtons = append(agreementButtons, models.InlineKeyboardButton{
		Text: "Нет", CallbackData: "refuse"},
	)

	buttons = append(buttons, agreementButtons)

	markup := models.InlineKeyboardMarkup{InlineKeyboard: buttons}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: chatId, Text: contentString, ReplyMarkup: markup})
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}
}

func (a App) sendWebhookResult(text string, chatId int) {
	ctx := context.Background()
	a.processGigachatAnswer(ctx, a.b, text, a.cfg.Bot.MainUserId)
}
