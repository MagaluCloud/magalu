# Internal route for update status of a instance when receive a update from URP.

## Usage:
```bash
### Note
This route is used only for internal proposes.
```

## Product catalog:
- Usage:
- ./mgc virtual-machine-xaas instances internal-instances urp update [urp-id] [flags]

## Other commands:
- Flags:
- --error string              Error
- -h, --help                      help for update
- --instance-type-id string   Instance Type Id
- --state enum                InstanceV1State (one of "deleted", "new", "running", "stopped" or "suspended")
- --status enum               InstanceV1Status (one of "attach_nic_pending", "attaching_nic", "completed", "creating", "creating_error", "creating_error_quota", "creating_error_quota_disk", "creating_error_quota_floating_ip", "creating_error_quota_instance", "creating_error_quota_ram", "creating_error_quota_vcpu", "creating_network_error", "deleted", "deleting", "deleting_error", "deleting_pending", "detach_nic_pending", "detaching_nic", "provisioning", "rebooting", "rebooting_pending", "retyping", "retyping_confirmed", "retyping_error", "retyping_error_quota", "retyping_pending", "starting", "starting_pending", "stopping", "stopping_pending", "suspending" or "suspending_pending")
- --urp-id string             Urp Id (required)
- -v, --version                   version for update

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --env enum               Environment to use (one of "pre-prod" or "prod") (default "prod")
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri         Manually specify the server to use
```

