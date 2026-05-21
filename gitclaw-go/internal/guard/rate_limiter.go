package guard

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// slidingWindow tracks call counts within a time window using atomic counters.
type slidingWindow struct {
	count     atomic.Int64
	windowEnd atomic.Int64 // unix nano of current window end
}

// RateLimiter enforces a maximum number of tool calls per sliding window
// per session. Uses atomic counters — no locks on the hot path.
type RateLimiter struct {
	windowSize   time.Duration
	maxPerWindow int64
	windows      sync.Map // sessionID → *slidingWindow
}

// RateLimiterConfig configures the rate limiter.
type RateLimiterConfig struct {
	// WindowSize is the sliding window duration (e.g., 60s).
	WindowSize time.Duration

	// MaxPerWindow is the maximum tool calls allowed within the window.
	MaxPerWindow int
}

// NewRateLimiter creates a new RateLimiter with the given config.
func NewRateLimiter(cfg RateLimiterConfig) *RateLimiter {
	if cfg.WindowSize == 0 {
		cfg.WindowSize = 60 * time.Second
	}
	if cfg.MaxPerWindow == 0 {
		cfg.MaxPerWindow = 100
	}
	return &RateLimiter{
		windowSize:   cfg.WindowSize,
		maxPerWindow: int64(cfg.MaxPerWindow),
	}
}

func (r *RateLimiter) Name() string { return "rate_limiter" }

func (r *RateLimiter) Check(req ToolRequest) ToolResponse {
	w := r.getOrCreateWindow(req.SessionID)
	now := time.Now().UnixNano()

	// Check if we've moved past the current window
	windowEnd := w.windowEnd.Load()
	if now > windowEnd {
		// Rotate window: reset counter, set new window end
		newEnd := now + r.windowSize.Nanoseconds()
		if w.windowEnd.CompareAndSwap(windowEnd, newEnd) {
			w.count.Store(1) // this is our first call in the new window
			return ToolResponse{Allowed: true}
		}
		// CAS failed — another goroutine rotated first, fall through to increment
	}

	count := w.count.Add(1)
	if count > r.maxPerWindow {
		return ToolResponse{
			Allowed: false,
			Reason: fmt.Sprintf(
				"rate limit exceeded: %d/%d calls in %v window",
				count, r.maxPerWindow, r.windowSize,
			),
		}
	}

	return ToolResponse{Allowed: true}
}

func (r *RateLimiter) getOrCreateWindow(sessionID string) *slidingWindow {
	if v, ok := r.windows.Load(sessionID); ok {
		return v.(*slidingWindow)
	}

	w := &slidingWindow{}
	w.windowEnd.Store(time.Now().Add(r.windowSize).UnixNano())
	actual, _ := r.windows.LoadOrStore(sessionID, w)
	return actual.(*slidingWindow)
}
