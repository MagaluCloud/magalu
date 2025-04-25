# Sync

This command uploads any file from the local path to the bucket if it is not already present or has modified time changed.

## Usage:
```
mgc object-storage objects sync [local] [bucket] [flags]
```

## Examples:
```
mgc object-storage objects sync --bucket="my-bucket/dir/" --local="./"
```

## Flags:
```
    --batch-size integer   Limit of items per batch to delete (range: 1 - 1000)
    --bucket uri           Bucket path (required)
    --delete               Deletes any item at the bucket not present on the local
-h, --help                 help for sync
    --local uri            Local path (required)
```

## Global Flags:
```
    --api-key string           Use your API key to authenticate with the API
    --chunk-size integer       Chunk size to consider when doing multipart requests. Specified in Mb (range: 8 - 5120)
-U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                               use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                               a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
-t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                               Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
    --debug                    Display detailed log information at the debug level
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. Use '--output=help' to know more details.
-r, --raw                      Output raw data, without any formatting or coloring
    --region enum              Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
    --server-url uri           Manually specify the server to use
    --workers integer          Number of routines that spawn to do parallel operations within object_storage (min: 1)
```

