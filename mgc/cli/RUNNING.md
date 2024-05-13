# Running the CLI

First download or build the `mgc` CLI. For more details on the build process
see [DEVELOPMENT.md](./DEVELOPMENT.md)

## Autocomplete

The following autocompletion scripts are provided, please copy these files to your system (`/usr`, `/usr/local` or some user-local folder):

- Fish: `share/fish/vendor_completions.d/mgc`
- Bash: `share/bash-completion/completions/mgc`
- Zsh: `share/zsh/site-functions/mgc`

> **NOTE:**
> Make sure that the `mgc` CLI is reachable via your `PATH` environment variable. It
> should be called as `mgc` instead of `./mgc` when using it.

You can then autocomplete commands by using tab completion.

## Flag usage guidelines

Running `./mgc --cli.show-cli-globals` will provide the complete set of `--cli.` flags, such as `--cli.log`,
`--cli.output` and more. These are global flags that can be used to specify certain
behaviors for each command, such as producing logs in the case of `--cli.log`. Further
instructions will be provided beside each flag. Due to how groups and executors are
handled, **for any given command, these flags should always be placed before any flags
that are specific to the command itself in order to ensure proper functionality**.
For example:
`./mgc virtual-machine instances delete --cli.log "*:*" --id ...`. Notice how
`--cli.log` is going before `--id`.

### Passing array to flags

Some commands have arrays as parameters. To pass these objects to the CLI one should pass it as a JSON enclosed by single
quotes. Like so:

```
./mgc network port ports list --port-id-list='["a", "b", "c"]'
```

There's also the option to save the JSON in a file, and pass this file as the parameter. For more information see [Input
handling with prefixes](#input-handling-with-prefixes)

Alternatively, one can specify as comma-separated values (both delimiters ';' and ',' are accepted), multiple flags are allowed.
The following are equivalent to the example above:

```
./mgc network port ports list --port-id-list=a,b,c
./mgc network port ports list --port-id-list=a --port-id-list=b,c
```

### Passing objects to flags

Some commands have objects as parameters. To pass these objects to the CLI one should pass it as a JSON enclosed by single
quotes. Like so:

```
./mgc dbaas instances create --volume '{"size":10,"type":"CLOUD_NVME"}'
```

There's also the option to save the JSON in a file, and pass this file as the parameter. For more information see [Input
handling with prefixes](#input-handling-with-prefixes)

Alternatively, one can specify `key=value` pairs as comma-separated values (both delimiters ';' and ',' are accepted), multiple flags are allowed.
The following are equivalent to the example above:

```
./mgc dbaas instances create --volume size=10,type=CLOUD_NVME
./mgc dbaas instances create --volume size=10 --volume type=CLOUD_NVME
```

## Input handling with prefixes

When providing flags to a command, we can use one of the following prefixes to get the
parameter values:

- `@` to load file content as JSON;
- `%` to load file content as a string;
- `#` to load value directly as a string;
- `help` (exact value) to show the flag's detailed usage. Use `#help` to provide the `help` string as a value.
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

## Profiles

The CLI and Terraform uses profiles to hold authentication information, and system runtime configurations. They allow
the user to save different configurations for the CLI/Terraform, and switch between them in an easy and reliable manner.

Profiles exist as directories at `$XDG_CONFIG_HOME/mgc`. In Unix systems, `$XDG_CONFIG_HOME` will default to
`$HOME/.config`, and in Windows `%AppData%` (the examples will consider a Unix environment). A profile named `foo` would
be at `$HOME/.config/mgc/foo`. The current profile is stored in the `$HOME/.config/mgc/current` file, with the name of the
profile in it. Out of the box the current profile is the `default` profile. The user is allowed to have as many profiles
as they want. All of them will use their own authentication and system configuration file: `auth.yaml` and `cli.yaml`,
respectively. Thus the directory structure at `$HOME/.config/mgc` with two profiles `foo` and `bar`, will look like this:

```
$HOME/.config/mgc/
├─ current    # Holds the name of the current profile
├─ foo/
│  ├─ auth.yaml
│  ├─ cli.yaml
├─ bar/
│  ├─ cli.yaml
│  ├─ auth.yaml
```

A set of commands under `./mgc profile` is offered to allow the user to get and set the current profile, and create,
list or delete new ones.

## Authentication

To get a token, one can run

```shell
./mgc auth login
```

A browser will open, redirecting the user to id.magalu.com. After completing authentication,
the token will be saved to `$HOME/.config/mgc/<CURRENT_PROFILE>/auth.yaml` and reused by the CLI in further actions.

> **NOTE:**
> Upon logging in, you will be able to select any of the tenant delegated to you as your primary tenant.
If you have multiple Tenants linked to your account, you can choose and
switch between them as needed:

> ```shell
> ./mgc auth tenant list # to list all avaiable tenants for current login
> ./mgc auth tenant set <id> # to set active Tenant to be used for all subsequencial requests
> ```

Execute the `select` command without id parameter to choose from
the available tenant interactively without needing to `list` them before

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
./mgc config set --key=region --value=br-se1 # change server region
./mgc config get --key=region
```

Configurations can be stored in the configuration file at `$HOME/.config/mgc/<CURRENT_PROFILE>/cli.yaml` or in environment
variables. The `config set` command saves the key value pair in the file. The `config get` command will ALWAYS check if
there's a environment variable set with the `MGC_` prefix first. This means that if there's a `foo` key in the file and
a `MGC_FOO` environment variable set, `config get --key foo` will return the environment variable value.

> **NOTE:**
> Case sensitivity is not supported for environment variables. This means
> if one stored a configuration object in an environment variable it would
> be case insensitive. Thus, it is recommended to save complex configurations
> in the configuration file using the `config set` command.

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

> **NOTE:**
> Always add the log flag before any of the command's flags, for example:
> `./mgc virtual-machine instances delete --cli.log "*:*" --id ...`. Notice
> how  `--cli.log` is going before `--id`. For more information, see the
> [flag usage guidelines](#flag-usage-guidelines).

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

It's also possible to configure the logger using the ```MGC_LOGGING``` environment variable. The code below shows the same
logger configuration as before being set thorugh the environment variable:

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

The environment variable configuration has precedence over other logger configurations. This means that even if one has
set a previous logger configuration using ```./mgc config set --key=logging --value=[...]``` the CLI will look for the
configuration in ```MGC_LOGGING``` first.

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
> parts, export `MGC_SDK_LOG_SENSITIVE=1`. To log the actual payloads
> to be sent and received, one must explicitly use
> `MGC_SDK_LOG_HTTP_PAYLOAD=progressive` (log as soon as data is received)
> or
> `MGC_SDK_LOG_HTTP_PAYLOAD=final` (accumulates all data and log at the end, when it's closed).

## Output formats

The CLI can output the requests in the following formats using the `-o` flag:

- json
- yaml
- table
- jsonpath
- template

Some commands have a default output format defined in the OpenAPI spec using an
`x-mgc` extension. The format will be overridden if the user specify a new
output format.

### JSON

```sh
./mgc virtual-machine instances list -o json

{
 "instances": [
    {
      "id": "12db59e3-8715-47af-a15f-1a595f6647ec",
      ...
    }
  ]
}

```

### YAML

```sh
./mgc virtual-machine instances list -o yaml

instances:
    - id: 12db59e3-8715-47af-a15f-1a595f6647ec
      created_at: "2023-10-04T18:00:45Z"
      error: null
      ...
```

### Table

```sh
./mgc virtual-machine instances list -o table='ID:$.instances[*].id,STATE:$.instances[*].state'

+--------------------------------------+-----------+
| ID                                   |   STATE   |
+--------------------------------------+-----------+
| 12db59e3-8715-47af-a15f-1a595f6647ec |  running  |
+--------------------------------------+-----------+
```

Table may have options (all after `=` sign). If no options are specified,
then columns will be inferred from the actual data. However, columns may
be specified in a sequence delimited by `,`, in the following format:

```
name:jsonPath
name:jsonPath:parents
name:jsonPath:parents:subTableOptions
```

Where:
- `name`: the column name (title/header) to be displayed **(required)**
- `jsonPath`: selects the data to be displayed **(required)**
- `parents`: space-delimited list of column parents, which is useful
  to group multiple columns, if they have the same parent (optional)
- `subTableOptions`: is a **quoted** string with table options
  for the inner table, in case the data is either an object or a array.
  Note that the internal JSON Path will be relative to that data,
  say the data is `{"parent": {"child": [1,2,3]}}`, if the parent
  JSON Path is `$.parent`, then this would be `$.child[*]`; analogously
  if the parent is `$.parent.child`, this would be `$[*]`.
  To make it a string, you **MUST QUOTE**
  (recommended to use the back-quote) otherwise the `,` of the internal
  columns will be handled as being of the external. (optional)

If the first column string (before `,`) is the special purpose string
`transpose`, then the table will be build along the vertical axis,
each column will become a row.

### JSONPath

```sh
./mgc virtual-machine instances list -o jsonpath='$.instances[*].id'

[
 "12db59e3-8715-47af-a15f-1a595f6647ec"
]
```

### Template

```sh
./mgc virtual-machine instances list -o template='{{range .instances}}ID:{{.id}}{{end}}'

ID:12db59e3-8715-47af-a15f-1a595f6647ec%
```

## Retry/Polling

To execute a command it is possible to define its finish condition.

For this we can use the `-U` flag, which consists of 3 parts, the number of
retries(r), the interval between attempts(i) and a condition to validate the
success(c) in their respective orders `r,i,c`.

```sh
./mgc virtual-machine instances get --id=111e3457-9013-40bf-b117-de23fc34a38c -U=3,30s,jsonpath='$.status == "active"'
... polling until the `get` request output matches the condition ...
{
 "id": "111e3457-9013-40bf-b117-de23fc34a38c",
 ...
 "status": "active",
}

./mgc virtual-machine instances get --id=111e3457-9013-40bf-b117-de23fc34a38c -U=3,30s,template='{{if eq .status "active"}}true{{end}}'
... polling until the `get` request output matches the condition ...
{
 "id": "111e3457-9013-40bf-b117-de23fc34a38c",
 ...
 "status": "active",
}
```

## State Transition

Some request also have terminate conditions already defined by the developers.
Those conditions can be integrations with websockets, callback urls, pub/sub,
retry until and similar systems.

For the commands where the terminate information is already defined, we can
provide the `-w` flag to wait for the condition.

```sh
./mgc virtual-machine instances get --id=111e3457-9013-40bf-b117-de23fc34a38c -w
... Wait until the `get` request output status be "active", "shutoff" or "error"...
{
 "id": "111e3457-9013-40bf-b117-de23fc34a38c",
 ...
 "status": "error",
}```

## Examples

Under the folder [examples/](./examples), there are some shell scripts chaining multiple
CLI requests. For example, to create a VM, create a DISK, and attach both, run:

```shell
./examples/create-vm-with-disk.sh
```

## Select Command

The `select` command is available for resources that support both the
`LIST` and `SET` commands. This command simplifies the process of setting
values by automatically listing all the available options and prompting
the user to select and set the desired value.

### Example

```sh
./mgc auth scopes select
```

This command lists all the available scopes. The user
can then select the ones they want to add. To finish the operation
the user selects the `Done` option.

```sh
./mgc auth scopes select
```

```yaml
Select multiple entries to be used with "set":
Filter: Type to filter choices
  ▸ ⏎ Done
      "block-storage.read"
      "block-storage.write"
      "cpo:read"
      "cpo:write"
      "dbaas.read"
      "dbaas.write"
      "mke.read"
      "mke.write"
      "network.read"
      "network.write"
```

> **NOTE:**
> Notice that the options that are already predefined for the user
have a `✔` in front of them. The user can then choose to keep them
or unmark them

## Links

Links are an easy and simple way of chaining operations in the cli.
Through them  you can automatically use the output of a command as parameter
for other operations.  The cli maps the fields on its own

To list the possible links from a command, simply add the `--cli.list-links` flag.

Example:

```sh
virtual-machine instances create --cli.list-links
```

It will return the following available links:

```yaml
DESCRIPTION                         NAME

 Read Virtual Machine instance       get
 Reboot a Virtual Machine instance   reboot
 Rename a Virtual Machine instance   rename
 Retype a Virtual Machine instance   retype
 Start a Virtual Machine instance    start
 Stop a Virtual Machine instance     stop
 Suspend a Virtual Machine instance  suspend
 Delete Virtual Machine instance     delete
```

To use links, you must use the character `!` followed by the link name.

Example: To create a virtual machine instance, get the data, rename it and then suspend it, use:

```yaml
virtual-machine instances create  --name <exampleName> --image.id <imageId> --machine-type.id <machineTypeId> --ssh-key-name <sshKey> ! get ! rename --name <newName> ! stop
```

In addition, for some commands that have a `get` link and that link is a
`TerminatorExecutor` you can use the `--cli.watch` flag to wait until
the operation is completed.
t's the same behavior as using `! get -w`.
