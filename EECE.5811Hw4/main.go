package main

import (
    "fmt"
    "runtime"
    "sync"
    "sync/atomic"
    "time"
    "unsafe"
)



/******************************************************************************
 *                           LOCK-BASED QUEUE                                 *
 ******************************************************************************/

// A simple node for the lock-based queue
type node struct {
    value interface{}
    next  *node
}

// LockQueue uses a single mutex to protect head and tail
type LockQueue struct {
    head *node
    tail *node
    mu   sync.Mutex
}

// NewLockQueue initializes an empty queue with a sentinel node
func NewLockQueue() *LockQueue {
    sentinel := &node{} // sentinel node, no real value
    return &LockQueue{
        head: sentinel,
        tail: sentinel,
    }
}

// Enqueue adds a new element at the tail of the queue
func (q *LockQueue) Enqueue(val interface{}) {
    newNode := &node{value: val}

    q.mu.Lock()
    defer q.mu.Unlock()

    // Link the new node, then update tail
    q.tail.next = newNode
    q.tail = newNode
}

// Dequeue removes and returns the head of the queue
// Returns (nil, false) if empty
func (q *LockQueue) Dequeue() (interface{}, bool) {
    q.mu.Lock()
    defer q.mu.Unlock()

    first := q.head.next
    if first == nil {
        return nil, false
    }
    q.head = first
    return first.value, true
}

/******************************************************************************
 *                        MICHAEL & SCOTT LOCK-FREE QUEUE                     *
 ******************************************************************************/

// lfNode is a node used by the lock-free queue, updated with atomic operations
type lfNode struct {
    value interface{}
    next  unsafe.Pointer // *lfNode
}

// LFQueue maintains head and tail pointers atomically
type LFQueue struct {
    head unsafe.Pointer // *lfNode
    tail unsafe.Pointer // *lfNode
}

// NewLFQueue initializes a lock-free queue with a sentinel node
func NewLFQueue() *LFQueue {
    sentinel := &lfNode{}
    q := &LFQueue{}
    // Atomically store pointers to the sentinel node
    atomic.StorePointer(&q.head, unsafe.Pointer(sentinel))
    atomic.StorePointer(&q.tail, unsafe.Pointer(sentinel))
    return q
}

// Enqueue implements Michael & Scott's non-blocking enqueue
func (q *LFQueue) Enqueue(val interface{}) {
    newNode := &lfNode{value: val}
    for {
        tailPtr := atomic.LoadPointer(&q.tail)
        tail := (*lfNode)(tailPtr)
        nextPtr := atomic.LoadPointer(&tail.next)

        // Check if tail is still the last node
        if tailPtr == atomic.LoadPointer(&q.tail) {
            // If next is nil, we can try linking our newNode there
            if nextPtr == nil {
                if atomic.CompareAndSwapPointer(&tail.next, nextPtr, unsafe.Pointer(newNode)) {
                    // Enqueue is done; try to advance tail pointer
                    atomic.CompareAndSwapPointer(&q.tail, tailPtr, unsafe.Pointer(newNode))
                    return
                }
            } else {
                // Tail was not pointing to last node, help move tail forward
                atomic.CompareAndSwapPointer(&q.tail, tailPtr, nextPtr)
            }
        }
    }
}

// Dequeue implements Michael & Scott's non-blocking dequeue
func (q *LFQueue) Dequeue() (interface{}, bool) {
    for {
        headPtr := atomic.LoadPointer(&q.head)
        tailPtr := atomic.LoadPointer(&q.tail)
        head := (*lfNode)(headPtr)
        tail := (*lfNode)(tailPtr)
        firstPtr := atomic.LoadPointer(&head.next)
        first := (*lfNode)(firstPtr)

        // Check if queue is empty
        if head == tail {
            // If next pointer is nil, queue is truly empty
            if first == nil {
                return nil, false
            }
            // Otherwise tail is falling behind, try to advance it
            atomic.CompareAndSwapPointer(&q.tail, tailPtr, firstPtr)
        } else {
            // Read the value before CAS
            val := first.value
            // Try to swing head to the next node
            if atomic.CompareAndSwapPointer(&q.head, headPtr, firstPtr) {
                return val, true
            }
        }
    }
}

/******************************************************************************
 *                               BENCHMARKING                                 *
 ******************************************************************************/

// benchQueue is a simple interface for Enqueue/Dequeue
type benchQueue interface {
    Enqueue(interface{})
    Dequeue() (interface{}, bool)
}

// benchmarkScenario runs a single scenario with given #producers, #consumers, #ops
func benchmarkScenario(q benchQueue, producers, consumers, ops int) time.Duration {
    wg := &sync.WaitGroup{}
    start := time.Now()

    // Start producers
    for p := 0; p < producers; p++ {
        wg.Add(1)
        go func(pid int) {
            defer wg.Done()
            for i := 0; i < ops; i++ {
                // Example: enqueue a simple string
                q.Enqueue(fmt.Sprintf("P%d-Val%d", pid, i))
            }
        }(p)
    }

    // Start consumers
    for c := 0; c < consumers; c++ {
        wg.Add(1)
        go func(cid int) {
            defer wg.Done()
            for i := 0; i < ops; i++ {
                q.Dequeue()
            }
        }(c)
    }

    wg.Wait()
    return time.Since(start)
}

/******************************************************************************
 *                                   MAIN                                      *
 ******************************************************************************/

func main() {
    // Let Go use all available CPU cores
    runtime.GOMAXPROCS(runtime.NumCPU())

    // We'll run multiple scenarios:
    // - concurrencyLevels: (producers, consumers) pairs
    // - operations: how many ops each producer/consumer will do
    concurrencyLevels := [][]int{
        {1, 1},
        {2, 2},
        {4, 4},
        {8, 8},
    }
    operationCounts := []int{10000, 50000}

    fmt.Println("Comparison of Lock-based vs Lock-free Queues")
    fmt.Println("We will benchmark various (producers, consumers, operations) workloads.\n")

    // Print table header
    fmt.Printf("%-10s %-10s %-10s %-20s %-20s\n",
        "Prod", "Cons", "Ops", "Lock-Based (ms)", "Lock-Free (ms)")

    for _, cc := range concurrencyLevels {
        producers := cc[0]
        consumers := cc[1]

        for _, ops := range operationCounts {
            // 1. Test Lock-based queue
            lockQ := NewLockQueue()
            elapsedLock := benchmarkScenario(lockQ, producers, consumers, ops)

            // 2. Test Lock-free (Michael & Scott) queue
            lfQ := NewLFQueue()
            elapsedLF := benchmarkScenario(lfQ, producers, consumers, ops)

            // Print results in milliseconds for easier reading
            fmt.Printf("%-10d %-10d %-10d %-20.3f %-20.3f\n",
                producers, consumers, ops,
                float64(elapsedLock.Milliseconds()),
                float64(elapsedLF.Milliseconds()),
            )
        }
    }

    fmt.Println("\nDONE.")
    fmt.Println("Analyze the above timings to see how lock contention vs lock-free overhead scales.")
    fmt.Println("Typically, lock-free queues may outperform locked queues at high concurrency,")
    fmt.Println("but at low concurrency the simpler locked queue can be faster.")
}

