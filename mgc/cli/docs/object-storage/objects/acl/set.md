# Set

Set ACL information for the specified object

## Usage:
```
mgc object-storage objects acl set [dst] [flags]
```

## Flags:
```
    --dst uri                     The full object URL to set the ACL information (required)
    --grant-write array(object)   Allows grantees to create objects in the bucket
                                  Use --grant-write=help for more details
-h, --help                        help for set
    --private                     Owner gets FULL_CONTROL. Delegated users have access. No one else has access rights
    --public-read                 Owner gets FULL_CONTROL. Everyone else has READ rights
```

## Global Flags:
```
    --api-key string           Use your API key to authenticate with the API
    --chunk-size integer       Chunk size to consider when doing multipart requests. Specified in Mb (range: 8 - 5120) (default 8)
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
    --workers integer          Number of routines that spawn to do parallel operations within object_storage (min: 1) (default 5)
```

