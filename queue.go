package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// -----------------------------------------------------------------------------
//  1) Common Interface
// -----------------------------------------------------------------------------

type ConcurrentQueue interface {
    Enqueue(value int)
    Dequeue() (int, bool)
}

// -----------------------------------------------------------------------------
//  2) Baseline Queue (Coarse-Grained Locking)
// -----------------------------------------------------------------------------

type BaselineQueue struct {
    head  *node
    tail  *node
    lock  sync.Mutex
}

type node struct {
    value int
    next  *node
}

// NewBaselineQueue creates a queue with a dummy node and a single lock.
func NewBaselineQueue() *BaselineQueue {
    dummy := &node{}
    return &BaselineQueue{
        head: dummy,
        tail: dummy,
    }
}

func (q *BaselineQueue) Enqueue(value int) {
    q.lock.Lock()
    defer q.lock.Unlock()

    newNode := &node{value: value}
    q.tail.next = newNode
    q.tail = newNode
}

func (q *BaselineQueue) Dequeue() (int, bool) {
    q.lock.Lock()
    defer q.lock.Unlock()

    if q.head.next == nil {
        return 0, false // empty
    }
    front := q.head.next
    q.head.next = front.next
    if q.head.next == nil {
        q.tail = q.head
    }
    return front.value, true
}

// -----------------------------------------------------------------------------
//  3) Hand-Over-Hand Queue (Lock Coupling)
// -----------------------------------------------------------------------------

type hohNode struct {
    value int
    next  *hohNode
    nlock sync.Mutex
}

type HOHQueue struct {
    head *hohNode
    tail *hohNode
}

// NewHOHQueue initializes the queue with a dummy node.
func NewHOHQueue() *HOHQueue {
    dummy := &hohNode{}
    return &HOHQueue{
        head: dummy,
        tail: dummy,
    }
}

// Enqueue uses a small loop to safely lock the current tail node, verify
// it is still the tail, then append the new node.
func (q *HOHQueue) Enqueue(value int) {
    newNode := &hohNode{value: value}
    
    for {
        // Capture current tail pointer
        oldTail := q.tail
        
        // Lock oldTail
        oldTail.nlock.Lock()
        
        // Check if q.tail is still oldTail
        if q.tail != oldTail {
            // Another goroutine changed the tail before we locked,
            // so unlock and retry.
            oldTail.nlock.Unlock()
            continue
        }
        
        // Now oldTail is definitely the real tail; link in the new node
        oldTail.next = newNode
        
        // Update the queue's tail pointer to newNode
        q.tail = newNode
        
        // Unlock the oldTail
        oldTail.nlock.Unlock()
        break
    }
}

// Dequeue locks the dummy head, then the actual front node, removes it,
// and unlocks in reverse order.
func (q *HOHQueue) Dequeue() (int, bool) {
    // Lock the dummy head
    q.head.nlock.Lock()
    defer q.head.nlock.Unlock()

    if q.head.next == nil {
        return 0, false // empty
    }

    // Lock the next node so we can remove it safely
    nextNode := q.head.next
    nextNode.nlock.Lock()
    defer nextNode.nlock.Unlock()

    val := nextNode.value
    q.head.next = nextNode.next

    // If we removed the last node, reset tail
    if q.head.next == nil {
        q.tail = q.head
    }

    return val, true
}

// -----------------------------------------------------------------------------
//  4) Benchmark Harness
// -----------------------------------------------------------------------------

type workload struct {
    numEnqueue int
    numDequeue int
}

// runBenchmark spawns numGoroutines, each of which performs enqueue
// operations followed by dequeue operations. It returns total time taken.
func runBenchmark(q ConcurrentQueue, numGoroutines int, wl workload) time.Duration {
    var wg sync.WaitGroup
    wg.Add(numGoroutines)

    start := time.Now()

    for g := 0; g < numGoroutines; g++ {
        go func() {
            defer wg.Done()
            // Perform enqueues
            for i := 0; i < wl.numEnqueue; i++ {
                q.Enqueue(rand.Intn(1_000_000))
            }
            // Perform dequeues
            for i := 0; i < wl.numDequeue; i++ {
                _, _ = q.Dequeue()
            }
        }()
    }

    wg.Wait()
    return time.Since(start)
}

func main() {
    rand.Seed(time.Now().UnixNano())

    // Some example workloads
    workloads := []struct {
        name string
        wl   workload
    }{
        {"Write-Heavy", workload{numEnqueue: 100000, numDequeue: 10000}},
        {"Read-Heavy",  workload{numEnqueue: 10000,  numDequeue: 100000}},
        {"Balanced",    workload{numEnqueue: 50000,  numDequeue: 50000}},
    }
    numGoroutinesList := []int{1, 2, 4, 8, 16}

    fmt.Println("Baseline Queue vs Hand-Over-Hand Queue Benchmark\n")

    // Compare BaselineQueue and HOHQueue in a table of (Workload x #Goroutines)
    for _, w := range workloads {
        fmt.Printf("=== Workload: %s (Enq=%d, Deq=%d) ===\n",
            w.name, w.wl.numEnqueue, w.wl.numDequeue)
        for _, g := range numGoroutinesList {
            // Baseline
            bq := NewBaselineQueue()
            durationB := runBenchmark(bq, g, w.wl)

            // Hand-over-hand
            hq := NewHOHQueue()
            durationH := runBenchmark(hq, g, w.wl)

            fmt.Printf("Goroutines=%2d | Baseline=%v | HOH=%v\n",
                g, durationB, durationH)
        }
        fmt.Println()
    }
}
