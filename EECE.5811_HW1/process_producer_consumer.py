import multiprocessing
import time

def producer(pipe_conn, count):
    """Produce `count` integers and send them through a pipe."""
    for i in range(count):
        pipe_conn.send(i)
    # Signal that we are done
    pipe_conn.send(None)
    pipe_conn.close()

def consumer(pipe_conn):
    """Consume messages from the pipe until None is received."""
    while True:
        item = pipe_conn.recv()
        if item is None:
            break
        # Do some simulated work if desired
    # print("Consumer done.")

def main():
    num_messages = 1000000

    # Create a pipe
    parent_conn, child_conn = multiprocessing.Pipe()

    # Create processes
    p_producer = multiprocessing.Process(target=producer, args=(child_conn, num_messages))
    p_consumer = multiprocessing.Process(target=consumer, args=(parent_conn,))

    start_time = time.time()

    # Start processes
    p_producer.start()
    p_consumer.start()

    # Wait for both to finish
    p_producer.join()
    p_consumer.join()

    elapsed = time.time() - start_time
    print(f"Process-based P/C: Sent {num_messages} messages in {elapsed:.3f} seconds")

if __name__ == "__main__":
    main()
