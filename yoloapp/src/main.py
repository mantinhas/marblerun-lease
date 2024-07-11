import torch

imgs=[
    "https://github.com/dharmx/walls/blob/main/nature/a_beach_with_waves_and_cliffs.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_beach_with_waves_and_rocks.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_body_of_water_with_a_lit_up_tower_in_the_middle.png?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_building_in_the_woods.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_bush.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_fern.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_flower.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_flower_01.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_plant.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_a_plant.png?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_leaves.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_leaves_01.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_leaves_02.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_moss_on_a_branch.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_pink_flowers.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_close_up_of_red_leaves.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_deer_with_antlers_grazing_on_grass.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_fog_over_water_with_a_black_background.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_foggy_landscape_with_trees_and_grass.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_forest_of_trees_with_fog.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_forest_with_moss_and_trees.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_group_of_bamboo_trees.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_group_of_bamboo_trees_01.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_group_of_black_leaves.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_group_of_dark_leaves.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_group_of_fish_swimming_in_a_pond.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_group_of_green_leaves.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_group_of_pink_tulips.png?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_group_of_trees_with_green_leaves.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_group_of_white_flowers.png?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_house_in_the_woods.png?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_lake_surrounded_by_trees.png?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_lake_with_snow_covered_mountains_in_the_background.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_lake_with_trees_and_clouds.png?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_large_rocks_in_a_desert.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_log_in_a_body_of_water_surrounded_by_trees.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_mountain_range_with_snow_on_top.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_pile_of_logs_in_a_snowy_field.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_pile_of_rocks.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_pond_with_lily_pads_and_a_wooden_fence.png?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_river_between_rocks_with_trees.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_river_running_through_a_forest.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_river_running_through_a_forest_01.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_river_running_through_a_rocky_canyon.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_road_with_lights_on_the_side_of_a_body_of_water.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_road_with_trees_in_the_background.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_rocky_beach_with_water_in_the_background.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_rocky_shore_with_waves_crashing_on_rocks.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_sand_dunes_in_the_desert.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_small_plant_in_a_pot.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_snowy_mountain_with_a_body_of_water.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_stairs_leading_to_a_rocky_canyon_with_Flume_Gorge_in_the_background.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_stone_stairs_in_a_forest.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_stream_in_the_woods.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_tall_mountain_with_clouds.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_train_tracks_in_a_tunnel.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_tree_with_pink_flowers.png?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_waterfall_over_rocks.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/a_wooden_stairs_in_a_forest.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/plants_leaves_in_the_water.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/rocks_in_the_water.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/train_tracks_in_a_forest.jpg?raw=true",
    "https://github.com/dharmx/walls/blob/main/nature/waves_crashing_waves_on_rocks.jpg?raw=true"
]

# Model
model = torch.hub.load('ultralytics/yolov5', 'yolov5s', pretrained=True)

# Inference
for x in range(5):
    result = model(imgs[x])
    # Results
    result.print()

