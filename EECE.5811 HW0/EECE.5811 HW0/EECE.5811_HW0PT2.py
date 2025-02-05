#!/usr/bin/env python3
"""
EECE.5811: HW0 - Part 2
Stack Implementation (Array of size 100)

Group Members: Vina Dang, Layann Shaban
Date: 2/5/2025
"""

class Stack:
    """
    A simple stack class that can hold up to 100 integers.
    """
    
    def __init__(self):
        """
        Constructor: Initializes the stack with an array of size 100
        and sets top to -1 indicating an empty stack.
        """
        self.data = [None] * 100  # fixed-size array
        self.top = -1             # top index, -1 means empty

    def push(self, value):
        """
        Pushes an integer onto the stack.
        Raises an IndexError if the stack is already full.
        """
        if self.top >= 99:
            raise IndexError("Stack overflow: cannot push onto a full stack.")
        self.top += 1
        self.data[self.top] = value

    def pop(self):
        """
        Pops and returns the top integer from the stack.
        Raises an IndexError if the stack is empty.
        """
        if self.top < 0:
            raise IndexError("Stack underflow: cannot pop from an empty stack.")
        value = self.data[self.top]
        self.data[self.top] = None  # Clean Data Up
        self.top -= 1
        return value

def stack_test():
    """
    Demonstrates basic usage of the Stack class by pushing and popping a few values.
    Prints the values pushed/popped in a vertical format (one per line).
    """
    print("\n--- Stack Test Start ---")
    
    # Create a stack
    s = Stack()
    
    # Define a list of integers we plan to push
    values_to_push = [10, 20, 30]
    
    print("Pushed the following values onto the stack :")
    for val in values_to_push:
        s.push(val)
        print(val)
    
    # Pop two values from the stack
    popped_val1 = s.pop()
    popped_val2 = s.pop()
    print("\nPopped the following values :")
    print(popped_val1)
    print(popped_val2)
    
    print("--- Stack Test End ---\n")

def main():
    """
    Main function to run the stack test.
    """
    stack_test()

if __name__ == "__main__":
    main()
