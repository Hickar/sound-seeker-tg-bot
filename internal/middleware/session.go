package middleware

import (
	"sync"

	"gopkg.in/tucnak/telebot.v3"
)

const SessionKey = "session"

type SessionStore interface {
	Get(int64) (*Session, error)
	Set(int64, *Session)
	Delete(int64)
}

type Session struct {
	sync.RWMutex
	store   map[string]interface{}
}

func NewSession() *Session {
	return &Session{
		RWMutex: sync.RWMutex{},
		store: make(map[string]interface{}),
	}
}

func (s *Session) Get(key string) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()
	value, ok := s.store[key]

	return value, ok
}

func (s *Session) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

	if s.store == nil {
		s.store = make(map[string]interface{})
	}

	s.store[key] = value
}

func (s *Session) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.store, key)
}

func SessionMiddleware(store SessionStore) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			sessionId := ctx.Message().Sender.ID

			session, err := store.Get(sessionId)
			if err != nil || session == nil {
				session = NewSession()
			}

			ctx.Set(SessionKey, session)
			if err := next(ctx); err != nil {
				return err
			}
			store.Set(sessionId, session)

			return nil
		}
	}
}