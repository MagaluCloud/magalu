# User Guide: Installing and Using a Magalu Cloud Provider in Terraform

## Introduction

This guide provides step-by-step instructions on how to install and use a Magalu provider in Terraform to integrate resources into your infrastructure environment. Follow these steps carefully to ensure a successful integration.

### Prerequisites

Before you begin, make sure you have the following prerequisites installed:

- [Terraform](https://www.terraform.io/downloads.html)

or

- [OpenTofu](https://opentofu.org/docs/intro/install/)

## Installing the Magalu Provider

1.  **Download the Magalu Provider:**

    - Download the Magalu provider or compile it from the source code.
      - When compiling, make sure to add the `embed` value to the `-tags` flag, like so: `go build -tags embed`
      - When compliing, also make sure to export the `MGC_SDK_OPENAPI_DIR` and `MGC_SDK_BLUEPRINTS_DIR` variables
        pointing to the correct directories where the OpenAPI and Blueprint files are located. In this repo, they
        are under the `mgc/cli/openapis` and `mgc/cli/blueprints` directories, respectfully

2.  **Install in Terraform Directory:**

    > **NOTE:**
    > To use OpenTofu set the environment variable `MGC_OPENTF` to one before run install script
    >
    > ```shell
    > export MGC_OPENTF=1
    > ```

    - Choose the correct binary corresponding to your operating system, copy both the `binary` and `install.sh` to `~/terraform-mgc`.

    - Run the `install.sh` script from the same directory as your binary. This action will establish an override, directing Terraform to search for the provider within the local environment (this directory) rather than in the remote registry.

      ```sh
      ./install.sh
      ```

    > **NOTE:**
    > Executing a terraform plan/apply might result in the following warning message:
    >
    > "Warning: Provider development overrides are in effect"
    >
    > It is safe to disregard it for the time being.

## Configuring the Provider in Terraform

1. **Create a New Project Directory:**

   - Create a new directory for your project and navigate to it in the terminal.

2. **Initialize a New Terraform Configuration File:**

   - Start a new Terraform configuration file (e.g., `main.tf`) and add the following configuration for the Magalu provider:

     ```hcl
     terraform {
        required_providers {
            mgc = {
               version = "0.1"
               source = "magalucloud/mgc"
            }
        }
     }

     provider "mgc" {}
     ```

## Configuring Resources

1. **Add Resources Specific to the Magalu Provider:**

   - Add resources specific to the Magalu provider to your Terraform configuration file. Refer to the provider's documentation for details on supported resources.

     ```hcl
     # example
     resource "mgc_virtual-machine_instances" "myvm" {
        name = "my-tf-vm"

        machine_type = {
            name = "cloud-bs1.xsmall"
        }

        image = {
            name = "cloud-ubuntu-22.04 LTS"
        }

        ssh_key_name = "luizalabs-key"

        availability_zone = "br-ne-1c"
     }
     ```

2. **Magalu Resource Parameters:**
   - Magalu the resource parameters as needed for your environment.

## Initialization and Application

1. **Execute Commands in the Terminal:**

   - Execute the following commands in the terminal to initialize and plan the configurations:

     ```sh
     terraform plan # Preview changes before applying.

     # or using OpenTofu
     tofu plan
     ```

   - Apply your configuration

     ```sh
     terraform apply

     # or using OpenTofu
     tofu apply
     ```

2. **Confirm Changes:**
   - Confirm the application of changes when prompted.

## Verification and Management

1. **Check if Resources Were Created:**

   - Use the following command to check if resources were created successfully:

     ```sh
     terraform show

     # or using OpenTofu
     tofu show
     ```

2. **Modify or Remove Resources:**
   - Use commands such as `terraform apply`, `terraform destroy`, and others as needed to modify or remove resources.
