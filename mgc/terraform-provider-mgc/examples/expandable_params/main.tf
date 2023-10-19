terraform {
    required_providers {
        mgc = {
            version = "0.1"
            source = "magalucloud/mgc"
        }
    }
}

provider "mgc" {
    # This will be used later on to test the SDK loading functions
    apis = ["virtual-machine@1.60.0", "block-storage@1.52.0"]
}

resource "mgc_fast_items" "my_item" {
    name = "my_item"
    nested = {
        id = "8a4ab3b6-b4e5-4068-86b9-eabf171a40bf"
    }
}
