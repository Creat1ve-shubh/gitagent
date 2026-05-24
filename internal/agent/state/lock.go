package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Lock struct {
	Path string
	TTL  time.Duration
}

type lockInfo struct {
	PID       int       `json:"pid"`
	Timestamp time.Time `json:"timestamp"`
}

func NewLock(root string, ttl time.Duration) *Lock {
	return &Lock{Path: filepath.Join(root, ".gitagent", "lock.json"), TTL: ttl}
}

func (l *Lock) Acquire() error {
	if err := os.MkdirAll(filepath.Dir(l.Path), 0o755); err != nil {
		return err
	}

	if ok, err := l.tryAcquire(); ok || err != nil {
		return err
	}

	return errors.New("lock already held")
}

func (l *Lock) tryAcquire() (bool, error) {
	if _, err := os.Stat(l.Path); err == nil {
		stale, err := l.isStale()
		if err != nil {
			return false, err
		}
		if !stale {
			return false, nil
		}
		_ = os.Remove(l.Path)
	}

	info := lockInfo{PID: os.Getpid(), Timestamp: time.Now().UTC()}
	data, _ := json.Marshal(info)
	return true, os.WriteFile(l.Path, data, 0o600)
}

func (l *Lock) Release() error {
	return os.Remove(l.Path)
}

func (l *Lock) isStale() (bool, error) {
	data, err := os.ReadFile(l.Path)
	if err != nil {
		return false, err
	}
	var info lockInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return true, nil
	}
	return time.Since(info.Timestamp) > l.TTL, nil
}
