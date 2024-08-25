---
page_title: "Environment Variables MGC"
subcategory: "Guides"
description: |-
  How to use Environment Variables in MGC Cloud.
---

# Environment Variables MGC

## Introduction

This documentation describes how to configure and use the following environment variables for Terraform and the CLI:

- `MGC_API_KEY`
- `MGC_OBJ_KEY_ID`
- `MGC_OBJ_KEY_SECRET`
- `MGC_REGION`
- `MGC_ENV`

These environment variables are used for authentication and environment configuration when interacting with the provided infrastructure and services.

1. `MGC_API_KEY` - 
API key for authentication. [More information](/docs/terraform/how-to/auth#autenticação-com-api-key).

2. `MGC_OBJ_KEY_ID` - 
Key ID to access the Object Storage product. [More information](/docs/terraform/how-to/auth#autenticação-com-api-key).

3. `MGC_OBJ_KEY_SECRET` - 
*Secret* of the key to access the Object Storage product. [More information](/docs/terraform/how-to/auth#autenticação-com-api-key).

4. `MGC_REGION` - 
Specifies the region where resources will be created and managed.

5. `MGC_ENV` - 
Defines the operating environment to differentiate between different phases of development..


## Configuration in Terraform

Example:

```hcl
provider "mgc" {
  alias = "nordeste"
  region = var.mgc_region
  api_key = var.mgc_api_key
  object_storage = {
    key_pair = {
      key_id = var.mgc_obj_key_id
      key_secret = var.mgc_obj_key_secret
    }
  }
}

variable "mgc_api_key" {
  description = "API key for authentication."
}

variable "mgc_obj_key_id" {
  description = "Key ID to access the Object Storage product."
}

variable "mgc_obj_key_secret" {
  description = "Secret of the key to access the Object Storage product."
}

variable "mgc_region" {
  description = "Specifies the region where resources will be created and managed."
}

variable "mgc_env" {
  description = "Defines the operating environment"
}
```