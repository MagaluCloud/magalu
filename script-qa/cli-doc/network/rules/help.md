# ERROR	magalu.cloud/sdk/blueprint.network.rules	failed to create child	{"child": {"name":"","version":"","description":"","summary":"","isInternal":false,"scopes":null,"children":null,"parametersSchema":null,"configsSchema":null,"resultSchema":null,"parameters":null,"configs":null,"result":"","steps":null,"$ref":"http://magalu.cloud/sdk#/network/security_groups/rules/create"}, "error": "could not resolve \"/network/security_groups\": missing field: \"security_groups\""}
magalu.cloud/sdk/blueprint.createGroupChildren
	/home/gfz/git/mgc/magalu/mgc/sdk/blueprint/grouper.go:24
magalu.cloud/sdk/blueprint.newGrouper.func1
	/home/gfz/git/mgc/magalu/mgc/sdk/blueprint/grouper.go:13
magalu.cloud/core.(*GrouperLazyChildren[...]).getChildren
	/home/gfz/git/mgc/magalu/mgc/core/grouper_lazy_children.go:32
magalu.cloud/core.(*GrouperLazyChildren[...]).VisitChildren
	/home/gfz/git/mgc/magalu/mgc/core/grouper_lazy_children.go:53
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
magalu.cloud/cli/cmd.loadSdkCommandTree
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:132
magalu.cloud/cli/cmd.Execute
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/root.go:108
main.main
	/home/gfz/git/mgc/magalu/mgc/cli/main.go:11
runtime.main
	/usr/local/go/src/runtime/proc.go:271
ERROR	magalu.cloud/sdk/blueprint.network.rules	failed to create child	{"child": {"name":"","version":"","description":"","summary":"","isInternal":false,"scopes":null,"children":null,"parametersSchema":null,"configsSchema":null,"resultSchema":null,"parameters":null,"configs":null,"result":"","steps":null,"$ref":"http://magalu.cloud/sdk#/network/security_groups/rules/list"}, "error": "could not resolve \"/network/security_groups\": missing field: \"security_groups\""}
magalu.cloud/sdk/blueprint.createGroupChildren
	/home/gfz/git/mgc/magalu/mgc/sdk/blueprint/grouper.go:24
magalu.cloud/sdk/blueprint.newGrouper.func1
	/home/gfz/git/mgc/magalu/mgc/sdk/blueprint/grouper.go:13
magalu.cloud/core.(*GrouperLazyChildren[...]).getChildren
	/home/gfz/git/mgc/magalu/mgc/core/grouper_lazy_children.go:32
magalu.cloud/core.(*GrouperLazyChildren[...]).VisitChildren
	/home/gfz/git/mgc/magalu/mgc/core/grouper_lazy_children.go:53
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
magalu.cloud/cli/cmd.loadSdkCommandTree
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/load_cmd_tree.go:132
magalu.cloud/cli/cmd.Execute
	/home/gfz/git/mgc/magalu/mgc/cli/cmd/root.go:108
main.main
	/home/gfz/git/mgc/magalu/mgc/cli/main.go:11
runtime.main
	/usr/local/go/src/runtime/proc.go:271
VPC Rules

## Usage:
```bash
Usage:
  ./mgc network rules [flags]
```

## Product catalog:
- Flags:
- -h, --help   help for rules

## Other commands:
- Global Flags:
- --cli.show-cli-globals   Show all CLI global flags on usage text

## Flags:
```bash

```

