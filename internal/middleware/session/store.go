package session

type SessionStore interface {
	Get(int64) *Session
	Set(int64, *Session)
	Delete(int64)
}

type KVStore struct {
	store map[int64]*Session
}

func NewSessionStore() *KVStore {
	return &KVStore{store: make(map[int64]*Session)}
}

func (s *KVStore) Get(key int64) *Session {
	session, ok := s.store[key]
	if !ok || session == nil {
		return nil
	}

	return session
}

func (s *KVStore) Set(key int64, session *Session) {
	s.store[key] = session
}

func (s *KVStore) Delete(key int64) {
	delete(s.store, key)
}