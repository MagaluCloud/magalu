# Script to merge interal and external OpenAPI spec
# Replace external requestBody with internal, and save the new
# spec in a .yaml file
# 
# How to run:
# 
#   python3 fix_openapi.py internal_path external_path destination_path
# 
#   - internal_path: path to internal OpenAPI specification
#   - external_path: path to external OpenAPI specification
#   - destination_path: path to destination where the new OpenAPI spec will be saved


import yaml
import warnings
import sys

def fix_openapi(internal_path: str, external_path: str):
    internal = yaml.load(open(internal_path), Loader=yaml.CLoader)
    external = yaml.load(open(external_path), Loader=yaml.CLoader)
    final = {}

    final['openapi'] = external['openapi']
    final['info'] = external['info']
    final['paths'] = external['paths']

    for path in external['paths']:
        if path not in internal['paths']:
            warnings.warn(f"Path {path} not present in internal YAML", category=RuntimeWarning)
            return

        for action in external['paths'][path]:
            if action not in internal['paths'][path]:
                warnings.warn(f"Action {path} not present in internal YAML", category=RuntimeWarning)
                return
            
            if "requestBody" in external['paths'][path][action]:
                final['paths'][path][action]['requestBody'] = internal['paths'][path][action]['requestBody']
    
    final['components'] = external['components']
    final['tags'] = external['tags']
    final['servers'] = external['servers']

    with open(destination, 'w') as file:
        yaml.dump(
            yaml.safe_load(str(final)),
            file,
            sort_keys=False,
            indent=4
        )

internal_file = sys.argv[1]
external_file = sys.argv[2]
destination   = sys.argv[3]

fix_openapi(internal_file, external_file)
