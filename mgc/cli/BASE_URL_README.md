# Base URL Flag for CLI Commands

`--base-url`
This flag is used to override the default host of the API.

## Use Cases
- The user wants to test the application in a different environment than the default.
- The user wants to test the application in a local environment.

## Possibilities
- The user is allowed to change the URL schema (http/https).
- The user is allowed to change the URL host.
- The user is allowed to change the URL path.

## Notes
- Query parameters are not altered.
- If you need to change only the path, you must copy the traditional schema and host, and then include the new path.
- Schema (http/https) is always required.

## Usage Examples
For example, we'll use the command: `mgc virtual-machine instances list`.
The traditional route for this command is: `https://api.magalu.cloud/br-se1/compute/v1/instances`

### A) Changing the schema:
1. If the user wants to change the schema to `http://api.magalu.cloud/br-se1/compute/v1/instances`, they can use the command:
   ```
   mgc virtual-machine instances list --base-url=http
   ```
2. The modified route will be: `http://api.magalu.cloud/br-se1/compute/v1/instances`

### B) Changing the host of the URL:
1. If the user wants to change the host to `http://localhost:8080/v1/route`, they can use the command:
   ```
   mgc virtual-machine instances list --base-url=http://localhost:8080
   ```
2. The modified route will be: `http://localhost:8080/br-se1/compute/v1/instances`

### C) Changing the path:
1. If the user wants to change the path to `https://api.magalu.cloud/br-se1/compute/v2/test-route`, they can use the command:
   ```
   mgc virtual-machine instances list --base-url=https://api.magalu.cloud/br-se1/compute/v2/test-route
   ```
2. The modified route will be: `https://api.magalu.cloud/br-se1/compute/v2/test-route`

### D) Changing the host + path:
1. If the user wants to change the host and path to `http://localhost:7777/v2/test-route`, they can use the command:
   ```
   mgc virtual-machine instances list --base-url=http://localhost:7777/v2/test-route
   ```
2. The modified route will be: `http://localhost:7777/v2/test-route`