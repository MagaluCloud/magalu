package cmd

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
	"github.com/spf13/cobra"
)

func newDumpTreeCmd(sdk *mgcSdk.Sdk) *cobra.Command {
	return &cobra.Command{
		Use:     "dump-tree",
		Short:   "Print command tree",
		Long:    `Walks through the command tree, and prints name, description, version, children and schema for parameters and configs. Defaults to YAML output, but "-o json" and other formats may be used`,
		Hidden:  true,
		GroupID: "other",
		RunE: func(cmd *cobra.Command, args []string) error {
			root := sdk.Group()

			tree, err := collectAllChildren(root)
			if err != nil {
				return err
			}

			output := getOutputFlag(cmd)
			if output == "" {
				output = "yaml"
			}
			name, options := parseOutputFormatter(output)
			formatter, err := getOutputFormatter(name, options)
			if err != nil {
				return err
			}

			return formatter.Format(tree["children"], options, getRawOutputFlag(cmd))
		},
	}
}

func collectAllChildren(child core.Descriptor) (map[string]any, error) {
	node := map[string]any{}
	node["name"] = child.Name()
	node["description"] = child.Description()
	node["isInternal"] = child.IsInternal()
	node["groupId"] = child.GroupID()

	if executor, ok := child.(core.Executor); ok {
		node["parameters"] = executor.ParametersSchema()
		node["configs"] = executor.ConfigsSchema()
		node["result"] = executor.ResultSchema()

		return node, nil
	} else if grouper, ok := child.(core.Grouper); ok {
		children := []map[string]any{}
		_, err := grouper.VisitChildren(func(child core.Descriptor) (run bool, err error) {
			c, err := collectAllChildren(child)
			if err != nil {
				return false, err
			}
			children = append(children, c)

			return true, nil
		})

		if err != nil {
			return nil, fmt.Errorf("unable to visit all children from node %s: %w", child.Name(), err)
		}

		node["children"] = children

		return node, nil
	} else {
		return nil, fmt.Errorf("child %v not group/executor", child)
	}
}
