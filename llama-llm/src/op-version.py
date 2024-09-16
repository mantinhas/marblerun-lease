import torch
import time
import sys
import os
import json
import csv
import logging
import numpy as np
import llm
from llm import run_model, querys

import random

import pandas as pd

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))
from operations_api import operations


def setup():
    # Model
    for x in range(10):
        print("Warming up %d" % x)
        response = operations.request_operation(1)
        time.sleep(x * 0.1)
    return llm.setup()
    

#df = pd.read_csv("hf://datasets/fka/awesome-chatgpt-prompts/prompts.csv").to_dict()['prompt']
#querys = [ df[key] for key in sorted(df.keys()) ]

def benchmark2(model, nquerys, nloops):
    assert isinstance(nquerys, int), "Number of querys must be an integer"
    assert(nquerys<=len(querys)), f"Number of querys 'nquerys' must be smaller than available set totalling {len(querys)}"
    batch_sizes = (1,2,4,8,16,32)

    request_data = []
    model_data = []
    total_data = []

    for batch_size in batch_sizes:
        for i in range(nloops):
            start_time = time.perf_counter()
            
            for batch in range(nquerys//batch_size):
                
                start_request_time = time.perf_counter()
                response = operations.request_operation(batch_size)
                end_request_time = time.perf_counter()
                request_data.append({
                    "batch_size" : batch_size,
                    "experiment" : i,
                    "duration" : end_request_time - start_request_time,
                })

                for j in range(batch_size*batch, min(batch_size*(batch+1),nquerys), 1):
                    start_model_time = time.perf_counter()
                    run_model(model,querys[j])
                    end_model_time = time.perf_counter()

                    model_data.append(
                        {
                            "batch_size" : batch_size,
                            "experiment": i,
                            "query_id": j,
                            "duration" : end_model_time - start_model_time,
                        }
                    )
                    pretty = json.dumps(model_data[-1], indent=2)
                    print(pretty)

            end_time = time.perf_counter()

            total_data.append( {
                "environment": "epochlock_op_llama",
                "batch_size": batch_size,
                "experiment": i,
                "total_operations" : nquerys,
                "duration": end_time - start_time
            } )

    def save_as_csv(data, filename):
        with open("results/"+filename, 'a', newline='') as csvfile:
            fieldnames = list(data[0].keys())
            writer = csv.DictWriter(csvfile, fieldnames=fieldnames)

            writer.writeheader()
            for result in data:
                writer.writerow(result)

    save_as_csv(request_data, "request_data.csv")
    save_as_csv(model_data, "model_data.csv")
    save_as_csv(total_data, "total_data.csv")

def benchmark(model, nquerys, nloops):

    assert isinstance(nquerys, int), "Number of querys must be an integer"
    assert(nquerys<=len(querys)), f"Number of querys 'nquerys' must be smaller than available set totalling {len(querys)}"

    batch_sizes = (1,2,4,8,16)

    all_results = { batch_size: [] for batch_size in batch_sizes }
    request_data = []
    model_data = [ [] for x in range(nquerys)  ]
    
    logger = logging.getLogger(__name__)
    for batch_size in batch_sizes:
    
        # Benchmark
        for i in range(nloops):
            start_time = time.perf_counter()

            for batch in range(nquerys//batch_size):

                start_request_time = time.perf_counter()
                response = operations.request_operation(batch_size)
                #time.sleep(random.random())
                end_request_time = time.perf_counter()
                request_data.append(
                    {
                        "duration" : end_request_time-start_request_time,
                        "batch_size": batch_size
                    }
                )

                for j in range(batch_size*batch, batch_size*(batch+1), 1):
                    logger.info("Batch_size: %d\tLoop:%d/%d\tquery: %d/%d", (batch_size, i+1, nloops, j+1, nquerys))
                    
                    start_model_time = time.perf_counter()
                    run_model(model,querys[j])
                    end_model_time = time.perf_counter()

                    model_data[j].append(
                        {
                            "query_id": j,
                            "duration" : end_model_time - start_model_time,
                            "batch_size" : batch_size
                        }
                    )

            end_time = time.perf_counter()

            result = {}
            result["elapsed_time"] = end_time - start_time
            result["elapsed_time_per_op"] = result["elapsed_time"]/(nloops*nquerys)
            result["batch_size"] = batch_size
            result["environment_type"] = "epochlock_op_llama"
            all_results[batch_size].append(result)

    generateResultsDataCSV(all_results, batch_sizes)
    generateRequestDataCSV(request_data)
    generateModelDataCSV(model_data)

def generateResultsDataCSV(all_results, batch_sizes):
    with open('results/results.csv', 'w', newline='') as csvfile:
        fieldnames = ["environment_type","batch_size","elapsed_time","elapsed_time_per_op"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames, delimiter='\t')

        writer.writeheader()
        for batch_size in batch_sizes:
            for result in all_results[batch_size]:
                writer.writerow(result)

    with open('results/aggregate_results.csv', 'w', newline='') as csvfile:
        fieldnames = ["environment_type","batch_size"] + ["mean", "median", "std", "variance", "min", "max", "range", "percentile25", "percentile50", "percentile75"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames, delimiter='\t')

        writer.writeheader()
        for batch_size in batch_sizes:
            batch_size_results = all_results[batch_size]
            data = [measurement["elapsed_time"] for measurement in batch_size_results]

            result = { **{
                "environment_type": batch_size_results[0]["environment_type"],
                "batch_size": batch_size_results[0]["batch_size"]
            }, **generateStatisticsFromData(data)}

            writer.writerow(result)

def generateRequestDataCSV(request_data):
    with open("results/request_data.csv", 'w', newline='') as csvfile:
        fieldnames = ["duration", "batch_size"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames, delimiter='\t')

        writer.writeheader()
        for result in request_data:
            writer.writerow(result)

    with open("results/aggregate_request_data.csv", 'w', newline='') as csvfile:
        fieldnames = ["mean", "median", "std", "variance", "min", "max", "range", "percentile25", "percentile50", "percentile75"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames, delimiter='\t')

        data = [measurement["duration"] for measurement in request_data]

        result = generateStatisticsFromData(data)
        writer.writeheader()
        writer.writerow(result)



def generateModelDataCSV(model_data):
    with open("results/model_data.csv", 'w', newline='') as csvfile:
        fieldnames = ["query_id", "duration", "batch_size"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames, delimiter='\t')

        writer.writeheader()
        for i in range(len(model_data)):
            for result in model_data[i]:
                writer.writerow(result)

    with open("results/aggregate_model_data.csv", 'w', newline='') as csvfile:
        fieldnames = ["query_id"] + ["mean", "median", "std", "variance", "min", "max", "range", "percentile25", "percentile50", "percentile75"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames, delimiter='\t')

        writer.writeheader()
        for i in range(len(model_data)):
            data = [measurement["duration"] for measurement in model_data[i]]
            result = {**{
                "query_id": i
            }, **generateStatisticsFromData(data)}
            writer.writerow(result)


def generateStatisticsFromData(data):
    return {
        "mean" : round(float(np.mean(data)),4),
        "median" : round(float(np.median(data)),4),
        "std" : round(float(np.std(data)),4),
        "variance": round(float(np.var(data)),4),
        "min": round(float(np.min(data)),4),
        "max" : round(float(np.max(data)),4),
        "range" : round(float(np.ptp(data)),4),
        "percentile25": round(float(np.percentile(data, 25)),4),
        "percentile50": round(float(np.percentile(data, 50)),4),
        "percentile75": round(float(np.percentile(data, 75)),4)
    }
    

def main():
    model = setup()
    benchmark2(model, 16, 5)

if __name__ == "__main__":
    main()
