package app

import (
	"apisrv/pkg/gigaChat"
	"context"
	"github.com/go-telegram/bot"
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
	}
	GigaChat gigaChat.GigaChatConfig
	Database *pg.Options
}

type App struct {
	embedlog.Logger
	appName    string
	cfg        Config
	db         db.DB
	dbc        *pg.DB
	commonRepo db.CommonRepo
	echo       *echo.Echo
	vtsrv      zenrpc.Server
	b          *bot.Bot
	g          *gigaChat.GigaChat
}

func New(appName string, verbose bool, cfg Config, dbo db.DB, dbc *pg.DB) *App {
	a := &App{
		appName:    appName,
		cfg:        cfg,
		db:         dbo,
		dbc:        dbc,
		commonRepo: db.NewCommonRepo(dbo),
		echo:       echo.New(),
	}
	a.SetStdLoggers(verbose)
	a.echo.HideBanner = true
	a.echo.HidePort = true
	a.echo.IPExtractor = echo.ExtractIPFromRealIPHeader()
	opts := []bot.Option{}
	a.b, _ = bot.New(cfg.Bot.Token, opts...)
	a.g = gigaChat.NewGigaChat(cfg.GigaChat)
	//a.vtsrv = vt.New(a.db, a.Logger, a.cfg.Server.IsDevel)

	return a
}

// Run is a function that runs application.
func (a *App) Run() error {
	a.registerMetrics()
	a.registerHttpHandlers()
	a.registerBotHandlers()
	a.b.Start(context.TODO())
	//return a.runHTTPServer(a.cfg.Server.Host, a.cfg.Server.Port)
	return nil
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
