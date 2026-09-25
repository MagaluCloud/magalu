package telemetry

import (
	"strings"
	"time"
)

const UnknownAction = "unknown"

var nonInfraGroups = map[string]struct{}{
	"auth":      {},
	"config":    {},
	"profile":   {},
	"workspace": {},
	"telemetry": {},
}

type CommandInfo struct {
	Path           []string
	UnknownCommand bool
	OptionsSet     []string
	ResourceID     string
	Start          time.Time
	End            time.Time
	TenantID       string
	LastRequestID  string
}

func (c CommandInfo) Action() string {
	if c.UnknownCommand || len(c.Path) == 0 {
		return UnknownAction
	}
	segments := make([]string, len(c.Path))
	for i, s := range c.Path {
		segments[i] = strings.ReplaceAll(strings.ToLower(s), "-", "")
	}
	return strings.Join(segments, ".")
}

func (c CommandInfo) ResourceType() string {
	if c.UnknownCommand || len(c.Path) < 2 {
		return ""
	}
	segments := make([]string, len(c.Path)-1)
	for i, s := range c.Path[:len(c.Path)-1] {
		segments[i] = strings.ReplaceAll(strings.ToLower(s), "-", "_")
	}
	return strings.Join(segments, "_")
}

func (c CommandInfo) IsLogin() bool {
	return !c.UnknownCommand && len(c.Path) == 2 && c.Path[0] == "auth" && c.Path[1] == "login"
}

func (c CommandInfo) IsInfrastructure() bool {
	if c.UnknownCommand || len(c.Path) == 0 {
		return false
	}
	_, nonInfra := nonInfraGroups[c.Path[0]]
	return !nonInfra
}
