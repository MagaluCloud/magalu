---
sidebar_position: 0
---
# Objects

Object operations for Object Storage API

## Usage:
```
mgc object-storage objects [flags]
mgc object-storage objects [command]
```

## Commands:
```
acl          ACL related operations
copy         Copy an object from a bucket to another bucket
copy-all     Copy all objects from a bucket to another bucket
delete       Delete an object from a bucket
delete-all   Delete all objects from a bucket
download     Download an object from a bucket
download-all Download all objects from a bucket
head         Get object metadata
list         List all objects from a bucket
move         Moves one object from source to destination
move-dir     Moves objects from source to destination
object-lock  Object locking commands
presign      Generate a pre-signed URL for accessing an object
public-url   Get object public url
sync         Synchronizes a local path with a bucket
upload       Upload a file to a bucket
upload-dir   Upload a directory to a bucket
versions     Retrieve all versions of an object
```

## Flags:
```
-h, --help   help for objects
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
    --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
-o, --output string            Change the output format. Use '--output=help' to know more details.
-r, --raw                      Output raw data, without any formatting or coloring
```

