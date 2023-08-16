# Running the CLI

First download or build the `mgc` cli.

## Authentication

To get a token, one can run

```shell
./mgc auth login
```

A browser will open, redirecting the user to id.magalu.com. After completing authentication,
the token will be saved to `$HOME/.config/mgc/auth.yaml` and reused by the CLI in further actions:

To ensure it is working, perform a CLI command that requires authentication:

```shell
./mgc virtual-machine instances list
```

> **NOTE:**
> one can still use the env var to override the value of the token by setting:
> ```shell
> export MGC_SDK_ACCESS_TOKEN=""
> ```

## Configuration

```shell
./mgc config list  # list all known config and their types
./mgc config set --key=region --value=br-ne-2 # change server region
./mgc config get --key=region
```

## Examples

Under the folder [examples/](./examples), there are some shell scripts chaining multiple
CLI requests. For example, to create a VM, create a DISK, and attach both, run:

```shell
./examples/create-vm-with-disk.sh
```
