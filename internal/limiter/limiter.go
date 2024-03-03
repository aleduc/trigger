package limiter

import (
	"sync"

	"test_trigger/internal/realtime"
)

// SlidingWindow is responsible for rate limit sliding window algo.
type SlidingWindow struct {
	RealTime    realtime.Time
	size        int64
	limit       uint64
	windowStart int64
	counter     uint64   // not mandatory field.
	entries     []uint64 // can be replaced with list.
	mu          *sync.Mutex
}

func NewSlidingWindow(size uint64, limit uint64, t realtime.Time) *SlidingWindow {
	return &SlidingWindow{
		RealTime:    t,
		size:        int64(size),
		limit:       limit,
		windowStart: t.Now().Unix(),
		entries:     make([]uint64, size),
		mu:          &sync.Mutex{},
	}
}

// Allow returns false if limit exceeded.
func (s *SlidingWindow) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	currentTime := s.RealTime.Now().Unix()
	s.refreshCurrentState(currentTime)
	if s.counter+1 > s.limit {
		return false
	}
	s.entries[currentTime%s.size] += 1
	s.counter += 1
	return true
}

func (s *SlidingWindow) refreshCurrentState(currentTime int64) {
	toDelete := currentTime - s.windowStart - s.size + 1
	if toDelete <= 0 {
		return
	}
	if toDelete >= s.size {
		// There is reinit instead of ring buffer, just for simplification.
		s.entries = make([]uint64, s.size)
		s.counter = 0
		s.windowStart = currentTime
		return
	}
	localHead := s.windowStart % s.size
	s.windowStart = s.windowStart + toDelete
	for i := int64(0); i < toDelete; i++ {
		if localHead+i > s.size-1 {
			localHead = -i
		}
		s.counter -= s.entries[localHead+i]
		s.entries[localHead+i] = 0
	}
}
