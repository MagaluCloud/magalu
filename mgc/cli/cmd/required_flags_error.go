package cmd

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

type requiredFlagsError []*flag.Flag

var _ error = (*requiredFlagsError)(nil)

func (e requiredFlagsError) Error() string {
	switch n := len(e); n {
	case 0:
		panic("programming error: must never return empty errors")

	case 1:
		return fmt.Sprintf("missing required flag: %s", formatFlagUsage(e[0], true))

	default:
		s := "missing required flags: "
		for i, f := range e {
			if i > 0 {
				s += ", "
			}
			s += formatFlagUsage(f, false)
		}
		return s
	}
}

func formatFlagUsage(f *flag.Flag, showHelp bool) string {
	s := fmt.Sprintf("--%s=%s", f.Name, f.Value.Type())

	if showHelp {
		if description := getFlagDescription(f); description != "" {
			s += fmt.Sprintf(" (%s)", description)
		}
	}

	return s
}
