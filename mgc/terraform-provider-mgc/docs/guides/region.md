---
page_title: "Configure your Terraform Region"
subcategory: "Guides"
description: |-
    How to configure the region in Terraform for Magalu Cloud.
---

# Configurar região

Para configurar a região do provedor, utilize o seguinte código:

```hcl
provider "mgc" {
  region="br-se1"
}
```

Parâmetros

region: Define a região para o provedor MGC. No exemplo acima, br-se1 indica a região Sudeste.
Exemplo

No exemplo abaixo utilizamos 2 regiões br-se1 e br-ne1 e criamos um provider para cada, adicionando um alias para facilitar o controle dos recursos de block-storage:

```hcl
terraform {
  required_providers {
    mgc = {
      source = "magalucloud/mgc"
    }
  }
}

provider "mgc" {
  alias = "sudeste"
  region = "br-se1"
}

provider "mgc" {
  alias  = "nordeste"
  region = "br-ne1"
}

resource "mgc_block-storage_volumes" "block_storage" {
    provider = mgc.nordeste
    name = "volume-via-terraform"
    size = 30
    type = {
        name = "cloud_nvme20k"
    }
}

resource "mgc_block-storage_volumes" "block_storage-sudeste" {
    provider = mgc.sudeste
    name = "volume-via-terraform"
    size = 30
    type = {
        name = "cloud_nvme20k"
    }
}
```