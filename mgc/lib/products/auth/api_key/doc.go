/*
Package: api-key

# Summary

# Manage your ID Magalu API keys

# Description

ID Magalu API Keys are used for authentication across various platforms (CLI, SDK, Terraform, API requests). An API key has three components:

API Key: Used for Magalu API, CLI, SDK, and Terraform authentication.
Key Pair ID: Used for Object Storage authentication.
Key Pair Secret: Works with Key Pair ID for Object Storage authentication.

The API Key authenticates with the main Magalu services, while the Key Pair ID and Secret are specifically for Object Storage. Using these components correctly allows secure interaction with Magalu services and resources.

import "github.com/MagaluCloud/magalu/mgc/lib/products/auth/api_key"
*/
package apiKey
