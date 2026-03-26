package container

import "sync"

type Store struct {
	mu         sync.RWMutex
	containers map[string]*Container
}

func NewStore() *Store {
	return &Store{
		containers: make(map[string]*Container),
	}
}

func (s *Store) Add(c *Container) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.containers[c.ID] = c
}

func (s *Store) Get(id string) (*Container, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.containers[id]
	return c, ok
}

func (s *Store) List() []*Container {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []*Container
	for _, c := range s.containers {
		out = append(out, c)
	}
	return out
}

func (s *Store) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.containers, id)
}
