# ERROR	magalu.cloud/sdk/blueprint.network.public_ips	failed to create child	{"child": {"name":"","version":"","description":"","summary":"","isInternal":false,"scopes":null,"children":null,"parametersSchema":null,"configsSchema":null,"resultSchema":null,"parameters":null,"configs":null,"result":"","steps":null,"$ref":"http://magalu.cloud/sdk#/network/vpcs/public-ips/create"}, "error": "could not resolve \"/network/vpcs\": missing field: \"vpcs\""}
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
ERROR	magalu.cloud/sdk/blueprint.network.public_ips.port-attachment.create	failed to resolve blueprint references	{"name": "create", "description": "Attach a Public IP to a Port", "error": "parametersSchema: could not resolve \"/network/public_ips/attach\": missing field: \"attach\""}
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
magalu.cloud/cli/cmd.loadAllGrouperChildren.func1
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:82
magalu.cloud/core.(*GrouperLazyChildren[...]).VisitChildren
	/home/gfz/git/mgc/magalu/mgc/core/grouper_lazy_children.go:59
magalu.cloud/cli/cmd.loadAllGrouperChildren
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:77
magalu.cloud/cli/cmd.loadAllChildren
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:94
magalu.cloud/cli/cmd.loadCommandTree
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:159
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
ERROR	magalu.cloud/sdk/blueprint.network.public_ips.port-attachment.delete	failed to resolve blueprint references	{"name": "delete", "description": "Detach a Public IP from a Port", "error": "parametersSchema: could not resolve \"/network/public_ips/attach\": missing field: \"attach\""}
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
magalu.cloud/cli/cmd.loadAllGrouperChildren.func1
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:82
magalu.cloud/core.(*GrouperLazyChildren[...]).VisitChildren
	/home/gfz/git/mgc/magalu/mgc/core/grouper_lazy_children.go:59
magalu.cloud/cli/cmd.loadAllGrouperChildren
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:77
magalu.cloud/cli/cmd.loadAllChildren
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:94
magalu.cloud/cli/cmd.loadCommandTree
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:159
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
ERROR	magalu.cloud/sdk/blueprint.network.public_ips.port-attachment.get	failed to resolve blueprint references	{"name": "get", "description": "Check if a Public IP is attached to a Port", "error": "parametersSchema: could not resolve \"/network/public_ips/attach\": missing field: \"attach\""}
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
magalu.cloud/cli/cmd.loadAllGrouperChildren.func1
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:82
magalu.cloud/core.(*GrouperLazyChildren[...]).VisitChildren
	/home/gfz/git/mgc/magalu/mgc/core/grouper_lazy_children.go:59
magalu.cloud/cli/cmd.loadAllGrouperChildren
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:77
magalu.cloud/cli/cmd.loadAllChildren
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:94
magalu.cloud/cli/cmd.loadCommandTree
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:159
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
Manage the attachment between a Public IP and a Port

## Usage:
```bash
Usage:
  ./mgc network public-ips port-attachment [flags]
  ./mgc network public-ips port-attachment [command]
```

## Product catalog:
- Commands:
- create      Attach a Public IP to a Port
- delete      Detach a Public IP from a Port
- get         Check if a Public IP is attached to a Port

## Other commands:
- Flags:
- -h, --help   help for port-attachment

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
```

