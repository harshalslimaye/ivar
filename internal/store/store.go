package store

import "sync"

type Store struct {
	data sync.Map
}

var (
	instance *Store
	once     sync.Once
)

func (s *Store) Get(key string) interface{} {
	value, exists := s.data.Load(key)
	if !exists {
		return nil
	}

	return value
}

func (s *Store) Set(key string, value interface{}) {
	s.data.Store(key, value)
}

func GetStore() *Store {
	once.Do(func() {
		instance = &Store{}
	})
	return instance
}
