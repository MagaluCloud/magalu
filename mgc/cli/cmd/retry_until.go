package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/spf13/cobra"
)

const (
	retryUntilFlag                string = "cli.retry-until"
	retryUntilFlagFormat          string = "retries,interval,condition"
	retryUntilFlagConditionFormat string = "engine=value"
)

func addRetryUntilFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().StringP(
		retryUntilFlag,
		"U",
		"",
		fmt.Sprintf(`Retry the action with the same parameters until the given condition is met. The flag parameters
use the format: '%s', where 'retries' is a positive integer, 'interval' is
a duration (ex: 2s) and 'condition' is a '%s' pair such as "jsonpath=expression"`,
			retryUntilFlagFormat,
			retryUntilFlagConditionFormat,
		),
	)
}

func getRetryUntilFlag(cmd *cobra.Command) (spec *core.RetryUntil, err error) {
	var v string
	v, err = cmd.Root().PersistentFlags().GetString(retryUntilFlag)
	if err != nil {
		return
	}

	return parseRetryUntilFlag(v)
}

func parseRetryUntilFlag(v string) (r *core.RetryUntil, err error) {
	if v == "" {
		return nil, nil
	}

	p := strings.SplitN(v, ",", 3)
	if len(p) != 3 {
		err = fmt.Errorf("--%s value must be in the format %s", retryUntilFlag, retryUntilFlagFormat)
		return
	}

	cfg := &core.RetryUntilConfig{}

	if _, err = fmt.Sscanf(p[0], "%d", &cfg.MaxRetries); err != nil {
		err = fmt.Errorf("--%s: failed to parse retries: %w", retryUntilFlag, err)
		return
	}

	cfg.Interval, err = time.ParseDuration(p[1])
	if err != nil {
		err = fmt.Errorf("--%s: failed to parse interval: %w", retryUntilFlag, err)
		return
	}

	condition := p[2]
	p = strings.SplitN(condition, "=", 2)
	if len(p) != 2 {
		err = fmt.Errorf("--%s condition must be in the format %s", retryUntilFlag, retryUntilFlagConditionFormat)
		return
	}
	switch p[0] {
	default:
		err = fmt.Errorf("--%s unknown condition engine: %s, supported: jsonpath|template", retryUntilFlag, p[0])
		return

	case "jsonpath":
		cfg.JSONPathQuery = p[1]
	case "template":
		cfg.TemplateQuery = p[1]
	}

	return cfg.Build()
}
