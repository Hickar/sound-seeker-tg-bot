package bot

import (
	"errors"
	"time"

	"github.com/Hickar/sound-seeker-bot/internal/config"
	"github.com/Hickar/sound-seeker-bot/internal/controller"
	"github.com/Hickar/sound-seeker-bot/internal/middleware"
	"gopkg.in/tucnak/telebot.v3"
)

type App struct {
	bot telebot.Bot
}

type KVStore struct {
	store map[int64]*middleware.Session
}

func NewSessionStore() *KVStore {
	return &KVStore{store: make(map[int64]*middleware.Session)}
}

func (s *KVStore) Get(key int64) (*middleware.Session, error) {
	session, ok := s.store[key]
	if !ok || session == nil {
		return nil, errors.New("no such value is store")
	}

	return session, nil
}

func (s *KVStore) Set(key int64, session *middleware.Session) {
	s.store[key] = session
}

func (s *KVStore) Delete(key int64) {
	delete(s.store, key)
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

	bot.Use(middleware.SessionMiddleware(NewSessionStore()))

	if err := controller.NewMainRouter(bot); err != nil {
		return err
	}

	bot.Start()

	return nil
}
