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
