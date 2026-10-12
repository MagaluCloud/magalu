package telemetry

import (
	"strings"
	"time"
)

const (
	UnknownAction  = "unknown"
	UnknownProduct = "unknown"
)

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
	// Authenticated indica login ou API key. Sem nenhum dos dois o comando não é coletado
	Authenticated bool
}

func (c CommandInfo) Action() string {
	if len(c.Path) == 0 {
		return UnknownAction
	}
	segments := make([]string, len(c.Path))
	for i, s := range c.Path {
		segments[i] = strings.ReplaceAll(strings.ToLower(s), "-", "")
	}
	return strings.Join(segments, ".")
}

func (c CommandInfo) Product() string {
	if len(c.Path) == 0 {
		return UnknownProduct
	}
	if c.UnknownCommand && len(c.Path) < 2 {
		return UnknownProduct
	}
	return strings.ReplaceAll(strings.ToLower(c.Path[0]), "-", " ")
}

func (c CommandInfo) ResourceType() string {
	pathLen := len(c.Path)
	if pathLen < 2 {
		return ""
	}
	resourcePath := c.Path[:pathLen-1]
	segments := make([]string, len(resourcePath))
	for i, s := range resourcePath {
		segments[i] = strings.ReplaceAll(strings.ToLower(s), "-", "_")
	}
	return strings.Join(segments, "_")
}

func (c CommandInfo) IsLogin() bool {
	return !c.UnknownCommand && len(c.Path) == 2 && c.Path[0] == "auth" && c.Path[1] == "login"
}

func (c CommandInfo) IsInfrastructure() bool {
	if len(c.Path) == 0 {
		return false
	}
	if c.UnknownCommand && len(c.Path) < 2 {
		return false
	}
	_, nonInfra := nonInfraGroups[c.Path[0]]
	return !nonInfra
}
