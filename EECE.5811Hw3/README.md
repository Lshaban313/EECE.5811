# EECE.5811: HW3

> **Group Members**: Vina Dang, Layann Shaban  
> **Due Date**: February 26, 2025  

## Overview
```
This repository contains **two** concurrent queue implementations in Go:

1. **Baseline Queue** (coarse‐grained locking)  
2. **Hand‐Over‐Hand Queue** (fine‐grained, lock‐coupling)  

Both are found in a single file called **`queue.go`**, which also includes a simple benchmarking harness to compare each queue’s performance under different workloads.

---
```
## File Descriptions
```
1. **`queue.go`**  
   This file provides:
   - **`BaselineQueue`**: a queue protected by a single lock for enqueue and dequeue (coarse‐grained).  
   - **`HOHQueue`**: a queue that uses lock‐coupling (a fine‐grained approach) where each node has its own mutex.  
   - A **benchmarking** function that measures how long it takes for multiple goroutines to perform a specified number of enqueues and dequeues.


```
## How to Run
```
1. **Clone** or download this repository:
2.  **Run the Benchmark**: Using go run queue.go

```
## Notes
```
 You can edit the `main()` function to change:  
 numGoroutinesList := []int{1, 2, 4, 8, 16} or any other parameters to explore performance under different conditions.

```
## Design of the Program
```
1. **Baseline Queue (Coarse‐Grained Locking)**
   - **Data**: 
     - A single `sync.Mutex`.
     - A dummy head node and a tail pointer.
   - **Enqueue(value int)**:
     - Lock the mutex.
     - Append a new node at the tail.
     - Unlock.
   - **Dequeue() (int, bool)**:
     - Lock the mutex.
     - Remove the node from `head.next` (if any).
     - Unlock.

2. **Hand‐Over‐Hand Queue (Lock Coupling)**
   - **Data**:
     - Each node has its own `sync.Mutex`.
     - The queue keeps pointers to a dummy `head` and `tail`.
   - **Enqueue(value int)**:
     - Lock the current tail.
     - Confirm it is still the tail (no other concurrent update changed it).
     - Link in the new node and update the tail pointer.
     - Unlock the old tail.
   - **Dequeue() (int, bool)**:
     - Lock the dummy head.
     - If `head.next != nil`, lock that node as well.
     - Remove the node, update pointers, then unlock in reverse order.

3. **Benchmark Harness**
   - Spawns multiple goroutines.
   - Each goroutine performs a series of enqueues, then dequeues.
   - Measures total time for each queue implementation.
   - Prints performance results under different concurrency levels and workloads.
  ```
   ## Dependencies

   ```
   'Go 1.18+ (earlier versions may also work). Uses "math/rand", "sync", "time" from the standard library. No external modules required.
