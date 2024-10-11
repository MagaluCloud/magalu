package mock

import "github.com/spf13/cobra"

func MockCmd() *cobra.Command {

	mockCmd := &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Use:               "mock",
		Short:             "Mock",
	}

	mockCmd.AddCommand(mockCreateCmd)

	return mockCmd
}
