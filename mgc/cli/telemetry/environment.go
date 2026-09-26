package telemetry

import (
	"regexp"
	"strings"
)

type envVarMatch struct {
	envVar      string
	environment string
}

var agentEnvVars = []envVarMatch{
	{"CLAUDE_CODE_IS_COWORK", "cowork"},
	{"CLAUDECODE", "claude-code"},
	{"GEMINI_CLI", "gemini-cli"},
	{"QWEN_CODE", "qwen-code"},
	{"OPENCODE", "opencode"},
	{"KILO", "kilo"},
	{"CODEX_SANDBOX", "codex"},
	{"GOOSE_TERMINAL", "goose"},
	{"CLINE_ACTIVE", "cline"},
	{"ROO_ACTIVE", "roo"},
	{"OR_SITE_URL", "aider"},
	{"COPILOT_CLI", "github-copilot"},
	{"CURSOR_AGENT", "cursor"},
	{"ANTIGRAVITY_AGENT", "antigravity"},
	{"AUGMENT_AGENT", "augment"},
	{"JUNIE_DATA", "junie"},
	{"REPL_ID", "replit"},
}

var ciEnvVars = []envVarMatch{
	{"GITHUB_ACTIONS", "github_actions"},
	{"GITLAB_CI", "gitlab_ci"},
	{"BUILDKITE", "buildkite"},
	{"CIRCLECI", "circleci"},
	{"TRAVIS", "travis_ci"},
	{"JENKINS_URL", "jenkins"},
	{"TEAMCITY_VERSION", "teamcity"},
	{"TF_BUILD", "azure_pipelines"},
	{"BITBUCKET_BUILD_NUMBER", "bitbucket_pipelines"},
	{"CI", ""},
}

var unsafeEnvironmentChars = regexp.MustCompile(`[^a-z0-9._-]+`)

const maxEnvironmentLen = 64

func DetectExecutionContext(getenv func(string) string, isTerminal bool) (ExecutionContext, string) {
	if name, ok := detectAgent(getenv); ok {
		return ExecutionContextAgent, name
	}

	for _, m := range ciEnvVars {
		if getenv(m.envVar) != "" {
			return ExecutionContextCI, m.environment
		}
	}

	if !isTerminal {
		return ExecutionContextCI, ""
	}

	return ExecutionContextInteractive, ""
}

func detectAgent(getenv func(string) string) (string, bool) {
	aiAgent := getenv("AI_AGENT") != ""
	if name := sanitizeEnvironment(getenv("AI_AGENT")); name != "" {
		return name, true
	}

	for _, m := range agentEnvVars {
		if getenv(m.envVar) != "" {
			return m.environment, true
		}
	}

	// todo terminal define TERM_PROGRAM (iTerm, vscode…). Por isso compara o valor.
	if strings.EqualFold(strings.TrimSpace(getenv("TERM_PROGRAM")), "kiro") {
		return "kiro", true
	}

	// AGENT é definida pelo Amp (AGENT=amp) e pelo opencode (valor variável), por isso
	// compara o valor. Fica por último por ser um nome genérico, com risco de falso positivo.
	if v := getenv("AGENT"); v != "" {
		if strings.EqualFold(strings.TrimSpace(v), "amp") {
			return "amp", true
		}
		return "opencode", true
	}

	return "", aiAgent
}

func sanitizeEnvironment(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)
	v = unsafeEnvironmentChars.ReplaceAllString(v, "-")
	if len(v) > maxEnvironmentLen {
		v = v[:maxEnvironmentLen]
	}
	return strings.Trim(v, "-")
}
