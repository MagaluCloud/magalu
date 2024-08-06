# Delete policy document for the specified bucket

## Usage:
```bash
Usage:
  ./mgc object-storage buckets policy delete [dst] [flags]
```

## Product catalog:
- Examples:
- ./mgc object-storage buckets policy delete --dst="my-bucket"

## Other commands:
- Flags:
- --dst string   Name of the bucket to delete policy file from (required)
- -h, --help         help for delete

## Flags:
```bash
Global Flags:
      --chunk-size integer     Chunk size to consider when doing multipart requests. Specified in Mb (range: 8 - 5120) (default 8)
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-se1")
      --server-url uri         Manually specify the server to use
      --workers integer        Number of routines that spawn to do parallel operations within object_storage (min: 1) (default 5)
```
