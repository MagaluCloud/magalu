---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "mgc_availability_zones Data Source - terraform-provider-mgc"
subcategory: "Availability Zones"
description: |-
  List of available regions and availability zones.
---

# mgc_availability_zones (Data Source)

List of available regions and availability zones.

## Example Usage

```terraform
data "mgc_availability_zones" "availability_zones" {
}

output "availability_zones" {
  value = data.mgc_availability_zones.availability_zones.regions
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `regions` (Attributes List) (see [below for nested schema](#nestedatt--regions))

<a id="nestedatt--regions"></a>
### Nested Schema for `regions`

Read-Only:

- `availability_zones` (Attributes List) (see [below for nested schema](#nestedatt--regions--availability_zones))
- `region` (String)

<a id="nestedatt--regions--availability_zones"></a>
### Nested Schema for `regions.availability_zones`

Optional:

- `block_type` (String)

Read-Only:

- `availability_zone` (String)