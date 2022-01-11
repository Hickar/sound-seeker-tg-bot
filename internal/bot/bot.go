package bot

import (
	"time"

	"github.com/Hickar/sound-seeker-bot/internal/config"
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

	bot.Start()

	return nil
}
