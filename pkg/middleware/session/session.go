package session

import (
	"sync"

	"gopkg.in/tucnak/telebot.v3"
)

const sessionKey = "session"

type Session struct {
	sync.RWMutex
	Store map[string]interface{}
}

func NewSession() *Session {
	return &Session{
		RWMutex: sync.RWMutex{},
		Store:   make(map[string]interface{}),
	}

}

func NewSessionWithStore(store map[string]interface{}) *Session {
	return &Session{
		RWMutex: sync.RWMutex{},
		Store:   store,
	}
}

func (s *Session) Get(key string) interface{} {
	s.RLock()
	defer s.RUnlock()
	value, ok := s.Store[key]
	if !ok {
		return nil
	}

	return value
}

func (s *Session) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

	if s.Store == nil {
		s.Store = make(map[string]interface{})
	}

	s.Store[key] = value
}

func (s *Session) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Store, key)
}

func Middleware(store SessionStore) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			sessionId := ctx.Message().Sender.ID

			session := store.Get(sessionId)
			if session == nil {
				session = NewSession()
			}

			ctx.Set(sessionKey, session)
			if err := next(ctx); err != nil {
				return err
			}
			store.Set(sessionId, session)

			return nil
		}
	}
}

func GetSession(ctx telebot.Context) *Session {
	session, ok := ctx.Get(sessionKey).(*Session)
	if session == nil || !ok {
		return NewSession()
	}

	return session
}
