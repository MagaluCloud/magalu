# Get bucket public url

## Usage:
```bash
Usage:
  ./cli object-storage buckets public-url [dst] [flags]
```

## Product catalog:
- Examples:
- ./cli object-storage buckets public-url --dst="bucket1"

## Other commands:
- Flags:
- --dst uri   Path of the bucket to generate the public url (required)
- -h, --help      help for public-url

## Flags:
```bash
Global Flags:
      --chunk-size integer     Chunk size to consider when doing multipart requests. Specified in Mb (range: 8 - 5120) (default 8)
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-ne1")
      --server-url uri         Manually specify the server to use
      --workers integer        Number of routines that spawn to do parallel operations within object_storage (min: 1) (default 5)
```
