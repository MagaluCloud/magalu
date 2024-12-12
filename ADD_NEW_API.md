# Adding a New Specification to the Project

## Steps to Add a New Specification

### 1. Add Specification to spec_manipulator/specs.yaml

Use the following command:
```bash
./mgc/spec_manipulator/build.sh
```

Then add the specification:
```bash
./mgc/spec_manipulator/specs add https://petstore3.swagger.io/api/v3/openapi.json pet-store
```

### 2. Download the Specification

Download the specification:
```bash
./mgc/spec_manipulator/specs download
```
This command downloads the specification and saves it in `mgc/spec_manipulator/cli_specs`

### 3. Prepare the Specification

Validate and indent the specification:
```bash
./mgc/spec_manipulator/specs prepare
```

### 4. Downgrade Specification Version (Optional)

Convert specifications from version 3.1.x to 3.0.x:
```bash
./mgc/spec_manipulator/specs downgrade
```
This will generate a new specification with the "conv." prefix in the filename.

## Updating add_all_specs.sh Script

### Specification Addition Guidelines

- Include your specification, paying attention to the version
- If the original specification is v3.1.x, add the converted specification with the "conv." prefix
- Follow the existing pattern:
  ```bash
  $BASEDIR/add_specs***.sh NAME_IN_MENU URL_PATH SPEC_LOCAL_PATH UNIQUE_URL
  ```

### Adding Specifications Based on Scope

#### Regional Specifications
```bash
$BASEDIR/add_specs.sh audit audit mgc/spec_manipulator/cli_specs/conv.events-consult.openapi.yaml https://events-consult.jaxyendy.com/openapi-cli.json
```

#### Global Specifications
```bash
$BASEDIR/add_specs_without_region.sh profile profile mgc/spec_manipulator/cli_specs/conv.globaldb.openapi.yaml https://globaldb.jaxyendy.com/openapi-cli.json
```

### Final Step

Execute the script to finalize specification integration:
```bash
./scripts/add_all_specs.sh
```

## Output

After completing the script, two new files will be created in these directories:
- `openapi-customizations`
- `mgc/sdk/openapi/openapis`

## Build and Availability

Once the process is complete:
- CLI
- Terraform
- Library can be built
- The new API will be available for use

## Notes
- Ensure you follow the specified naming conventions
- Pay attention to the specification's version and scope (regional or global)
- Use the appropriate script for adding specifications