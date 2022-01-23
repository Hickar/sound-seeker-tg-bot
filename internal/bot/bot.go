package bot

import (
	"time"

	"github.com/Hickar/sound-seeker-bot/internal/config"
	"github.com/Hickar/sound-seeker-bot/internal/controller"
	ssnStore "github.com/Hickar/sound-seeker-bot/internal/session_store"
	"github.com/Hickar/sound-seeker-bot/internal/usecase"
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/scene"
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/session"
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

	sessionStore, err := ssnStore.NewSessionStore(conf.Redis)
	if err != nil {
		return err
	}

	sessionMiddleware := session.Middleware(sessionStore)

	mainScene, _ := scene.NewScene(conf.Bot, controller.SceneMainName, sessionMiddleware)
	postScene, _ := scene.NewScene(conf.Bot, controller.ScenePostName, sessionMiddleware)

	stage := scene.NewStage()
	stage.Register(mainScene, postScene)
	stageMiddleware := scene.Middleware(stage)

	bot.Use(sessionMiddleware, stageMiddleware)
	bot.Handle("/start", func(ctx telebot.Context) error {
		return stage.EnterScene(ctx, controller.SceneMainName)
	})

	// Little hack to make bot pass updates of "text" type to scenes underneath :P
	bot.Handle(telebot.OnText, func(context telebot.Context) error { return nil })

	mainController := controller.NewMainController(stage)
	controller.NewMainRouter(mainScene, mainController)

	postController := controller.NewChannelPostController(&usecase.PostUsecase{}, stage)
	controller.NewPostRouter(postScene, postController)

	bot.Start()

	return nil
}
