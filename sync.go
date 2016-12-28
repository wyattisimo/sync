// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"fmt"
	"sync"
	"time"
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

// LockIfPossible checks the current status of m.
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

// LockOrPanic locks m. If the lock is already in use, it panics.
func (m *Mutex) LockOrPanic() {
	if !m.LockIfPossible() {
		panic("sync: unable to acquire mutex lock")
	}
}

// LockOrPanicAfter locks m. If the lock is already in use,
// it blocks until the mutex is available or timeout expires.
// If timeout expires before the mutex is available, it panics.
func (m *Mutex) LockOrPanicAfter(timeout time.Duration) {
	haveLock := make(chan bool, 1)
	go func() {
		for {
			if m.LockIfPossible() {
				haveLock <- true
				return
			}
		}
	}()
	select {
	case <-haveLock:
		close(haveLock)
		return
	case <-time.After(timeout):
		panic(fmt.Sprintf("sync: unable to acquire mutex lock after %v", timeout))
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
		panic("sync: unlock of unlocked mutex")
	}
	m.locked = false
	m.mutex.Unlock()
}
