import threading
import queue
import time

def producer(q, count):
    """Produce `count` integer messages and put them into the queue."""
    for i in range(count):
        q.put(i)
    # Signal the consumer we are done (optional technique in some designs)
    q.put(None)  # Special marker for "end of data"

def consumer(q):
    """Consume messages from the queue until None is received."""
    while True:
        item = q.get()
        if item is None:
            # End of data
            break
        # Simulate some work (optional):
        # result = item * item
    # print("Consumer done.")

def main():
    # Number of messages to send:
    num_messages = 1000000

    # Create a shared FIFO queue:
    q = queue.Queue()

    # Create producer and consumer threads
    t_producer = threading.Thread(target=producer, args=(q, num_messages))
    t_consumer = threading.Thread(target=consumer, args=(q,))

    start_time = time.time()

    # Start threads
    t_producer.start()
    t_consumer.start()

    # Wait for both to finish
    t_producer.join()
    t_consumer.join()

    elapsed = time.time() - start_time
    print(f"Thread-based P/C: Sent {num_messages} messages in {elapsed:.3f} seconds")

if __name__ == "__main__":
    main()
