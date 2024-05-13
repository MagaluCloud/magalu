# ERROR	magalu.cloud/sdk/blueprint.block-storage.volume-attachment.create	failed to resolve blueprint references	{"name": "create", "description": "Attach a volume to a virtual machine instance", "error": "parametersSchema: could not resolve \"/block-storage/volumes\": missing field: \"volumes\""}
magalu.cloud/sdk/blueprint.(*executor).resolve
	/home/gfz/git/mgc/magalu/mgc/sdk/blueprint/executor.go:56
magalu.cloud/sdk/blueprint.(*executor).ParametersSchema
	/home/gfz/git/mgc/magalu/mgc/sdk/blueprint/executor.go:79
magalu.cloud/cli/cmd.newExecutorCmdFlags
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/cmd_flags.go:444
magalu.cloud/cli/cmd.addAction
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:259
magalu.cloud/cli/cmd.addChildDesc
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:18
magalu.cloud/cli/cmd.loadGrouperChild
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:56
magalu.cloud/cli/cmd.loadChild
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:70
magalu.cloud/cli/cmd.loadCommandTree
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:170
magalu.cloud/cli/cmd.loadCommandTree
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:228
magalu.cloud/cli/cmd.loadCommandTree
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:228
magalu.cloud/cli/cmd.loadSdkCommandTree
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:132
magalu.cloud/cli/cmd.Execute
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/root.go:108
main.main
	/home/gfz/git/mgc/magalu/mgc/cli/main.go:11
runtime.main
	/usr/local/go/src/runtime/proc.go:271
Attach a volume to a virtual machine instance

## Usage:
```bash
Usage:
  ./mgc block-storage volume-attachment create [flags]
```

## Product catalog:
- Flags:
- -h, --help   help for create

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text

## Flags:
```bash

```

