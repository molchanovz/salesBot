// Code generated from jsonrpc schema by rpcgen v2.4.4; DO NOT EDIT.
package sales

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vmkteam/zenrpc/v2"
)

var (
	// Always import time package. Generated models can contain time.Time fields.
	_ time.Time
)

type Client struct {
	rpcClient *rpcClient

	Sales *svcSales
}

func NewDefaultClient(endpoint string) *Client {
	header := http.Header{}
	header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJib3RVc2VySUQiOjYsInBob25lIjoiODk5NTc5OTM5NzkiLCJleHAiOjE3MzE4NDM2NTF9.bt6_8aOOd-uV4tqjYUCv7dWjCjI_8TDFZOIw7gSIPDs")
	return NewClient(endpoint, header, &http.Client{})
}

func NewClient(endpoint string, header http.Header, httpClient *http.Client) *Client {
	c := &Client{
		rpcClient: newRPCClient(endpoint, header, httpClient),
	}

	c.Sales = newClientSales(c.rpcClient)

	return c
}

type BotUser struct {
	BotUserID      int     `json:"botUserId"`
	FirstName      *string `json:"firstName,omitempty"`
	LastName       *string `json:"lastName,omitempty"`
	Phone          string  `json:"phone"`
	StatusID       int     `json:"statusId"`
	TgidBotUser    *int    `json:"tgidBotUser,omitempty"`
	TokenExpiredAt *string `json:"tokenExpiredAt,omitempty"`
	WebhookURL     *string `json:"webhookURL,omitempty"`
}

type Chat struct {
	Admins           []int   `json:"admins"`
	ChatDescription  *string `json:"chatDescription,omitempty"`
	ChatID           int     `json:"chatId"`
	ChatName         *string `json:"chatName,omitempty"`
	ChatTitle        *string `json:"chatTitle,omitempty"`
	DailyMsg         *int    `json:"dailyMsg,omitempty"`
	Keyword          *string `json:"keyword,omitempty"`
	LinkToLinkedChat *string `json:"linkToLinkedChat,omitempty"`
	MemberCount      *int    `json:"memberCount,omitempty"`
	MsgCount         *int    `json:"msgCount,omitempty"`
	ParentTgID       int     `json:"parentTgId"`
	TgidChat         *int    `json:"tgidChat,omitempty"`
	ThemeID          int     `json:"themeId"`
}

type Message struct {
	ChatID       int     `json:"chatId"`
	ChatMemberID int     `json:"chatMemberId"`
	ContentText  *string `json:"contentText,omitempty"`
	ContentType  *string `json:"contentType,omitempty"`
	MediaURL     *string `json:"mediaURL,omitempty"`
}

type Response struct {
	Response *string `json:"response,omitempty"`
}

type SortMessage struct {
	ChatID   int       `json:"chatId"`
	Messages []Message `json:"messages"`
}

type TaskManager struct {
	CreatedAt string `json:"createdAt"`
	TaskID    int    `json:"taskId"`
	TypeID    string `json:"typeId"`
}

type svcSales struct {
	client *rpcClient
}

func newClientSales(client *rpcClient) *svcSales {
	return &svcSales{
		client: client,
	}
}

var (
	ErrSalesAddWebhookURL404 = zenrpc.NewError(404, fmt.Errorf("user was not found"))
)

// AddWebhookURL добавит или обновить ссылку на webhook
func (c *svcSales) AddWebhookURL(ctx context.Context, webhookURL string) (res *Response, err error) {
	_req := struct {
		WebhookURL string
	}{
		WebhookURL: webhookURL,
	}

	err = c.client.call(ctx, "sales.AddWebhookURL", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesAddWebhookURL404
		}
	}

	return
}

// CheckTask проверка статуса задачи с id = taskId(место в очереди/выполняется/завершена)
func (c *svcSales) CheckTask(ctx context.Context, taskId int) (res *Response, err error) {
	_req := struct {
		TaskID int
	}{
		TaskID: taskId,
	}

	err = c.client.call(ctx, "sales.CheckTask", _req, &res)

	return
}

// CheckTokenExpiration возвращает пользователю оставшийся срок годности токена. Если осталось меньше суток - вывод в часах
func (c *svcSales) CheckTokenExpiration(ctx context.Context) (res *Response, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "sales.CheckTokenExpiration", _req, &res)

	return
}

var (
	ErrSalesDeleteWebhookURL404 = zenrpc.NewError(404, fmt.Errorf("user was not found"))
)

// DeleteWebhookURL удалить ссылку на webhook
func (c *svcSales) DeleteWebhookURL(ctx context.Context) (res *Response, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "sales.DeleteWebhookURL", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesDeleteWebhookURL404
		}
	}

	return
}

var (
	ErrSalesFindChatsByTheme404 = zenrpc.NewError(404, fmt.Errorf("not found"))
)

// FindChatsByTheme поиск чатов по теме с id=themeId
func (c *svcSales) FindChatsByTheme(ctx context.Context, themeId int) (res *Response, err error) {
	_req := struct {
		ThemeID int
	}{
		ThemeID: themeId,
	}

	err = c.client.call(ctx, "sales.FindChatsByTheme", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesFindChatsByTheme404
		}
	}

	return
}

// FirstProfileRegisterUserBot отправление кода авторизации в Telegram пользователю с tgid = userTgId (взаимодействует с tgBot)
func (c *svcSales) FirstProfileRegisterUserBot(ctx context.Context, userTgId int) (res *Response, err error) {
	_req := struct {
		UserTgID int
	}{
		UserTgID: userTgId,
	}

	err = c.client.call(ctx, "sales.FirstProfileRegisterUserBot", _req, &res)

	return
}

var (
	ErrSalesGetAllChats404 = zenrpc.NewError(404, fmt.Errorf("user or chat was not found"))
)

// GetAllChats находит все чаты, в которых состоит пользователь и добавляет в БД
func (c *svcSales) GetAllChats(ctx context.Context) (res *Response, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "sales.GetAllChats", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesGetAllChats404
		}
	}

	return
}

// GetAllMessages возвращает все сообщения ОФЛАЙН МЕТОД
func (c *svcSales) GetAllMessages(ctx context.Context) (res []SortMessage, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "sales.GetAllMessages", _req, &res)

	return
}

var (
	ErrSalesGetAllProfiles404 = zenrpc.NewError(404, fmt.Errorf("users was not found"))
)

// GetAllProfiles возвращает информацию о всех доступных пользователях
func (c *svcSales) GetAllProfiles(ctx context.Context) (res []BotUser, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "sales.GetAllProfiles", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesGetAllProfiles404
		}
	}

	return
}

var (
	ErrSalesGetChatByTgID404 = zenrpc.NewError(404, fmt.Errorf("user was not found"))
)

// GetChatByTgId возвращает чат из БД по tgid ОФЛАЙН МЕТОД
func (c *svcSales) GetChatByTgID(ctx context.Context, chatTgId int) (res *Chat, err error) {
	_req := struct {
		ChatTgID int
	}{
		ChatTgID: chatTgId,
	}

	err = c.client.call(ctx, "sales.GetChatByTgId", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesGetChatByTgID404
		}
	}

	return
}

var (
	ErrSalesGetChatMessages404 = zenrpc.NewError(404, fmt.Errorf("messages was not found"))
)

// GetChatMessages возвращает сообщения из чата с id = chatId, page - страница (дефолт 0 = все страницы), pageSize - количество возвращаемых сообщений (дефолт = 25) ОФЛАЙН МЕТОД
func (c *svcSales) GetChatMessages(ctx context.Context, chatId int, page *int, pageSize *int) (res []Message, err error) {
	_req := struct {
		ChatID   int
		Page     *int
		PageSize *int
	}{
		ChatID: chatId, Page: page, PageSize: pageSize,
	}

	err = c.client.call(ctx, "sales.GetChatMessages", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesGetChatMessages404
		}
	}

	return
}

var (
	ErrSalesGetMessageByID404 = zenrpc.NewError(404, fmt.Errorf("message was not found"))
)

// GetMessageById возвращает сообщение с id = messageId ОФЛАЙН МЕТОД
func (c *svcSales) GetMessageByID(ctx context.Context, messageId int) (res *Message, err error) {
	_req := struct {
		MessageID int
	}{
		MessageID: messageId,
	}

	err = c.client.call(ctx, "sales.GetMessageById", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesGetMessageByID404
		}
	}

	return
}

// GetMyInfo возвращает информацию о пользователе ОФЛАЙН МЕТОД
func (c *svcSales) GetMyInfo(ctx context.Context) (res *BotUser, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "sales.GetMyInfo", _req, &res)

	return
}

// GetMyWebhookUrl возвращает информацию о WebhookUrl ОФЛАЙН МЕТОД
func (c *svcSales) GetMyWebhookUrl(ctx context.Context) (res *Response, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "sales.GetMyWebhookUrl", _req, &res)

	return
}

var (
	ErrSalesGetQueue404 = zenrpc.NewError(404, fmt.Errorf("user was not found"))
)

// GetQueue посмотреть очередь активных задач пользователя
func (c *svcSales) GetQueue(ctx context.Context) (res []TaskManager, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "sales.GetQueue", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesGetQueue404
		}
	}

	return
}

var (
	ErrSalesGetUserChats404 = zenrpc.NewError(404, fmt.Errorf("chats was not found"))
)

// GetUserChats возвращает информацию о чатах, в которых состоит пользователь ОФЛАЙН МЕТОД
func (c *svcSales) GetUserChats(ctx context.Context) (res []Chat, err error) {
	_req := struct {
	}{}

	err = c.client.call(ctx, "sales.GetUserChats", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesGetUserChats404
		}
	}

	return
}

// PauseParsingTask приостанавливает задачу (id = taskID) парсинга чата.
func (c *svcSales) PauseParsingTask(ctx context.Context, taskId int) (res *Response, err error) {
	_req := struct {
		TaskID int
	}{
		TaskID: taskId,
	}

	err = c.client.call(ctx, "sales.PauseParsingTask", _req, &res)

	return
}

// ProfileAddUserBot добавляет нового пользователя в базу данных (взаимодействует с tgBot)
func (c *svcSales) ProfileAddUserBot(ctx context.Context, userTgId int, userFirstName string, userLastName string, phone string) (res *BotUser, err error) {
	_req := struct {
		UserTgID      int
		UserFirstName string
		UserLastName  string
		Phone         string
	}{
		UserTgID: userTgId, UserFirstName: userFirstName, UserLastName: userLastName, Phone: phone,
	}

	err = c.client.call(ctx, "sales.ProfileAddUserBot", _req, &res)

	return
}

// ProfileAddUserBotToken добавляет в бд срок годности токена и возвращает токен пользователю. Подходит так же при обновлении токена(в случае продления подписки).
func (c *svcSales) ProfileAddUserBotToken(ctx context.Context, userTgId int) (res *Response, err error) {
	_req := struct {
		UserTgID int
	}{
		UserTgID: userTgId,
	}

	err = c.client.call(ctx, "sales.ProfileAddUserBotToken", _req, &res)

	return
}

// ProfileRegisterUserBot окончание авторизации пользователя. Ввод tgid = userTgId (регестрируемого пользователя) и кода авторизации из Telegram = code(взаимодействует с tgBot)
func (c *svcSales) ProfileRegisterUserBot(ctx context.Context, userTgId int, code string) (res *Response, err error) {
	_req := struct {
		UserTgID int
		Code     string
	}{
		UserTgID: userTgId, Code: code,
	}

	err = c.client.call(ctx, "sales.ProfileRegisterUserBot", _req, &res)

	return
}

// RestartParsingTask возобновляет выполнение задачи парсинга чата (снимает с паузы)
func (c *svcSales) RestartParsingTask(ctx context.Context, taskId int) (res *Response, err error) {
	_req := struct {
		TaskID int
	}{
		TaskID: taskId,
	}

	err = c.client.call(ctx, "sales.RestartParsingTask", _req, &res)

	return
}

// RestoreToken восстанавливает токен при потере
func (c *svcSales) RestoreToken(ctx context.Context, userTgId int) (res *Response, err error) {
	_req := struct {
		UserTgID int
	}{
		UserTgID: userTgId,
	}

	err = c.client.call(ctx, "sales.RestoreToken", _req, &res)

	return
}

var (
	ErrSalesSendTextMessage404 = zenrpc.NewError(404, fmt.Errorf("user or chat was not found"))
)

// SendTextMessage отправляем текстовое сообщение в чат с id = chatId
func (c *svcSales) SendTextMessage(ctx context.Context, chatId int, message string) (res *Response, err error) {
	_req := struct {
		ChatID  int
		Message string
	}{
		ChatID: chatId, Message: message,
	}

	err = c.client.call(ctx, "sales.SendTextMessage", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesSendTextMessage404
		}
	}

	return
}

var (
	ErrSalesSendTextMessageByTgChatID404 = zenrpc.NewError(404, fmt.Errorf("user or chat was not found"))
)

// SendTextMessageByTgChatId отправляем текстовое сообщение в чат с id = chatId (имеется в виду tgId чата)
func (c *svcSales) SendTextMessageByTgChatID(ctx context.Context, chatId int, message string) (res *Response, err error) {
	_req := struct {
		ChatID  int
		Message string
	}{
		ChatID: chatId, Message: message,
	}

	err = c.client.call(ctx, "sales.SendTextMessageByTgChatId", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesSendTextMessageByTgChatID404
		}
	}

	return
}

var (
	ErrSalesStartParseChatByID404 = zenrpc.NewError(404, fmt.Errorf("user or chat was not found"))
)

// StartParseChatById начинает парсинг чата с id = chatId
func (c *svcSales) StartParseChatByID(ctx context.Context, chatId int) (res *Response, err error) {
	_req := struct {
		ChatID int
	}{
		ChatID: chatId,
	}

	err = c.client.call(ctx, "sales.StartParseChatById", _req, &res)

	switch v := err.(type) {
	case *zenrpc.Error:
		if v.Code == 404 {
			err = ErrSalesStartParseChatByID404
		}
	}

	return
}

type rpcClient struct {
	endpoint string
	cl       *http.Client

	requestID uint64
	header    http.Header
}

func newRPCClient(endpoint string, header http.Header, httpClient *http.Client) *rpcClient {
	return &rpcClient{
		endpoint: endpoint,
		header:   header,
		cl:       httpClient,
	}
}

func (rc *rpcClient) call(ctx context.Context, methodName string, request, result interface{}) error {
	// encode params
	bts, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("encode params: %w", err)
	}

	requestID := atomic.AddUint64(&rc.requestID, 1)
	requestIDBts := json.RawMessage(strconv.Itoa(int(requestID)))

	req := zenrpc.Request{
		Version: zenrpc.Version,
		ID:      &requestIDBts,
		Method:  methodName,
		Params:  bts,
	}

	res, err := rc.Exec(ctx, req)
	if err != nil {
		return err
	}

	if res == nil {
		return nil
	}

	if res.Error != nil {
		return res.Error
	}

	if res.Result == nil {
		return nil
	}

	if result == nil {
		return nil
	}

	return json.Unmarshal(*res.Result, result)
}

// Exec makes http request to jsonrpc endpoint and returns json rpc response.
func (rc *rpcClient) Exec(ctx context.Context, rpcReq zenrpc.Request) (*zenrpc.Response, error) {
	if n, ok := ctx.Value("JSONRPC2-Notification").(bool); ok && n {
		rpcReq.ID = nil
	}

	c, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("json marshal call failed: %w", err)
	}

	buf := bytes.NewReader(c)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rc.endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header = rc.header.Clone()
	req.Header.Add("Content-Type", "application/json")

	if xRequestID, ok := ctx.Value("X-Request-Id").(string); ok && req.Header.Get("X-Request-Id") == "" && xRequestID != "" {
		req.Header.Add("X-Request-Id", xRequestID)
	}

	// Do request
	resp, err := rc.cl.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, fmt.Errorf("make request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response (%d)", resp.StatusCode)
	}

	var zresp zenrpc.Response
	if rpcReq.ID == nil {
		return &zresp, nil
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body (%s) read failed: %w", bb, err)
	}

	if err = json.Unmarshal(bb, &zresp); err != nil {
		return nil, fmt.Errorf("json decode failed (%s): %w", bb, err)
	}

	return &zresp, nil
}
