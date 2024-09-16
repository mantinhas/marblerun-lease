import csv
import time

def measure_request(wait=0):
    if wait!=0:
        time.sleep(wait)
    start_request_time = time.perf_counter()
    response = operations.request_operation(1)
    end_request_time = time.perf_counter()

    return end_request_time - start_request_time


def benchmark():
    data = []
    
    for i in range(0,100):
        wait = 0
        duration = measure_request(wait)
        data.append({
            "id" : i,
            "wait" : wait,
            "duration" : duration
        })
    for i in range(100,200):
        wait = 0.1
        duration = measure_request(wait)
        data.append({
            "id" : i,
            "wait" : wait,
            "duration" : duration
        })
    for i in range(200,300):
        wait = 1
        duration = measure_request(wait)
        data.append({
            "id" : i,
            "wait" : wait,
            "duration" : duration
        })

    return data

def save_data(data):
    with open("results/request_data.csv", 'a', newline='') as csvfile:
        fieldnames = ["id","wait", "duration"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)

        writer.writeheader()
        for result in data:
            writer.writerow(result)

def main():
    data = benchmark()
    save_data(data)
    
main()