# EECE.5811: HW2

> **Group Members**: Vina Dang, Layann Shaban  
> **Due Date**: February 19, 2025  

## Overview
This repository contains **two** lock implementations in Go:

1. **`TicketLock.go`** – Implements a Ticket Lock (FIFO approach)  
2. **`CasSpinLock.go`** – Implements a Compare‐and‐Swap (CAS) Spin Lock  

There is also a **`main.go`** file that runs a simple benchmark to compare each lock’s average waiting time under different numbers of goroutines.

---

## File Descriptions

1. **`TicketLock.go`**  
   Implements a **Ticket Lock** (inspired by Figure 28.7 from OSTEP).  
   - Tracks `ticket` (the next ticket number) and `turn` (the currently served ticket).  
   - Each goroutine fetches a “ticket” via an atomic Add, then spins until `turn == myTicket`.

2. **`CasSpinLock.go`**  
   Implements a **CAS** (Compare‐and‐Swap) Spin Lock.  
   - A single integer indicates lock state (0 = unlocked, 1 = locked).  
   - Threads spin until `atomic.CompareAndSwapInt32(...)` flips that integer from 0 to 1.

3. **`main.go`**  
   A benchmarking harness that:
   - Spawns a specified number of goroutines (`-g` flag).
   - Each goroutine performs a set of lock/unlock operations (`-i` flag).
   - Records total and average waiting time for each lock type.

---

## How to Run

1. **Clone** or download this repository:
   ```
   git clone https://github.com/Lshaban313/EECE.5811.git
    ```
2. Run the Benchmark:
   ```
   go run main.go
``
   Use Command-Line Flags:
```
-g <num> sets how many goroutines.
-i <num> sets how many lock/unlock cycles per goroutine.
```
## Design of the Program 

1. Ticket Lock (TicketLock.go)
```
Data:
ticket int32, turn int32.
Lock():
myTurn := atomic.AddInt32(&ticket, 1) - 1
Spin until atomic.LoadInt32(&turn) == myTurn.
Unlock():
atomic.AddInt32(&turn, 1)
```
2. CAS Spin Lock (CasSpinLock.go)
```
Data:
locked int32 (0 = unlocked, 1 = locked).
Lock():
Spin in a loop doing atomic.CompareAndSwapInt32(&locked, 0, 1).
Unlock():
atomic.StoreInt32(&locked, 0)
```
3. main.go Benchmark
```
Creates a specified number of goroutines (-g).
Each goroutine attempts -i lock/unlock operations.
Measures waiting time:
Record timestamp before Lock().
Acquire lock.
Immediately calculates elapsed time and release lock.
Collect stats and print average.
```
## Dependecies
```
Go 1.18+ (earlier versions may also work).
sync/atomic from the Go standard library is used for low‐level atomic operations.
No external Go modules are required beyond the standard library.



