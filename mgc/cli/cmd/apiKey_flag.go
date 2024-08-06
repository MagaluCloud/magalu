package cmd

import "github.com/spf13/cobra"

const (
	apiKeyFlag = "api-key"
)

type APIKeyParameters struct {
	Key string
}

func (a APIKeyParameters) GetAPIKey() string {
	return a.Key
}

func addApiKeyFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().String(
		apiKeyFlag,
		"",
		"Use your API key to authenticate with the API",
	)
}

func getApiKeyFlag(cmd *cobra.Command) string {
	apiKey, err := cmd.Root().PersistentFlags().GetString(apiKeyFlag)
	if err != nil {
		return ""
	}
	return apiKey
}
