EECE.5811: HW0
Group Members
Vina Dang
Layann Shaban
Due Date
February 5th, 2025

Overview
This repository contains solutions for HW0, which involves:

Producer-Consumer: Two processes (in Python) that communicate using queues, ensuring the specific output order of five numbers from a Producer to a Consumer.
Stack Data Structure: A simple fixed-size stack with push and pop operations, demonstrated by a small test program.
Both files  demonstrate basic understanding of inter‐process communication, synchronization, and array‐based data structures.

File Descriptions
EECE.5811_HW0PT1.py

Purpose: Implements the Producer–Consumer pattern using Python’s multiprocessing.
Main Components:
producer(q_data, q_ack): Generates and prints numbers 1 through 5, sending each to the consumer via q_data and waiting for an “ACK.”
consumer(q_data, q_ack): Receives five numbers from q_data, prints each, and sends an “ACK” back via q_ack.
main(): Sets up two queues (q_data and q_ack), spawns the producer and consumer processes, then waits for both to finish.

EECE.5811_HW0PT2.py

Purpose: Implements a Stack class with a capacity of 100 integers. Demonstrates push and pop.
Main Components:
Stack class:
push(value): Raises an IndexError if stack is full (top >= 99).
pop(): Raises an IndexError if stack is empty (top < 0).
stack_test(): A short demo that pushes values onto the stack, pops some of them, and prints the results (in a vertical format).


How to Run
Clone this repository (or download the ZIP and extract):

git clone https://github.com/YourUserName/EECE.5811_HWO.git
cd EECE.5811_HWO
Producer-Consumer (EECE.5811_HW0PT1.py)

Ensure you have Python 3 installed.
Run the script:                                        
python EECE.5811_HW0PT1.py


Stack (EECE.5811_HW0PT2.py)
Run the script:
python EECE.5811_HW0PT2.py
