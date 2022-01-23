package session

type SessionStore interface {
	Get(int64) *Session
	Set(int64, *Session)
	Delete(int64)
}

type MemorySessionStore struct {
	store map[int64]*Session
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{store: make(map[int64]*Session)}
}

func (s *MemorySessionStore) Get(key int64) *Session {
	session, ok := s.store[key]
	if !ok || session == nil {
		return nil
	}

	return session
}

func (s *MemorySessionStore) Set(key int64, session *Session) {
	s.store[key] = session
}

func (s *MemorySessionStore) Delete(key int64) {
	delete(s.store, key)
}
