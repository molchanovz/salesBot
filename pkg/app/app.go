package app

import (
	"apisrv/pkg/amoCRM"
	"apisrv/pkg/gigaChat"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"time"

	"apisrv/pkg/db"
	"apisrv/pkg/embedlog"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/vmkteam/zenrpc/v2"
)

type Config struct {
	Bot struct {
		Token         string
		SupportChatId string
		MainUserId    int
	}
	Client struct {
		Token string
	}
	GigaChat gigaChat.GigaChatConfig
	AmoCRM   amoCRM.AmoCRMConfig
	Database *pg.Options
}

type App struct {
	embedlog.Logger
	appName string
	cfg     Config
	db      db.DB
	dbc     *pg.DB
	sr      db.SalesbotRepo
	echo    *echo.Echo
	vtsrv   zenrpc.Server
	b       *bot.Bot
	g       *gigaChat.GigaChat
	crm     *amoCRM.AmoCRM
}

func New(appName string, verbose bool, cfg Config, dbo db.DB, dbc *pg.DB) *App {
	var err error
	a := &App{
		appName: appName,
		cfg:     cfg,
		db:      dbo,
		dbc:     dbc,
		sr:      db.NewSalesbotRepo(dbo),
		echo:    echo.New(),
	}
	a.SetStdLoggers(verbose)
	a.echo.HideBanner = true
	a.echo.HidePort = true
	a.echo.IPExtractor = echo.ExtractIPFromRealIPHeader()
	opts := []bot.Option{
		bot.WithDefaultHandler(a.defaultHandler),
	}
	a.b, err = bot.New(cfg.Bot.Token, opts...)
	if err != nil {
		a.Logger.Errorf("%v", err)
	}
	a.g = gigaChat.NewGigaChat(cfg.GigaChat)
	a.crm = amoCRM.NewAmoCRM(cfg.AmoCRM)
	//a.vtsrv = vt.New(a.db, a.Logger, a.cfg.Server.IsDevel)

	return a
}

// Run is a function that runs application.
func (a *App) Run() error {
	a.registerMetrics()
	a.registerBotHandlers()
	a.registerHandlers()
	go a.b.Start(context.TODO())
	//return a.runHTTPServer(a.cfg.Server.Host, a.cfg.Server.Port)
	return a.runHTTPServer("0.0.0.0", 8078)
}

//// VTTypeScriptClient returns TypeScript client for VT.
//func (a *App) VTTypeScriptClient() ([]byte, error) {
//	gen := rpcgen.FromSMD(a.vtsrv.SMD())
//	tsSettings := typescript.Settings{ExcludedNamespace: []string{NSVFS}, WithClasses: true}
//	return gen.TSCustomClient(tsSettings).Generate()
//}

// Shutdown is a function that gracefully stops HTTP server.
func (a *App) Shutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := a.echo.Shutdown(ctx); err != nil {
		a.Errorf("shutting down server err=%q", err)
	}
}

func (a *App) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}
	req := "Наша компания продаёт двери. \n\nМожно ли человеку предлагать купить дверь если он написал такое сообщение? Мы должны предлагать купить дверь, только если видим в запросе явную нужду в покупке/замене двери. Если из контекста кажется, что дверь можно предложить, но в сообщении явно не написано о желании купить/заменить дверь, то отвечай нет. В сообщении должно быть явно указано, что речь идет об обновлении, замене или покупке чего-то, связанного с дверями. В ответ напиши одно слово, \"да\" либо \"нет\" без точек на конце.\nВот сообщение:\n"
	resp, err := a.g.SendRequest(req + update.Message.Text)
	if err != nil {
		a.Logger.Errorf("%v", err)
	}

	if len(resp.Choices) < 1 {
		a.Logger.Errorf("Пустой ответ: %+v", resp.Choices)
		return
	}

	switch resp.Choices[0].Message.Content {
	case "Да", "да":
		req = "Ты - обычный сосед, который хочет по-дружески посоветовать установку дверей от ООО \"ДВ1\".  Ты должен сообщить, что знаешь владельца, сам пользовался и человеку могут предложить скидку, если он напишет, что он \"от Виктора\". Твоё сообщение должно быть похоже на повседневное человеческое, без лишней дружелюбности. Старайся писать максимально кратко, не более 200 символов. В твоем сообщении должен учитываться контекст сообщения от человека, вычлени из его сообщения жалобу или проблему, если видишь ее, и опиши, как бы продукт компании мог с ней помочь\n\nВот сообщение от человека:\n"
		resp, err = a.g.SendRequest(req + update.Message.Text)
		if err != nil {
			a.Logger.Errorf("%v", err)
		}
	case "Нет", "нет":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ParseMode: models.ParseModeMarkdown,
			Text:      resp.Choices[0].Message.Content + "\n запрос:\n```скопировать_сообщение " + update.Message.Text + "```",
			ChatID:    update.Message.Chat.ID,
		})
		return
	default:
		a.Logger.Errorf("Некорректный ответ бота на запрос! Запрос: %s\n, Ответ: %s", update.Message.Text, resp.Choices[0].Message.Content)
		return
	}

	//res := strings.Builder{}
	//
	//res.WriteString("Запрос:\n\n" +
	//	update.Message.Text)
	//
	//var generatedText string
	//for _, content := range resp.Choices {
	//	generatedText += content.Message.Content
	//}
	//
	//res.WriteString("\n\nОтвет:\n\n" + generatedText)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		Text:   resp.Choices[0].Message.Content,
		ChatID: update.Message.Chat.ID,
	})
	if err != nil {
		a.Logger.Errorf("%v", err)
		return
	}
}
