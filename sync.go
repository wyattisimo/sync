package sync

import (
	"sync"
)

// This Mutex is an extension of Go's standard sync.Mutex.
// It can be locked and unlocked like a normal mutex. Additionally,
// its status can be checked without attempting a lock or unlock,
// and it can be locked only if it is not currently locked.
type Mutex struct {
	mutex  sync.Mutex
	locked bool
}

// Lock locks m. If the lock is already in use,
// the calling goroutine blocks until the mutex is available.
func (m *Mutex) Lock() {
	for {
		m.mutex.Lock()
		if !m.locked {
			m.locked = true
			m.mutex.Unlock()
			return
		}
		m.mutex.Unlock()
	}
}

// Lock if possible checks the current status of m.
// If m is currently unlocked, it locks m, and returns
// true, otherwise it returns false.
func (m *Mutex) LockIfPossible() bool {
	m.mutex.Lock()
	if m.locked {
		m.mutex.Unlock()
		return false
	} else {
		m.locked = true
		m.mutex.Unlock()
		return true
	}
}

// Locked checks to see if m is currently locked.
func (m *Mutex) Locked() bool {
	m.mutex.Lock()
	ret := m.locked
	m.mutex.Unlock()
	return ret
}

// Unlock unlocks m. It is a run-time error if
// m is not locked on entry to Unlock.
func (m *Mutex) Unlock() {
	m.mutex.Lock()
	if !m.locked {
		panic("sync: unlock of locked mutex")
	}
	m.locked = false
	m.mutex.Unlock()
}
