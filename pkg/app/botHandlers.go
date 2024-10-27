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
	editMessagePattern        = "/edit"
	CallBackPatternAgreement  = "agree_"
	CallBackPatternRefusement = "refuse_"
)

func (a *App) registerBotHandlers() {
	a.b.RegisterHandler(bot.HandlerTypeMessageText, editMessagePattern, bot.MatchTypePrefix, a.editMessageHandler)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, CallBackPatternAgreement, bot.MatchTypePrefix, a.handleAgree)
	a.b.RegisterHandler(bot.HandlerTypeCallbackQueryData, CallBackPatternRefusement, bot.MatchTypePrefix, a.handleRefuse)
}

// Handler с editMessagePattern, редактирование сгенерированного сообщения
func (a *App) editMessageHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	req := update.Message.Text
	if req == "" {
		a.Logger.Errorf("Запрос пустой")
		return
	}

	parts := strings.SplitN(req, " ", 2)
	firstPart := parts[0]
	handlerMsgId := strings.SplitN(firstPart, "_", 2)
	msgId, _ := strconv.Atoi(handlerMsgId[1])

	request := parts[1]

	c := sales.NewDefaultClient("http://91.222.239.37:8080/v1/rpc/", a.cfg.Client.Token)

	message, err := a.sr.GigachatmessageByID(ctx, msgId)
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	if message == nil {
		a.Logger.Errorf("empty")
		return
	}

	info, err := c.Sales.SendTextMessageByTgChatID(ctx, int(*message.Tgid), request, pointer(false))
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	fmt.Printf("Инфо при отправке сообщения: %+v", info)
}

// Handler нажатия кнопки, соглашение с отправкой сгенерированного сообщения
func (a *App) handleAgree(ctx context.Context, b *bot.Bot, update *models.Update) {
	callBackData := update.CallbackQuery.Data[len(CallBackPatternAgreement):]
	params, err := NewCallbackDataParams(callBackData)
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	a.Logger.Printf("Наш юзер %d", params.TgID)

	c := sales.NewDefaultClient("http://91.222.239.37:8080/v1/rpc/", a.cfg.Client.Token)

	message, err := a.sr.GigachatmessageByID(ctx, params.MessageId)
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	if message == nil {
		a.Logger.Errorf("empty")
		return
	}

	info, err := c.Sales.SendTextMessageByTgChatID(ctx, params.TgID, message.Message, pointer(false))
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	fmt.Printf("Инфо: %+v", info)
}

// Handler нажатия кнопки, отказ от отправки сгенерированного сообщения
func (a *App) handleRefuse(ctx context.Context, b *bot.Bot, update *models.Update) {
	callBackData := update.CallbackQuery.Data[len(CallBackPatternRefusement):]
	params, err := NewCallbackDataParams(callBackData)
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	message, err := a.sr.GigachatmessageByID(ctx, params.MessageId)
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	if message == nil {
		a.Logger.Errorf("empty")
		return
	}

	var cmdForEdit strings.Builder
	cmdForEdit.WriteString("Чтобы изменить сообщение скопируйте текст ниже, напишите новый ответ и отправьте в чат\n")
	cmdForEdit.WriteString("`/edit_" + strconv.Itoa(params.MessageId) + " `")

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.CallbackQuery.From.ID, Text: cmdForEdit.String(), ParseMode: models.ParseModeMarkdown})
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

}

// Функция обращения к API gigaChat
func (a *App) processGigachatAnswer(ctx context.Context, b *bot.Bot, text string, chatId int64) {

	//Проверка есть ли сообщение от этого пользователя TODO
	message, _ := a.sr.OneGigachatmessage(ctx, &db.GigachatmessageSearch{Tgid: &chatId})
	if message != nil {
		a.Logger.Errorf("От этого человека уже было сообщение, пропускаем")
		return
	}

	yesNoPrompt := "Наша компания продаёт двери. \n\nМожно ли человеку предлагать купить дверь если он написал такое сообщение? Мы должны предлагать купить дверь, только если видим в запросе явную нужду в покупке/замене двери. Если из контекста кажется, что дверь можно предложить, но в сообщении явно не написано о желании купить/заменить дверь, то отвечай нет. В сообщении должно быть явно указано, что речь идет об обновлении, замене или покупке чего-то, связанного с дверями. В ответ напиши одно слово, \"да\" либо \"нет\" без точек на конце.\nВот сообщение:\n"

	resp, err := a.g.SendRequest(yesNoPrompt + text)
	if err != nil {
		a.Logger.Errorf("%v", err)
	}

	if len(resp.Choices) < 1 {
		a.Logger.Errorf("Пустой ответ: %+v\nЗапрос:\n%s", resp.Choices, text)
		return
	}

	switch resp.Choices[0].Message.Content {
	case "Да", "да", "да.":
		advertisePrompt := "Ты - обычный сосед, который хочет по-соседски посоветовать установку дверей от \"Синдорсервис\".  Ты должен сообщить, что знаешь владельца, сам пользовался и человеку могут предложить скидку, если он напишет, что он \"от Виктора\". Твоё сообщение должно быть похоже на повседневное человеческое, без лишней дружелюбности. Старайся писать максимально кратко, не более 200 символов. В твоем сообщении должен учитываться контекст сообщения от человека, вычлени из его сообщения жалобу или проблему, если видишь ее, и опиши, как бы продукт компании мог с ней помочь\n\nВот сообщение от человека:\n"
		resp, err = a.g.SendRequest(advertisePrompt + text)
		if err != nil {
			a.Logger.Errorf("%v", err)
		}
	case "Нет", "нет", "нет.":
		// TODO: для дебага
		b.SendMessage(ctx, &bot.SendMessageParams{
			ParseMode: models.ParseModeMarkdown,
			Text:      resp.Choices[0].Message.Content + "\n запрос:\n```скопировать_сообщение " + text + "```",
			ChatID:    a.cfg.Bot.MainUserId,
		})
		return
	default:
		a.Logger.Errorf("Некорректный ответ бота на запрос! Запрос: %s\n, Ответ: %s", text, resp.Choices[0].Message.Content)
		return
	}

	res := strings.Builder{}

	res.WriteString("Начальный запрос\n\n```Скопировать " +
		text +
		"```\n\nСгенерированный ответ\n\n")

	var generatedText string
	for _, content := range resp.Choices {
		generatedText += content.Message.Content + " "
	}

	res.WriteString("```Скопировать " + generatedText + "```\n\nОтправить ответ?")

	var buttons [][]models.InlineKeyboardButton
	var agreementButtons []models.InlineKeyboardButton

	//Добавление сообщения в БД
	newMessage, err := a.sr.AddGigachatmessage(ctx, &db.Gigachatmessage{Message: generatedText, Tgid: &chatId, Request: text})
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}

	if newMessage == nil {
		a.Logger.Errorf("empty msg")
		return
	}

	jsonParams := "{  \"messageId\": " + strconv.Itoa(newMessage.ID) + ",\n  \"TgID\": " + strconv.Itoa(int(chatId)) + "\n}"

	agreementButtons = append(agreementButtons, models.InlineKeyboardButton{
		Text: "Да", CallbackData: CallBackPatternAgreement + jsonParams},
	)

	agreementButtons = append(agreementButtons, models.InlineKeyboardButton{
		Text: "Изменить", CallbackData: CallBackPatternRefusement + jsonParams},
	)

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
	// TODO: норм обработка
	if strings.Contains(strings.ToLower(message.Message), "двер") {
		a.processGigachatAnswer(ctx, a.b, message.Message, *message.SenderTgId)
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

func pointer[T any](in T) *T {
	return &in
}
