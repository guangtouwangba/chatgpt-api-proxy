import yaml
import os
def convert_yaml_to_env(yaml_file, env_file):
    with open(yaml_file, 'r') as f:
        yaml_data = yaml.safe_load(f)

    with open(env_file, 'w') as f:
        for key, value in flatten_dict(yaml_data).items():
            line = f"{key.upper()}={value}\n"
            f.write(line)

def flatten_dict(d, parent_key='', sep='_'):
    items = []
    for k, v in d.items():
        new_key = f"{parent_key}{sep}{k}" if parent_key else k
        if isinstance(v, dict):
            items.extend(flatten_dict(v, new_key, sep=sep).items())
        else:
            items.append((new_key, v))
    return dict(items)

# read yaml file from parent directory ../../conf
yaml_file = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'conf', 'config-prod.yml')
# write env file to parent directory ../../
env_file = os.path.join(os.path.dirname(os.path.dirname(__file__)), '.env')

convert_yaml_to_env(yaml_file, env_file)
