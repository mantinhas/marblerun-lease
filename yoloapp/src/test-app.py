import sys
import os
import time


sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))
from operations_api import operations

while True:
    i = input("Enter to operate:")
    if i.isdigit():
        response = operations.request_operation(int(i))
    else:
        print("Invalid input")
