package session

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (s *Session) ShouldGet(key string, dest interface{}) error {
	s.RLock()
	defer s.RUnlock()
	raw, ok := s.Get(key).(map[string]interface{})
	if !ok {
		return errors.New("can't make type assertion")
	}

	jsonTmp, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonTmp))

	return json.Unmarshal(jsonTmp, dest)
}

func (s *Session) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()

	if s.Store == nil {
		s.Store = make(map[string]interface{})
	}

	//val, err := json.Marshal(value)
	//if err != nil {
	//	panic(err)
	//}

	s.Store[key] = value
}

func (s *Session) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.Store, key)
}

func (s *Session) Clear() {
	s.Lock()
	defer s.Unlock()
	for key, _ := range s.Store {
		delete(s.Store, key)
	}
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