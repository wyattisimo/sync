// Copyright 2013 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"sync"
	"testing"
	"time"
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

func TestLockOrPanic_ShouldLock(t *testing.T) {
	m := new(Mutex)
	m.LockOrPanic()
	if !m.locked {
		t.Error("LockOrPanic() did not lock mutex")
	}
}

func TestLockOrPanic_ShouldPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("LockOrPanic() did not panic after unable to acquire mutex lock")
		}
	}()
	m := new(Mutex)
	m.locked = true
	m.LockOrPanic()
}

func TestLockOrPanicAfter_ShouldEventuallyLock(t *testing.T) {
	m := new(Mutex)
	m.locked = true
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("LockOrPanicAfter() did not acquire mutex lock before time expired")
			}
			wg.Done()
		}()
		m.LockOrPanicAfter(10 * time.Millisecond)
	}()
	time.Sleep(5 * time.Millisecond)
	m.Unlock()
	wg.Wait()
}

func TestLockOrPanicAfter_ShouldEventuallyPanic(t *testing.T) {
	m := new(Mutex)
	m.locked = true
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("LockOrPanicAfter() did not panic after time expired")
			}
			wg.Done()
		}()
		m.LockOrPanicAfter(1 * time.Millisecond)
	}()
	time.Sleep(5 * time.Millisecond)
	m.Unlock()
	wg.Wait()
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
