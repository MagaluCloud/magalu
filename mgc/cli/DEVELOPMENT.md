# Development

## Install

Install dependencies using:

```sh
go install
```

## Build

Build the command line using:

```sh
go build -o mgc # basic build without embedded openapi
go build -tags "embed" -ldflags "-s -w" -o mgc # stripped build, with embedded openapi
```

Alternatively, during development one may run without building:

```sh
go run main.go
```

> **NOTE:**
> consider using `scripts/build_release.sh` to build with recommended
> flags for all supported platforms.
