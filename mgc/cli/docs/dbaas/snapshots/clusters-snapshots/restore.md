---
sidebar_position: 5
---
# Restore

Create a new cluster from snapshot.

## Usage:
```
mgc dbaas snapshots clusters-snapshots restore [cluster-id] [snapshot-id] [flags]
```

## Examples:
```
mgc dbaas snapshots clusters-snapshots restore --volume.size=30
```

## Flags:
```
    --backup-retention-days integer   Backup Retention Days: The number of days that a particular backup is kept until its deletion.
    --backup-start-at time            Backup Start At: Start time (UTC timezone) which is allowed to start the automated backup process.
    --cluster-id uuid                 Value referring to cluster Id. (required)
-h, --help                            help for restore
    --instance-type-id uuid           Instance Type Id (required)
    --name string                     Name (max character count: 100) (required)
    --snapshot-id uuid                Value referring to snapshot Id. (required)
    --volume object                   Cluster Volume Request (properties: size and type)
                                      Use --volume=help for more details
    --volume.size integer             Cluster Volume Request: The size of the volume (in GiB). (range: 20 - 50000)
                                      This is the same as '--volume=size:integer'.
    --volume.type enum                Cluster Volume Request: The type of the volume.
                                      Note: Preview volume types CLOUD_NVME30K, CLOUD_NVME40K, and CLOUD_NVME50K require tenant permission.
                                       (one of "CLOUD_NVME15K", "CLOUD_NVME20K", "CLOUD_NVME30K", "CLOUD_NVME40K" or "CLOUD_NVME50K")
                                      This is the same as '--volume=type:enum'.
```

## Global Flags:
```
    --api-key string           Use your API key to authenticate with the API
-U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                               use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                               a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
-t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                               Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
    --debug                    Display detailed log information at the debug level
    --env enum                 Environment to use (one of "pre-prod" or "prod") (default "prod")
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. Use '--output=help' to know more details.
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
```

