package main

import (
    "sync/atomic"
    "runtime"
)

// TicketLock implements a simple FIFO ticket lock.
type TicketLock struct {
    ticket int32
    turn   int32
}

// NewTicketLock returns a pointer to a new ticket lock.
func NewTicketLock() *TicketLock {
    return &TicketLock{
        ticket: 0,
        turn:   0,
    }
}

// Lock acquires the ticket lock, spinning until the caller's ticket is up.
func (tl *TicketLock) Lock() {
    // Atomically fetch and increment the ticket counter
    myTurn := atomic.AddInt32(&tl.ticket, 1) - 1
    // Spin while waiting for your turn
    for atomic.LoadInt32(&tl.turn) != myTurn {
        // Optional: yield CPU to other goroutines
        runtime.Gosched()
    }
}

// Unlock releases the ticket lock, allowing the next waiting goroutine in line to proceed.
func (tl *TicketLock) Unlock() {
    atomic.AddInt32(&tl.turn, 1)
}
