#!/usr/bin/env python3
"""
EECE.5811: HW0 - Part 1
Process and Inter-Process Communication 

Group Members: Vina Dang, Layann Shaban
Date: 2/5/2025

"""

import multiprocessing

def producer(q_data, q_ack):
    """
    Producer process:
      - Generates numbers from 1 to 5
      - Prints them in the format: "Producer: X"
      - Sends each number through q_data
      - Waits for an "ACK" from consumer before continuing
    """
    for i in range(1, 6):
        # Print the produced number
        print(f"Producer: {i}")
        
        # Send the number to the consumer
        q_data.put(i)
        
        # Wait for acknowledgment 
        ack_msg = q_ack.get()
        # Assume ack_msg is "ACK" 


def consumer(q_data, q_ack):
    """
    Consumer process:
      - Receives numbers from the producer
      - Prints them in the format: "Consumer: X"
      - After each number, sends "ACK" to the producer
    """
    for _ in range(5):
        # Get the number from the producer
        num = q_data.get()
        
        # Print the consumed number
        print(f"Consumer: {num}")
        
        # Send an acknowledgment back to producer
        q_ack.put("ACK")

def main():
    """
    Main function that sets up the queues, processes, and runs them.
    """
    # Create the queues for inter-process communication
    q_data = multiprocessing.Queue()
    q_ack = multiprocessing.Queue()
    
    # Create producer and consumer processes
    p_prod = multiprocessing.Process(target=producer, args=(q_data, q_ack))
    p_cons = multiprocessing.Process(target=consumer, args=(q_data, q_ack))
    
    # Start the processes
    p_cons.start()
    p_prod.start()
    
    # Wait for both processes to finish
    p_prod.join()
    p_cons.join()

if __name__ == "__main__":
    main()

