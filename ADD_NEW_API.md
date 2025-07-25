# Adding a New Specification to the Project

## Steps to Add a New Specification

### 1. Add Specification
Then add the specification:
```bash
make add-spec SPEC_URL="https://petstore3.swagger.io/api/v3/openapi.json" SPEC_NAME="pet-store"
```

### 2. Download the Specification

Download the specification:
```bash
make download-specs
```
This command downloads the specification and saves it in `mgc/spec_manipulator/cli_specs`

### 3. Prepare the Specification

Validate and indent the specification:
```bash
make prepare-specs
```

### 4. Downgrade Specification Version (Optional)

Convert specifications from version 3.1.x to 3.0.x:
```bash
make downgrade-specs
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
make refresh-specs
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
- Ensure the specification have the following structure:
```
{
  "openapi": "3.1.0", // This is the version of the specification we accept 3.1.x and 3.0.x
  "info": {
    "title": "Virtual Machine Api Product - v1",
    "description": "Virtual Machine Api Product",
    "contact": {
      "name": "IaaS Products"
    },
    "version": "v1"
  },
  "paths": {
    "/v1/images": {
      "get": {
        "tags": [
          "images"
        ],
        "summary": "Retrieves all images.",
        "description": "Retrieve a list of images allowed for the current region.",
        "operationId": "list_images_v1_v1_images_get",
        "parameters": [
          {
            "name": "_limit",
            "in": "query",
            "required": false,
            "schema": {
              "exclusiveMinimum": 0,
              "type": "integer",
              "title": " Limit",
              "maximum": 2147483647,
              "default": 50
            }
          },
          {
            "name": "_offset",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "title": " Offset",
              "maximum": 2147483647,
              "minimum": 0,
              "default": 0
            }
          },
          {
            "name": "_sort",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "title": " Sort",
              "pattern": "^(^[\\w-]+:(asc|desc)(,[\\w-]+:(asc|desc))*)?$",
              "default": "platform:asc,end_life_at:desc"
            }
          },
          {
            "name": "name",
            "in": "query",
            "description": "name of the image",
            "required": false,
            "schema": {
              "anyOf": [
                {
                  "type": "string"
                },
                {
                  "type": "null"
                }
              ],
              "title": "Name",
              "description": "name of the image"
            }
          },
          {
            "name": "availability-zone",
            "in": "query",
            "description": "br-ne1-a",
            "required": false,
            "schema": {
              "anyOf": [
                {
                  "type": "string"
                },
                {
                  "type": "null"
                }
              ],
              "title": "Availability-Zone",
              "description": "br-ne1-a"
            }
          },
          {
            "name": "x-tenant-id",
            "in": "header",
            "required": true,
            "schema": {
              "type": "string",
              "title": "X-Tenant-Id"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successful Response",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/PaginateListImageExternalV1"
                }
              }
            }
          },
          "422": {
            "description": "Validation Error",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/HTTPValidationError"
                }
              }
            }
          }
        },
        "security": [
          {
            "OAuth2": [ // OAuth2 IS REQUIRED --> go to IDM CLI to manage your product api and scopes
              "scope:read",
              "scope:write"
            ]
          }
        ]
      }
    },
  },
  "components": {
    "schemas": {
      "Architecture": {
        "type": "string",
        "title": "Architecture",
        "enum": [
          "x86/64"
        ]
      },
    }
  },
  "tags": [ // This tags is used to create the menu in the cli
    {
      "name": "instances",
      "description": "Operations with instances, including create, delete, start, stop, reboot and other actions."
    },
  ]
}