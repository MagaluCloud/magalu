# Upload a file to a bucket

## Usage:
```bash
Usage:
  ./cli object-storage objects upload [src] [dst] [flags]
```

## Product catalog:
- Examples:
- ./cli object-storage objects upload --dst="my-bucket/dir/file.txt" --src="./file.txt"

## Other commands:
- Flags:
- --dst uri    Full destination path in the bucket with desired filename (required)
- -h, --help       help for upload
- --src file   Source file path to be uploaded (required)

## Flags:
```bash
Global Flags:
      --chunk-size integer     Chunk size to consider when doing multipart requests. Specified in Mb (range: 8 - 5120) (default 8)
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-ne1")
      --server-url uri         Manually specify the server to use
      --workers integer        Number of routines that spawn to do parallel operations within object_storage (min: 1) (default 5)
```
