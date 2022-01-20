package bot

import (
	"time"

	"github.com/Hickar/sound-seeker-bot/internal/config"
	"github.com/Hickar/sound-seeker-bot/internal/controller"
	"github.com/Hickar/sound-seeker-bot/internal/middleware/scene"
	"github.com/Hickar/sound-seeker-bot/internal/middleware/session"
	"gopkg.in/tucnak/telebot.v3"
)

type App struct {
	bot telebot.Bot
}

func Start(conf config.Config) error {
	poller := telebot.LongPoller{Timeout: time.Second * time.Duration(conf.Bot.Timeout)}

	settings := telebot.Settings{
		Token:       conf.Bot.Token,
		Poller:      &poller,
		Synchronous: !conf.Bot.Concurrent,
		Verbose:     conf.Bot.Debug,
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		return err
	}

	sessionStore := session.NewSessionStore()
	sessionMiddleware := session.Middleware(sessionStore)

	mainScene, _ := scene.NewScene(conf.Bot, "main", sessionMiddleware)
	postScene, _ := scene.NewScene(conf.Bot, "post", sessionMiddleware)

	stage := scene.NewStage()
	stage.Register(mainScene, postScene)
	stageMiddleware := scene.Middleware(stage)

	bot.Use(sessionMiddleware, stageMiddleware)
	bot.Handle("/start", func(ctx telebot.Context) error {
		return scene.GetScene(ctx).EnterScene(ctx, "main")
	})
	bot.Handle(telebot.OnText, func(context telebot.Context) error {return nil})

	if err := controller.NewMainRouter(mainScene); err != nil {
		return err
	}

	if err := controller.NewPostRouter(postScene); err != nil {
		return err
	}

	bot.Start()

	return nil
}
