import torch
import time
import sys
import os
import json
import csv
import logging
import numpy as np

import random

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))
from operations_api import operations


def setup():
    # Model
    torch.hub.set_dir("torch_dir")
    model = torch.hub.load('ultralytics/yolov5', 'yolov5n', pretrained=True)
    return model

def run_model(model, arg):
    # Inference
    result = model(arg)
    # Results
    result.print()


def benchmark(model, nimages, nloops):
    imgs=[ "https://github.com/dharmx/walls/blob/main/nature/a_beach_with_waves_and_cliffs.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_beach_with_waves_and_rocks.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_body_of_water_with_a_lit_up_tower_in_the_middle.png?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_building_in_the_woods.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_bush.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_fern.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_flower.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_flower_01.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_plant.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_plant.png?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_leaves.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_leaves_01.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_leaves_02.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_moss_on_a_branch.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_pink_flowers.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_red_leaves.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_deer_with_antlers_grazing_on_grass.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_fog_over_water_with_a_black_background.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_foggy_landscape_with_trees_and_grass.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_forest_of_trees_with_fog.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_forest_with_moss_and_trees.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_group_of_bamboo_trees.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_group_of_bamboo_trees_01.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_group_of_black_leaves.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_group_of_dark_leaves.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_group_of_fish_swimming_in_a_pond.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_group_of_green_leaves.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_group_of_pink_tulips.png?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_group_of_trees_with_green_leaves.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_group_of_white_flowers.png?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_house_in_the_woods.png?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_lake_surrounded_by_trees.png?raw=true",
            "https://github.com/dharmx/walls/blob/main/nature/a_lake_with_snow_covered_mountains_in_the_background.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_lake_with_trees_and_clouds.png?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_large_rocks_in_a_desert.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_log_in_a_body_of_water_surrounded_by_trees.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_mountain_range_with_snow_on_top.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_pile_of_logs_in_a_snowy_field.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_pile_of_rocks.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_pond_with_lily_pads_and_a_wooden_fence.png?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_river_between_rocks_with_trees.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_river_running_through_a_forest.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_river_running_through_a_forest_01.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_river_running_through_a_rocky_canyon.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_road_with_lights_on_the_side_of_a_body_of_water.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_road_with_trees_in_the_background.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_rocky_beach_with_water_in_the_background.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_rocky_shore_with_waves_crashing_on_rocks.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_sand_dunes_in_the_desert.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_small_plant_in_a_pot.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_snowy_mountain_with_a_body_of_water.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_stairs_leading_to_a_rocky_canyon_with_Flume_Gorge_in_the_background.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_stone_stairs_in_a_forest.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_stream_in_the_woods.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_tall_mountain_with_clouds.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_train_tracks_in_a_tunnel.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_tree_with_pink_flowers.png?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_waterfall_over_rocks.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/a_wooden_stairs_in_a_forest.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/plants_leaves_in_the_water.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/rocks_in_the_water.jpg?raw=true", "https://github.com/dharmx/walls/blob/main/nature/train_tracks_in_a_forest.jpg?raw=true", 
            "https://github.com/dharmx/walls/blob/main/nature/waves_crashing_waves_on_rocks.jpg?raw=true", "https://www.nbc.com/sites/nbcblog/files/2022/07/the-office-how-to-watch.jpg" ]

    assert isinstance(nimages, int), "Number of images must be an integer"
    assert(nimages<=len(imgs)), f"Number of images 'nimages' must be smaller than available set totalling {len(imgs)}"

    batch_sizes = (1,2,4,8,16,32)

    all_results = { batch_size: [] for batch_size in batch_sizes }
    request_data = []
    model_data = [ [] for x in range(nimages)  ]
    
    logger = logging.getLogger(__name__)
    for batch_size in batch_sizes:
    
        # Benchmark
        for i in range(nloops):
            start_time = time.perf_counter()

            for batch in range(nimages//batch_size):

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
                    logger.info("Batch_size: %d\tLoop:%d/%d\tImage: %d/%d", (batch_size, i+1, nloops, j+1, nimages))
                    
                    start_model_time = time.perf_counter()
                    run_model(model,imgs[j])
                    end_model_time = time.perf_counter()

                    model_data[j].append(
                        {
                            "image_id": j,
                            "duration" : end_model_time - start_model_time,
                            "batch_size" : batch_size
                        }
                    )

            end_time = time.perf_counter()

            result = {}
            result["elapsed_time"] = end_time - start_time
            result["elapsed_time_per_op"] = result["elapsed_time"]/(nloops*nimages)
            result["batch_size"] = batch_size
            result["environment_type"] = "epochlock_op"
            all_results[batch_size].append(result)

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


    with open("results/model_data.csv", 'w', newline='') as csvfile:
        fieldnames = ["image_id", "duration", "batch_size"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames, delimiter='\t')

        writer.writeheader()
        for i in range(len(model_data)):
            for result in model_data[i]:
                writer.writerow(result)

    with open("results/aggregate_model_data.csv", 'w', newline='') as csvfile:
        fieldnames = ["image_id"] + ["mean", "median", "std", "variance", "min", "max", "range", "percentile25", "percentile50", "percentile75"]
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames, delimiter='\t')

        writer.writeheader()
        for i in range(len(model_data)):
            data = [measurement["duration"] for measurement in model_data[i]]
            result = {**{
                "image_id": i
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
    benchmark(model, 64, 5)

if __name__ == "__main__":
    main()
