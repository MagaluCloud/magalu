---
page_title: "Persisting Terraform State in MGC"
subcategory: "Guides"
description: |-
    How to persist Terraform state in MGC Cloud using object storage.
---

# Using S3 Backend for Terraform State in Magalu Cloud

This guide explains how to configure Terraform to use the S3 backend for storing the Terraform state file (`terraform.tfstate`) in Magalu Cloud using its S3-compatible object storage. This approach enhances your Terraform projects with better state management, collaboration, and security.

## Prerequisites

- Terraform installed on your local machine.
- Access to Magalu Cloud with permissions to create and manage object buckets.
- An object bucket created in Magalu Cloud for storing the Terraform state file.

## Configuration Steps

1. **Create an S3 Bucket in Magalu Cloud:**

   Before configuring Terraform, ensure you have an object bucket in Magalu Cloud. This bucket will be used to store the Terraform state files.

2. **Obtain Your Access Key and Secret Key:**

   For Terraform to access your bucket, you need an access key and a secret key from Magalu Cloud. These keys are used for authentication.

3. **Configure the Terraform Backend:**

   In your Terraform configuration file (`main.tf`), specify the S3 backend settings under the `backend "s3"` block. Here's an example configuration tailored for Magalu Cloud's S3 implementation:

   ```terraform
   terraform {
     required_providers {
       mgc = {
         source = "magalucloud/mgc"
       }
     }
     backend "s3" {
       bucket                      = "your-bucket-name"
       key                         = "terraform.tfstate"
       secret_key                  = "your-secret-key"
       access_key                  = "your-access-key"
       region                      = "your-region"
       skip_region_validation      = true
       skip_credentials_validation = true
       skip_requesting_account_id  = true
       skip_s3_checksum            = true
       endpoints {
         s3 = "https://your-region.magaluobjects.com/"
       }
     }
   }
    ```
