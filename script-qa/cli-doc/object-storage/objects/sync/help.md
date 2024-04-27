# This command uploads any file from the source to the destination if it's not present or has a different size. Additionally any file in the destination not present on the source is deleted.

## Usage:
```bash
Usage:
  ./cli object-storage objects sync [src] [dst] [flags]
```

## Product catalog:
- Examples:
- ./cli object-storage objects sync --dst="s3://my-bucket/dir/" --src="./"

## Other commands:
- Flags:
- --batch-size integer   Limit of items per batch to delete (range: 1 - 1000) (default 1000)
- --delete               Deletes any item at the destination not present on the source
- --dst uri              Full destination path to sync with the source path (required)
- -h, --help                 help for sync
- --src uri              Source path to sync the remote with (required)

## Flags:
```bash
Global Flags:
      --chunk-size integer     Chunk size to consider when doing multipart requests. Specified in Mb (range: 8 - 5120) (default 8)
      --cli.show-cli-globals   Show all CLI global flags on usage text
      --region enum            Region to reach the service (one of "br-mgl1", "br-ne1" or "br-se1") (default "br-ne1")
      --server-url uri         Manually specify the server to use
      --workers integer        Number of routines that spawn to do parallel operations within object_storage (min: 1) (default 5)
```

