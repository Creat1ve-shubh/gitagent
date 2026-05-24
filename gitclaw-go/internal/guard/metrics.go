package guard

import (
	"sync"
	"sync/atomic"
)

// Metrics tracks guard pipeline statistics using atomic counters
// for lock-free concurrent access.
type Metrics struct {
	allowed  atomic.Int64
	rejected atomic.Int64

	// Per-guard rejection counts
	mu                sync.RWMutex
	rejectionsByGuard map[string]*atomic.Int64
}

// MetricsSnapshot is a point-in-time copy of guard metrics.
type MetricsSnapshot struct {
	TotalAllowed      int64            `json:"total_allowed"`
	TotalRejected     int64            `json:"total_rejected"`
	RejectionsByGuard map[string]int64 `json:"rejections_by_guard"`
}

// NewMetrics creates a new Metrics tracker.
func NewMetrics() *Metrics {
	return &Metrics{
		rejectionsByGuard: make(map[string]*atomic.Int64),
	}
}

// RecordAllow increments the allowed counter.
func (m *Metrics) RecordAllow() {
	m.allowed.Add(1)
}

// RecordRejection increments the rejected counter and the per-guard counter.
func (m *Metrics) RecordRejection(guardName string) {
	m.rejected.Add(1)

	m.mu.RLock()
	counter, ok := m.rejectionsByGuard[guardName]
	m.mu.RUnlock()

	if !ok {
		m.mu.Lock()
		// Double-check after acquiring write lock
		counter, ok = m.rejectionsByGuard[guardName]
		if !ok {
			counter = &atomic.Int64{}
			m.rejectionsByGuard[guardName] = counter
		}
		m.mu.Unlock()
	}

	counter.Add(1)
}

// Snapshot returns a point-in-time copy of all metrics.
func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	byGuard := make(map[string]int64, len(m.rejectionsByGuard))
	for name, counter := range m.rejectionsByGuard {
		byGuard[name] = counter.Load()
	}

	return MetricsSnapshot{
		TotalAllowed:      m.allowed.Load(),
		TotalRejected:     m.rejected.Load(),
		RejectionsByGuard: byGuard,
	}
}
