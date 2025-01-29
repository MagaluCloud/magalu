# ███╗   ███╗ ██████╗  ██████╗     ██████╗██╗     ██╗
	████╗ ████║██╔════╝ ██╔════╝    ██╔════╝██║     ██║
	██╔████╔██║██║  ███╗██║         ██║     ██║     ██║
	██║╚██╔╝██║██║   ██║██║         ██║     ██║     ██║
	██║ ╚═╝ ██║╚██████╔╝╚██████╗    ╚██████╗███████╗██║
	╚═╝     ╚═╝ ╚═════╝  ╚═════╝     ╚═════╝╚══════╝╚═╝
       
Magalu Cloud CLI is a command-line interface for the Magalu Cloud. 
It allows you to interact with the Magalu Cloud to manage your resources.

## Usage:
```bash
Usage:
  ./mgc [flags]
  ./mgc [command]
```

## Product catalog:
- Products:
- audit              Cloud Events API Product.
- block-storage      Block Storage API Product
- container-registry Magalu Container Registry product API.
- dbaas              DBaaS API Product.
- kubernetes         APIs related to the Kubernetes product.
- network            VPC Api Product
- object-storage     Operations for Object Storage
- permissions        RESTful API for permissions

## Other commands:
- virtual-machine    Virtual Machine Api Product

## Flags:
```bash
Settings:
  auth               Actions with ID Magalu to log in, API Keys, refresh tokens, change tenants and others
  config             Manage CLI Configuration values
  profile            Manage account settings, including SSH keys and related configurations
  workspace          Manage workspaces for isolated auth and config settings
```

