package cmd

import (
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

type ColorScheme struct {
	headers        *color.Color
	commands       *color.Color
	cmdDescription *color.Color
	execName       *color.Color
	flags          *color.Color
	flagsDataType  *color.Color
	flagsDesc      *color.Color
	aliases        *color.Color
	example        *color.Color
	isRequired     *color.Color
}

var (
	lnStartWithFlagRe = regexp.MustCompile(`^\s*(--?\S+)`)
	flagRe            = regexp.MustCompile(`(-\w+,\s|--\S+)`)
	flagDescRe        = regexp.MustCompile(`(\s{3,})\w.*`)
	flagDataTypeRe    = regexp.MustCompile(`(\s\w.*\s{3,})`)
	flagStyleRe       = regexp.MustCompile(`(?i)(\.(InheritedFlags|LocalFlags)\.FlagUsages)`)
	flagIsRequiredRe  = regexp.MustCompile("(required)")
)

func styleHeaders(template string, config *ColorScheme) string {
	if config.headers != nil {
		cobra.AddTemplateFunc("headerStyle", config.headers.SprintFunc())
		template = strings.NewReplacer(
			"Usage:", `{{headerStyle "Usage:"}}`,
			"Aliases:", `{{headerStyle "Aliases:"}}`,
			"Examples:", `{{headerStyle "Examples:"}}`,
			"Available Commands:", `{{headerStyle "Available Commands:"}}`,
			"Global Flags:", `{{headerStyle "Global Flags:"}}`,
			"Additional help topics:", `{{headerStyle "Additional help topics:"}}`,
			"Flags:", `{{headerStyle "Flags:"}}`,
			"{{.Title}}", `{{headerStyle .Title}}`,
			"Available links:", `{{headerStyle "Available Links:"}}`,
			"Additional Commands:", `{{headerStyle "Additional Commands:"}}`,
		).Replace(template)
	}
	return template
}

func styleCommands(template string, config *ColorScheme) string {
	if config.commands != nil {
		cobra.AddTemplateFunc("commandStyle", config.commands.SprintFunc())
		template = strings.ReplaceAll(template, "{{rpad .Name .NamePadding }}", "{{rpad .Name .NamePadding | commandStyle}}")
	}
	return template
}

func styleCmdDesc(template string, config *ColorScheme) string {
	if config.cmdDescription != nil {
		cobra.AddTemplateFunc("cmdDescStyle", config.cmdDescription.SprintFunc())
		template = strings.ReplaceAll(template, `{{.Short}}`, `{{cmdDescStyle .Short}}`)
	}
	return template
}

func styleFlags(template string, config *ColorScheme) string {
	var flagColor, flagDescColor, flagDataTypeColor, flagIsRequiredColor func(a ...any) string
	if config.flags != nil {
		flagColor = config.flags.SprintFunc()
	}
	if config.flagsDesc != nil {
		flagDescColor = config.flagsDesc.SprintFunc()
	}
	if config.flagsDataType != nil {
		flagDataTypeColor = config.flagsDataType.SprintFunc()
	}
	if config.isRequired != nil {
		flagIsRequiredColor = config.isRequired.SprintFunc()
	}

	if flagColor != nil || flagDescColor != nil || flagDataTypeColor != nil {
		fmtShortAndFullFlag := func(s string) string {
			for _, flag := range flagRe.FindAllString(s, -1) {
				s = strings.Replace(s, flag, flagColor(flag), 1)
			}
			return s
		}

		fmtFlagDataType := func(s string) string {
			for _, t := range flagDataTypeRe.FindAllString(s, -1) {
				s = strings.Replace(s, t, flagDataTypeColor(t), 1)
			}
			return s
		}

		fmtFlagDesc := func(s string) string {
			for _, desc := range flagDescRe.FindAllString(s, -1) {
				trimmedDesc := strings.TrimSpace(desc)
				s = strings.Replace(s, trimmedDesc, flagDescColor(trimmedDesc), 1)
			}
			return s
		}

		fmtFlagRequired := func(s string) string {
			for _, desc := range flagIsRequiredRe.FindAllString(s, -1) {
				trimmedDesc := strings.TrimSpace(desc)
				s = strings.Replace(s, trimmedDesc, flagIsRequiredColor(trimmedDesc), -1)
			}
			return s
		}

		flagStyleFunc := func(s string) string {
			lines := strings.Split(s, "\n")
			for i, line := range lines {
				if line == "" {
					continue
				}

				if ok := lnStartWithFlagRe.MatchString(line); ok {
					if flagDescColor != nil {
						line = fmtFlagDesc(line)
					}
					if flagColor != nil {
						line = fmtShortAndFullFlag(line)
					}
					if flagDataTypeColor != nil {
						line = fmtFlagDataType(line)
					}
					if flagIsRequiredColor != nil {
						line = fmtFlagRequired(line)
					}
					lines[i] = line
				} else {
					if flagDescColor != nil {
						line = fmtFlagDesc(line)
					}
					if flagIsRequiredColor != nil {
						line = fmtFlagRequired(line)
					}
					lines[i] = line
				}
			}
			s = strings.Join(lines, "\n")
			return s
		}

		cobra.AddTemplateFunc("flagStyle", flagStyleFunc)
		template = flagStyleRe.ReplaceAllString(template, `flagStyle $1`)
	}
	return template
}

func styleAliases(template string, config *ColorScheme) string {
	if config.aliases != nil {
		cobra.AddTemplateFunc("aliasStyle", config.aliases.SprintFunc())
		template = strings.ReplaceAll(template, `{{.NameAndAliases}}`, `{{aliasStyle .NameAndAliases}}`)
	}
	return template
}

func styleExample(template string, config *ColorScheme) string {
	if config.example != nil {
		cobra.AddTemplateFunc("exampleStyle", config.example.SprintFunc())
		template = strings.ReplaceAll(template, `{{.Example}}`, `{{exampleStyle .Example}}`)
	}
	return template
}

func styleExecName(template string, config *ColorScheme) string {
	if config.execName != nil {
		execColor := config.execName.SprintFunc()
		execNameFunc := func(s string) string {
			spl := strings.Split(s, " ")
			if len(spl) == 0 {
				return s
			}
			spl[0] = execColor(spl[0])
			return strings.Join(spl, " ")
		}
		cobra.AddTemplateFunc("execNameStyle", execNameFunc)
		template = strings.ReplaceAll(template, `{{.CommandPath}}`, `{{execNameStyle .CommandPath}}`)
		template = strings.ReplaceAll(template, `{{.UseLine}}`, `{{execNameStyle .UseLine}}`)
	}
	return template
}

func configureOutputColor(rootCmd *cobra.Command) {
	if rootCmd == nil {
		return
	}

	darkScheme := &ColorScheme{
		headers:       color.New(color.FgHiBlue, color.Bold, color.Underline),
		commands:      color.New(color.FgHiWhite, color.Bold),
		execName:      color.New(color.FgHiWhite, color.Bold),
		flags:         color.New(color.FgHiWhite, color.Bold),
		flagsDataType: color.New(color.FgHiBlack, color.Italic),
		example:       color.New(color.Italic),
		isRequired:    color.New(color.FgHiRed),
	}

	lightScheme := &ColorScheme{
		headers:       color.New(color.FgHiBlue, color.Bold, color.Underline),
		commands:      color.New(color.FgHiBlack, color.Bold),
		execName:      color.New(color.FgHiBlack, color.Bold),
		flags:         color.New(color.FgHiBlack, color.Bold),
		flagsDataType: color.New(color.FgHiBlack, color.Italic),
		example:       color.New(color.Italic),
		isRequired:    color.New(color.FgHiRed),
	}

	colorScheme := lightScheme
	if termenv.HasDarkBackground() {
		colorScheme = darkScheme
	}

	template := rootCmd.UsageTemplate()
	template = styleHeaders(template, colorScheme)
	template = styleCommands(template, colorScheme)
	template = styleCmdDesc(template, colorScheme)
	template = styleExecName(template, colorScheme)
	template = styleFlags(template, colorScheme)
	template = styleAliases(template, colorScheme)
	template = styleExample(template, colorScheme)

	template += "\n"

	rootCmd.SetUsageTemplate(template)
}
