# ERROR	magalu.cloud/sdk/blueprint.network.ports	failed to create child	{"child": {"name":"","version":"","description":"","summary":"","isInternal":false,"scopes":null,"children":null,"parametersSchema":null,"configsSchema":null,"resultSchema":null,"parameters":null,"configs":null,"result":"","steps":null,"$ref":"http://magalu.cloud/sdk#/network/vpcs/ports/create"}, "error": "could not resolve \"/network/vpcs\": missing field: \"vpcs\""}
magalu.cloud/sdk/blueprint.createGroupChildren
	/home/gfz/git/mgc/magalu/mgc/sdk/blueprint/grouper.go:24
magalu.cloud/sdk/blueprint.newGrouper.func1
	/home/gfz/git/mgc/magalu/mgc/sdk/blueprint/grouper.go:13
magalu.cloud/core.(*GrouperLazyChildren[...]).getChildren
	/home/gfz/git/mgc/magalu/mgc/core/grouper_lazy_children.go:32
magalu.cloud/core.(*GrouperLazyChildren[...]).VisitChildren
	/home/gfz/git/mgc/magalu/mgc/core/grouper_lazy_children.go:53
magalu.cloud/cli/cmd.findChildByNameOrAliases
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:26
magalu.cloud/cli/cmd.loadGrouperChild
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:51
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
ERROR	magalu.cloud/sdk/blueprint.network.ports.security-group-attachment.delete	failed to resolve blueprint references	{"name": "delete", "description": "Detach a Security Group from a Port", "error": "parametersSchema: could not resolve \"/network/ports/attach\": missing field: \"attach\""}
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
Detach a Security Group from a Port

## Usage:
```bash
Usage:
  ./mgc network ports security-group-attachment delete [flags]
```

## Product catalog:
- Flags:
- -h, --help   help for delete

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text

## Flags:
```bash

```

