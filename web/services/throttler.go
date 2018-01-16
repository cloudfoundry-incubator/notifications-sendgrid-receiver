package services

import (
	"sync"
)

type Throttler struct {
	maxRequests    int64
	mu             sync.Mutex
	activeRequests int64
}

func NewThrottler(maxRequests int) *Throttler {
	return &Throttler{
		maxRequests: int64(maxRequests),
	}
}

func (t *Throttler) Throttle() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.activeRequests >= t.maxRequests {
		return true
	} else {
		t.activeRequests++
		return false
	}
}

func (t *Throttler) Finish() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.activeRequests--
}
