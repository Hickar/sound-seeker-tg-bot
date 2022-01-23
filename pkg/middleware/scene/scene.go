package scene

import (
	"fmt"

	"github.com/Hickar/sound-seeker-bot/internal/config"
	"gopkg.in/tucnak/telebot.v3"
)

type Scene struct {
	Name         string
	Stage        *Stage
	startHandler telebot.HandlerFunc
	*telebot.Bot
}

func NewScene(conf config.BotConfig, name string, middleware ...telebot.MiddlewareFunc) (*Scene, error) {
	botConf := telebot.Settings{
		Token:       conf.Token,
		Poller:      nil,
		Synchronous: !conf.Concurrent,
		Verbose:     conf.Debug,
	}

	bot, err := telebot.NewBot(botConf)
	if err != nil {
		return nil, fmt.Errorf("error during \"%s\" scene creation: %s", name, err.Error())
	}

	if len(middleware) > 0 {
		bot.Use(middleware...)
	}

	return &Scene{Name: name, Bot: bot}, nil
}

func newDefaultScene(stage *Stage) *Scene {
	return &Scene{Name: defaultScene, Stage: stage}
}

func (sc *Scene) HandleStart(handler telebot.HandlerFunc) {
	sc.startHandler = handler
}
