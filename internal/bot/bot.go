package bot

import (
	"net/http"
	"time"

	"github.com/Hickar/sound-seeker-bot/internal/config"
	"github.com/Hickar/sound-seeker-bot/internal/controller"
	"github.com/Hickar/sound-seeker-bot/internal/repository"
	localDatasource "github.com/Hickar/sound-seeker-bot/internal/repository/local_datasource"
	remoteDatasource "github.com/Hickar/sound-seeker-bot/internal/repository/remote_datasource"
	"github.com/Hickar/sound-seeker-bot/internal/usecase"
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/scene"
	"github.com/Hickar/sound-seeker-bot/pkg/middleware/session"
	"github.com/Hickar/sound-seeker-bot/pkg/postgres"
	"github.com/go-redis/redis/v8"
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

	sessionStore, err := session.NewRedisSessionStore(redis.Options{
		Addr:     conf.Redis.Host,
		Password: conf.Redis.Password,
		DB:       conf.Redis.Db,
	})
	if err != nil {
		return err
	}
	defer sessionStore.Stop()

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

	// HTTP/DB clients, repositories
	httpClient := http.Client{}
	db, err := postgres.New(conf.Db)
	if err != nil {
		return err
	}
	sqlDb, _ := db.DB()
	defer sqlDb.Close()

	albumRepo := repository.NewAlbumRepo(
		localDatasource.New(db),
		remoteDatasource.NewDiscogsDatasource(&httpClient),
		remoteDatasource.NewSpotifyDatasource(&httpClient, conf.Spotify),
	)

	// Scenes, controllers and usecases
	mainController := controller.NewMainController(stage)
	controller.NewMainRouter(mainScene, mainController)

	postUsecase := usecase.NewPostUsecase(albumRepo)
	postController := controller.NewChannelPostController(postUsecase, stage)
	controller.NewPostRouter(postScene, postController)

	bot.Start()

	return nil
}
