## Magalu Cloud Provider

The MGC provider allows you to use Terraform to manage your resources on Magalu Cloud.

# Development build

## Dependencies and build

```shell
go mod tidy
make build
```

## Using local provider

```terraform
terraform {
  required_providers {
    mgc = {
      source  = "terraform.local/local/mgc"
      version = "1.0.0"
    }
  }
}
```

