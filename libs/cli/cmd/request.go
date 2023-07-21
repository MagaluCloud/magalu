package cmd

import (
	"fmt"
	"strings"

	"github.com/profusion/magalu/libs/parser"
	"github.com/spf13/cobra"
)

func GetRequestBody(cmd *cobra.Command, flags []*parser.Param) *strings.Reader {
	fields := []string{}

	for _, flag := range flags {
		if cmd.Flags().Changed(flag.Name) {
			value := cmd.Flags().Lookup(flag.Name).Value.String()
			fields = append(fields, fmt.Sprintf("\"%v\": %v,", flag.Name, value))
		}
	}

	body := append([]string{}, "{ ", strings.Join(fields, ""), " }")

	return strings.NewReader(strings.Join(body, ""))
}
