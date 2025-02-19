package main

import (
    "sync/atomic"
    "runtime"
)

// CASLock implements a simple compare-and-swap spin lock.
type CASLock struct {
    locked int32 // 0 => unlocked, 1 => locked
}

// NewCASLock constructs a new CAS spin lock.
func NewCASLock() *CASLock {
    return &CASLock{locked: 0}
}

// Lock acquires the lock by spinning until locked flips from 0 to 1.
func (l *CASLock) Lock() {
    for !atomic.CompareAndSwapInt32(&l.locked, 0, 1) {
        // Optional: yield to other goroutines
        runtime.Gosched()
    }
}

// Unlock releases the lock by storing 0 (unlocked) into locked.
func (l *CASLock) Unlock() {
    atomic.StoreInt32(&l.locked, 0)
}
