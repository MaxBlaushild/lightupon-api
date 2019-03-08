import yaml
import traveling_salesman

"""
This is meant to be a throwaway thin wrapper around our custom traveling salesman module. It's hardwired to use 
the example_quest.yaml file, so you can just run "python run_traveling_salesman_on_quest.py"
"""

with open("example_quest.yaml", 'r') as stream:
    try:
        quest = yaml.load(stream)
    except yaml.YAMLError as exc:
        print(exc)

quest['posts'] = traveling_salesman.put_posts_into_tsp_order(quest['posts'])

print yaml.dump(quest, default_flow_style=False)