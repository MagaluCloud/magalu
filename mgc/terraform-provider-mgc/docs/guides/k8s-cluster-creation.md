---
page_title: "Requirements for K8s cluster creation"
subcategory: "Guides"
description: |-
    What is needed to create k8s cluster via Terraform in Magalu Cloud.
---

# API Key scopes required

Before creating clusters in Magalu Cloud via Terraform, make sure the used API Key has the following scopes:
- MKE: Resources creation
- ID Magalu: Manage service accounts

The "Manage service accounts" scope is required by Kubernetes so it can create needed resources in behalf of the user.

For more information on how to create an API Key with such scopes, please check the [API Key tutorial](https://docs.magalu.cloud/docs/devops-tools/api-keys/how-to/other-products/create-api-key).
