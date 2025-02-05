 EECE.5811: HW0  
**Group Members**: Vina Dang, Layann Shaban  
**Due Date**: February 5, 2025  

## Overview
This repository contains solutions for HW0, which involves two parts:
1. Producer–Consumer using Python’s `multiprocessing`.
2. A simple Stack data structure with push/pop.

## File Descriptions
- **EECE.5811_HW0PT1.py**: Implements the Producer–Consumer pattern.
- **EECE.5811_HW0PT2.py**: Implements a 100-element Stack with push and pop, plus a demo test.

## How to Run
1. **Clone** the repository:
   ```bash```
   git clone https://github.com/Lshaban313/EECE.5811_HWO.git
   cd EECE.5811_HWO

2. Run the Producer–Consumer Program
 ```bash```
  python EECE.5811_HW0PT1.py


4. Run the Stack Program
```bash```
python EECE.5811_HW0PT2.py

## Design of the Programs
 ```bash```
 Producer–Consumer (EECE.5811_HW0PT1.py)
 Producer:
Generates numbers 1 through 5.
Prints “Producer: X”.
Sends each number to the Consumer via a queue.
Waits for an “ACK” after each send.
Consumer:
Receives each number from the Producer.
Prints “Consumer: X”.
Sends an “ACK” back to the Producer.
main():
Creates two queues (q_data and q_ack) for communication.
Spawns Producer and Consumer processes with multiprocessing.Process.
Waits for both to complete using .join().
