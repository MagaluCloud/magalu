# ID Magalu API Keys are used for authentication across various platforms (CLI, SDK, Terraform, API requests). An API key has three components:

## Usage:
```bash
API Key: Used for Magalu API, CLI, SDK, and Terraform authentication.
Key Pair ID: Used for Object Storage authentication.
Key Pair Secret: Works with Key Pair ID for Object Storage authentication.
```

## Product catalog:
- The API Key authenticates with the main Magalu services, while the Key Pair ID and Secret are specifically for Object Storage. Using these components correctly allows secure interaction with Magalu services and resources.

## Other commands:
- Usage:
- mgc auth api-key [flags]
- mgc auth api-key [command]

## Flags:
```bash
Commands:
  create      Create a new API Key
  get         Get a specific API key by its ID
  list        List your account API keys
  revoke      Revoke an API key by its ID
```

