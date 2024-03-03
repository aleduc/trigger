package call

import (
	"context"
	"sync"
)

// Storage stores calls for processing.
// Implementation can be with real database, buffered channel, ring buffer like in limiter, etc.
// I've used slices for queue(not channel), since we always should respond fast regardless workers loading.
// []Meta can be []*Meta. Can save some memory allocations, and will be possible to free memory earlier in the following way s[0]=nil, but it also has cons.
// Ring buffer implementation doesn't fit within time frame.
// Context in input, error in output are for future implementation with database.
type Storage struct {
	toProcess []Meta
	statuses  map[ID]int
	mu        *sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{toProcess: make([]Meta, 0), statuses: make(map[ID]int), mu: &sync.Mutex{}}
}

// AddToQueueBack adds meta to the end of the queue.
func (s *Storage) AddToQueueBack(_ context.Context, meta Meta) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.toProcess = append(s.toProcess, meta)
	return nil
}

func (s *Storage) AddToQueueFront(_ context.Context, meta Meta) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.toProcess = append(s.toProcess, Meta{})
	copy(s.toProcess[1:], s.toProcess)
	s.toProcess[0] = meta
	return nil
}

func (s *Storage) Next(_ context.Context) (Meta, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.toProcess) == 0 {
		return Meta{}, false, nil
	}
	res := s.toProcess[0]
	s.toProcess = s.toProcess[1:]
	return res, true, nil
}

func (s *Storage) SaveStatus(_ context.Context, status int, meta Meta) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statuses[meta.ID] = status
	return nil
}

func (s *Storage) QueueLength(_ context.Context) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.toProcess), nil
}
