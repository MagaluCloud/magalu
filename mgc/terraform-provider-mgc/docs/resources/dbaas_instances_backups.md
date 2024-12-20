---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mgc_dbaas_instances_backups Resource - terraform-provider-mgc"
subcategory: "Database"
description: |-
  Manages a DBaaS instance backup
---

# mgc_dbaas_instances_backups (Resource)

Manages a DBaaS instance backup

## Example Usage

```terraform
# Create a full backup for a DBaaS instance
resource "mgc_dbaas_instances_backups" "example" {
  instance_id = mgc_dbaas_instances.my_instance.id
  mode       = "FULL"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `instance_id` (String) ID of the DBaaS instance to backup
- `mode` (String) Mode of the backup

### Read-Only

- `id` (String) Unique identifier for the backup

## Import

Import is supported using the following syntax:

```shell
terraform import mgc_dbaas_instances_backups.example "instance-123,backup-456"
```
