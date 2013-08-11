package sync

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// go test -test.bench ".*"
// has been buggy, so (just in case)
// manually make sure that other
// tests aren't running along side
// of the benchmarking function.
var benchMtx sync.Mutex

func TestLock(t *testing.T) {
	benchMtx.Lock()
	m := new(Mutex)
	m.Lock()
	if !m.locked {
		t.Error("Lock() did not lock mutex")
	}
	benchMtx.Unlock()
}

func TestUnlock(t *testing.T) {
	benchMtx.Lock()
	m := new(Mutex)
	m.locked = true
	m.Unlock()
	if m.locked {
		t.Error("Unlock() did not unlock mutex")
	}
	benchMtx.Unlock()
}

func TestLocked(t *testing.T) {
	benchMtx.Lock()
	m := new(Mutex)
	if m.Locked() {
		t.Error("Locked() reported that unlocked mutex was locked")
	}
	m.Lock()
	if !m.Locked() {
		t.Error("Locked() reported that locked mutex was unlocked")
	}
	benchMtx.Unlock()
}

func TestLockIfPossible(t *testing.T) {
	benchMtx.Lock()
	m := new(Mutex)
	if !m.LockIfPossible() {
		t.Error("LockIfPossible() failed to lock unlocked mutex")
	}
	if m.LockIfPossible() {
		t.Error("LockIfPossible() succeeded to lock locked mutex")
	}
	benchMtx.Unlock()
}

func BenchmarkMutex(b *testing.B) {
	benchMtx.Lock()
	t0 := time.Now()
	testMutex()
	t1 := time.Now()
	fmt.Printf("Took %v\n", t1.Sub(t0))
	benchMtx.Unlock()
}

func BenchmarkSyncMutex(b *testing.B) {
	benchMtx.Lock()
	t0 := time.Now()
	testSyncMutex()
	t1 := time.Now()
	fmt.Printf("Took %v\n", t1.Sub(t0))
	benchMtx.Unlock()
}

func testSyncMutex() {
	m := new(Mutex)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(m *Mutex) {
			for i := 0; i < 10000; i++ {
				m.Lock()
				m.Unlock()
			}
			wg.Done()
		}(m)
	}
	wg.Wait()
}

func testMutex() {
	m := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		for i := 0; i < 10000; i++ {
			m.Lock()
			m.Unlock()
		}
		wg.Done()
	}
	wg.Wait()
}
