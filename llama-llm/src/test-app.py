import sys
import os
import time
import json


sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))
from operations_api import operations

result = {}
while True:
    i = input("Enter to operate:")
    if i.isdigit():
        response = operations.request_operation(int(i))
        result[int(i)]=int(i)**2
    else:
        with open("results.json","w") as f:
            f.write(str(result))
            
