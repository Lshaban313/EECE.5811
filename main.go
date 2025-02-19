package main

import (
    "flag"
    "fmt"
    "sync/atomic"
    "time"
    "runtime"
)

// benchLock runs multiple goroutines, each repeatedly locking/unlocking,
// and measures how long on average each Lock() call waited.
func benchLock(lock interface {
    Lock()
    Unlock()
}, goroutines, iterations int) time.Duration {

    var totalWait int64  // accumulates total wait time (ns) across all goroutines
    startBarrier := make(chan struct{})
    done := make(chan struct{})
    finishedCount := int32(0)

    // Launch goroutines
    for g := 0; g < goroutines; g++ {
        go func() {
            // Wait for the "start signal"
            <-startBarrier

            localSum := int64(0)
            for i := 0; i < iterations; i++ {
                t0 := time.Now()
                lock.Lock()
                waitTime := time.Since(t0).Nanoseconds()
                
                lock.Unlock()

                localSum += waitTime
            }
            atomic.AddInt64(&totalWait, localSum)

            // Signal we've finished
            if atomic.AddInt32(&finishedCount, 1) == int32(goroutines) {
                close(done)
            }
        }()
    }

    // "Broadcast" start to all goroutines
    close(startBarrier)

    // Wait until everyone is done
    <-done

    //  average wait per lock acquisition
    totalOps := int64(goroutines * iterations)
    avgWaitNanos := totalWait / totalOps
    return time.Duration(avgWaitNanos) * time.Nanosecond
}

func main() {
    // Parse command-line flags
    goroutinesFlag := flag.Int("g", 4, "Number of goroutines")
    iterationsFlag := flag.Int("i", 10000, "Iterations (lock/unlock) per goroutine")
    flag.Parse()

    // Let us confirm how many OS threads Go can use
    runtime.GOMAXPROCS(runtime.NumCPU())

    g := *goroutinesFlag
    iters := *iterationsFlag

    fmt.Printf("Running with %d goroutines, %d iterations each...\n", g, iters)

    // 1. Benchmark the Ticket Lock
    tLock := NewTicketLock()
    avgTicketWait := benchLock(tLock, g, iters)
    fmt.Printf("TicketLock average wait: %v\n", avgTicketWait)

    // 2. Benchmark the CAS Spin Lock
    casLock := NewCASLock()
    avgCASWait := benchLock(casLock, g, iters)
    fmt.Printf("CASLock average wait:    %v\n", avgCASWait)
}
