# EECE.5811: HW4

> **Group Members**: Vina Dang, Layann Shaban  
> **Due Date**: March 5, 2025  

## Overview
```
This repository contains two concurrent queue implementations in Go that satisfy Problem 2’s requirements:

Lock‐Based Queue (mutex‐protected, coarse‐grained)
Michael & Scott Lock‐Free Queue (non‐blocking CAS)
All functionality and a benchmarking harness are located in a single file named main.go, which performs multiple test scenarios to compare performance under various workloads (different numbers of producers, consumers, and operations).
```

## File Descriptions
```
main.go
This file provides:
LockQueue: a queue protected by a single mutex for enqueue and dequeue.
LFQueue (Michael & Scott): a lock-free, non-blocking queue using atomic compare-and-swap for both enqueue and dequeue.
A benchmark that runs multiple concurrency scenarios (e.g., 1×1, 2×2, 4×4, 8×8 producers/consumers), each performing varying numbers of operations (10,000 or 50,000).
Printed timing results, so you can see how each queue scales with higher concurrency.
```

## How to Run
```
1. Clone or download this repository.
2. Open a terminal, navigate to the folder containing main.go.
3. Compile & Run:
4. Build: go build main.go
5. Execute: ./main
Alternatively, you can directly run without building:
go run main.go
```

## Notes
```
You can modify the concurrency levels or operation counts in the main() function to suit different workloads.
The code assumes there is no ABA issue, as per the assignment instructions. In practice, lock-free queues must ensure safe memory reclamation.
Real-world usage might also include extra instrumentation (like a "count" field), but we're ignoring that here.
```

## Design of the Program
```
Lock-Based Queue (LockQueue)

Data:
Two pointers: head and tail
A single sync.Mutex protecting both
Enqueue(value interface{}):
Lock the mutex
Link in a newly allocated node at the tail
Unlock
Dequeue() (interface{}, bool):
Lock the mutex
Remove the first node if present
Unlock
Michael & Scott Lock-Free Queue (LFQueue)

Data:
Two atomically updated pointers: head and tail
A sentinel dummy node
Enqueue(value interface{}):
Atomically load tail and tail.next
If tail.next == nil, compare-and-swap the new node onto the queue; otherwise, help move the tail pointer forward
Dequeue() (interface{}, bool):
Atomically load head and head.next
If empty, return false; otherwise, compare-and-swap the head forward and return the dequeued value
Benchmark Harness

Varies the number of producer and consumer goroutines
Each producer enqueues ops items
Each consumer dequeues ops items
Measures total time for each queue under each scenario

```

## Dependencies
```
Go 1.18+ (earlier versions likely work as well).
Imports: "fmt", "runtime", "sync", "sync/atomic", "time", "unsafe"
No external libraries beyond the Go standard library.


