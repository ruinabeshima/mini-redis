/* Defines a key-value store for the Redis server */

package store

import "sync"

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

/*
Acts as a constructor for the Store struct
Initialises the map (an uninitialised one is nil, and cannot be used)
make() stores the map in memory
Returns a ready to use pointer *Store
*/
func New() *Store {
	return &Store{data: make(map[string]string)}
}

func (s *Store) Set(key, value string) {
	// Exclusive lock: Only one goroutine can hold a write lock
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	// Multiple readers can read concurrently
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.data[key]
	return v, ok
}

func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[key]
	delete(s.data, key) // delete on a missing key is a safe no-op
	return exists
}
