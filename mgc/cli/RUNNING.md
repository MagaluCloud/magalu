# Running the CLI

First download or build the `mgc` cli.

## Input handling with prefixes

When providing flags to a command, we can use one of the following prefixes to get the
parameter values:

- `@` to load file content as JSON;
- `%` to load file content as a string;
- `#` to load value directly as a string;
- No prefix to load value directly as JSON;

For example, let's say you have this file:

```
[
	{
		"id": "76e3b8c4-407f-422b-a1e9-343a40faf1cd"
	}
]
```

It follows the JSON format, so you should load it using the `@` prefix: `--flag=@filename`.

However, if your file is like this, you may want to load it as a string instead:

```
keypair_name_here
```

So use the `%` prefix: `--flag=%filename`.

This is not different from using the `#` prefix with `--flag=#keypair_name_here`, but you
can use `--flag=#1234` to make sure `1234` is read as a string instead of number.

Lastly, if you use no prefix, the value will be interpreted as JSON first then as a string
if that fails. You can use this option if you are unsure about the input format.

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

> **NOTE:**
> The file `cli.yaml` is saved at a sub-folder `mgc` of
> `$XDG_CONFIG_HOME` (on Unix, defaults to `$HOME/.config`)
> or `%AppData%` (on Windows).

## Logging

This tool uses [Uber's Zap](https://github.com/uber-go/zap) and the
command line allows filtering with [zapfilter](https://github.com/moul/zapfilter).

The logging [configuration](https://pkg.go.dev/go.uber.org/zap#Config)
can be set to customize some behavior. In the following example
we set colored level names, the timestamp in nanoseconds and show
the `file:line` of the caller:

```shell
./mgc config set --key=logging --value='
{
  "encoderConfig": {
    "levelEncoder": "color",
    "timeEncoder": "nanos",
    "timeKey": "ts",
    "callerKey": "caller"
  }
}
'
```

> **NOTE:**
> The contents must be a **valid JSON**, be careful with quotes,
> braces and commas. It's strict: trailing commas and comments are not allowed.

Then one can run commands using `--cli.log` or `-l` followed by a [pattern](https://github.com/moul/zapfilter), which takes one of the forms below:
- `levels:namespaces`
- `namespaces`

Where `levels` is a comma-separated list of level names or `*` to show all.
Note that level names are **exact**, if you want to use that level or greater,
then add the `+` suffix. That is `debug:*` logs only the debug level,
while `debug+:*` will log
`debug`, `info`, `warn`, `error`, `dpanic`, `panic` and `fatal`.

And `namespaces` is a comma-separated list of namespaces to use, it does
accept `*` as wildcard. To negate a namespace, add the `-` prefix. That
is `magalu.cloud/sdk/openapi.mke.cluster.list` matches only that namespace,
while `magalu.cloud/sdk/openapi.*` will match all the OpenAPI and
`-magalu.cloud/sdk/openapi.mke.cluster.list` will **not** show that namespace.

### Configure via environment variable

It's also possible to configure the logger using the ```MGC_LOGGING``` environment variable. The code below shows the same logger configuration as before being set thorugh the environment variable:

```shell
export MGC_LOGGING='{
  "encoderConfig": {
    "levelEncoder": "color",
    "timeEncoder": "nanos",
    "timeKey": "ts",
    "callerKey": "caller"
  }
}
'
```

The environment variable configuration has precedence over other logger configurations. This means that even if one has set a previous logger configuration using ```./mgc config set --key=logging --value=[...]``` the CLI will look for the configuration in ```MGC_LOGGING``` first.

### Examples:

```shell
# logs everything (all levels, all namespaces)
./mgc mke cluster list -l "*:*"

# logs info/warn/error/... of sdk/openapi and debug (only that level) of core/http
./mgc mke cluster list -l "info+:magalu.cloud/sdk/openapi* debug:magalu.cloud/core/http"
```

> **NOTE:**
> HTTP requests and responses are logged using the `debug` level and the
> `magalu.cloud/core/http` namespace. The sensitive bits such as
> `Authorization` headers are redacted. If you want to see the sensitive
> parts, export `MGC_SDK_LOG_SENSITIVE=1`.


## Examples

Under the folder [examples/](./examples), there are some shell scripts chaining multiple
CLI requests. For example, to create a VM, create a DISK, and attach both, run:

```shell
./examples/create-vm-with-disk.sh
```
