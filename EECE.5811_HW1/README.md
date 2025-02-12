# EECE.5811: HW1

> **Group Members**: Vina Dang, Layann Shaban  
> **Due Date**: February 12, 2025  

## Overview
This repository contains **two** Producer–Consumer implementations in Python:

1. **`process_producer_consumer.py`** – Uses the `multiprocessing` module  
2. **`thread_producer_consumer.py`** – Uses the `threading` module

Both demonstrate sending integers from a Producer to a Consumer using either multiple processes or multiple threads.

---

## File Descriptions

- **`process_producer_consumer.py`**  
  Implements the Producer–Consumer pattern with **multiprocessing**.  
  - **Producer** sends a range of integers via a `Pipe`.
  - **Consumer** receives integers until it encounters `None`.  
  - Prints timing for passing `num_messages` items.

- **`thread_producer_consumer.py`**  
  Implements the Producer–Consumer pattern with **threading**.  
  - **Producer** sends integers into a shared `Queue`.
  - **Consumer** dequeues items until it encounters `None`.
  - Also prints elapsed time for processing.

---

## How to Run

1. **Clone** the repository:
   ```bash
   git clone https://github.com/Lshaban313/EECE.5811_HW1.git
   cd EECE.5811_HW1

2.Run the Multiprocessing Program
```bash
python process_producer_consumer.py
```
3.Run the Threading Program

```
python thread_producer_consumer.py
```
## Design of the Programs

1.Multiprocessing (process_producer_consumer.py)
Producer
```
Generates a range of integers (e.g., for i in range(count): …)
Sends each integer via a pipe
Sends None when done
```

Consumer
```
Receives each integer from the pipe
Stops upon encountering None
main()

Creates a pipe (parent_conn, child_conn = multiprocessing.Pipe())
Spawns producer/consumer processes with multiprocessing.Process
Joins both processes and calculates total runtime
```
Threading (thread_producer_consumer.py)
Producer
```
Generates integers and places them onto a queue
Sends None at the end
```
Consumer
```
Continuously removes integers from the queue
Stops on None
main()

Creates a queue.Queue()
Spawns producer/consumer threads using threading.Thread
Joins threads and prints timing results
```
## Dependencies
```
Python 3.x
multiprocessing, threading, queue – all come standard with Python
No external libraries required

