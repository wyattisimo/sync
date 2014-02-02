// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"sync"
	"testing"
)

func TestLock(t *testing.T) {
	m := new(Mutex)
	m.Lock()
	if !m.locked {
		t.Error("Lock() did not lock mutex")
	}
}

func TestUnlock(t *testing.T) {
	m := new(Mutex)
	m.locked = true
	m.Unlock()
	if m.locked {
		t.Error("Unlock() did not unlock mutex")
	}
}

func TestLocked(t *testing.T) {
	m := new(Mutex)
	if m.Locked() {
		t.Error("Locked() reported that unlocked mutex was locked")
	}
	m.Lock()
	if !m.Locked() {
		t.Error("Locked() reported that locked mutex was unlocked")
	}
}

func TestLockIfPossible(t *testing.T) {
	m := new(Mutex)
	if !m.LockIfPossible() {
		t.Error("LockIfPossible() failed to lock unlocked mutex")
	}
	if m.LockIfPossible() {
		t.Error("LockIfPossible() succeeded to lock locked mutex")
	}
}

func BenchmarkMutex(b *testing.B) {
	m := new(sync.Mutex)
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			m.Lock()
			m.Unlock()
		}
	}()
	for i := 0; i < b.N; i++ {
		m.Lock()
		m.Unlock()
	}
}

func BenchmarkSyncMutex(b *testing.B) {
	m := new(Mutex)
	b.ResetTimer()
	go func() {
		for i := 0; i < b.N; i++ {
			m.Lock()
			m.Unlock()
		}
	}()
	for i := 0; i < b.N; i++ {
		m.Lock()
		m.Unlock()
	}
}
