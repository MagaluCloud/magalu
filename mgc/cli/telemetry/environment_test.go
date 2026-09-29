package telemetry

import (
	"fmt"
	"strings"
	"testing"
)

func TestDetectExecutionContext(t *testing.T) {
	testCases := []struct {
		vars        map[string]string
		terminal    bool
		context     ExecutionContext
		environment string
	}{
		{map[string]string{"AI_AGENT": "github-copilot"}, true, ExecutionContextAgent, "github-copilot"},
		{map[string]string{"AI_AGENT": " Claude Code "}, true, ExecutionContextAgent, "claude-code"},
		{map[string]string{"AI_AGENT": "!!!"}, true, ExecutionContextAgent, ""},
		{map[string]string{"AI_AGENT": "!!!", "CLAUDECODE": "1"}, true, ExecutionContextAgent, "claude-code"},
		{map[string]string{"CLAUDECODE": "1"}, true, ExecutionContextAgent, "claude-code"},
		{map[string]string{"GEMINI_CLI": "1"}, true, ExecutionContextAgent, "gemini-cli"},
		{map[string]string{"QWEN_CODE": "1"}, true, ExecutionContextAgent, "qwen-code"},
		{map[string]string{"OPENCODE": "1"}, true, ExecutionContextAgent, "opencode"},
		{map[string]string{"AGENT": "1"}, true, ExecutionContextAgent, "opencode"},
		{map[string]string{"AGENT": "amp"}, true, ExecutionContextAgent, "amp"},
		{map[string]string{"KILO": "1"}, true, ExecutionContextAgent, "kilo"},
		{map[string]string{"CODEX_SANDBOX": "seatbelt"}, true, ExecutionContextAgent, "codex"},
		{map[string]string{"GOOSE_TERMINAL": "1"}, true, ExecutionContextAgent, "goose"},
		{map[string]string{"CLINE_ACTIVE": "true"}, true, ExecutionContextAgent, "cline"},
		{map[string]string{"ROO_ACTIVE": "true"}, true, ExecutionContextAgent, "roo"},
		{map[string]string{"OR_SITE_URL": "https://aider.chat"}, true, ExecutionContextAgent, "aider"},
		{map[string]string{"COPILOT_CLI": "1"}, true, ExecutionContextAgent, "github-copilot"},
		{map[string]string{"CURSOR_AGENT": "1"}, true, ExecutionContextAgent, "cursor"},
		{map[string]string{"ANTIGRAVITY_AGENT": "1"}, true, ExecutionContextAgent, "antigravity"},
		{map[string]string{"AUGMENT_AGENT": "1"}, true, ExecutionContextAgent, "augment"},
		{map[string]string{"CLAUDE_CODE_IS_COWORK": "1", "CLAUDECODE": "1"}, true, ExecutionContextAgent, "cowork"},
		{map[string]string{"JUNIE_DATA": "/tmp/x"}, true, ExecutionContextAgent, "junie"},
		{map[string]string{"TERM_PROGRAM": "kiro"}, true, ExecutionContextAgent, "kiro"},
		{map[string]string{"REPL_ID": "abc"}, true, ExecutionContextAgent, "replit"},

		{map[string]string{"CI": "true"}, true, ExecutionContextCI, ""},
		{map[string]string{"CI": "true", "GITHUB_ACTIONS": "true"}, true, ExecutionContextCI, "github_actions"},
		{map[string]string{"GITLAB_CI": "true"}, true, ExecutionContextCI, "gitlab_ci"},
		{map[string]string{"BUILDKITE": "true"}, true, ExecutionContextCI, "buildkite"},
		{map[string]string{"CIRCLECI": "true"}, true, ExecutionContextCI, "circleci"},
		{map[string]string{"TRAVIS": "true"}, true, ExecutionContextCI, "travis_ci"},
		{map[string]string{"JENKINS_URL": "http://jenkins"}, true, ExecutionContextCI, "jenkins"},
		{map[string]string{"TEAMCITY_VERSION": "2024.1"}, true, ExecutionContextCI, "teamcity"},
		{map[string]string{"TF_BUILD": "True"}, true, ExecutionContextCI, "azure_pipelines"},
		{map[string]string{"BITBUCKET_BUILD_NUMBER": "42"}, true, ExecutionContextCI, "bitbucket_pipelines"},

		{map[string]string{"CLAUDECODE": "1", "CI": "true"}, false, ExecutionContextAgent, "claude-code"},
		{nil, false, ExecutionContextCI, ""},
		{nil, true, ExecutionContextInteractive, ""},
		{map[string]string{"TERM_PROGRAM": "vscode"}, true, ExecutionContextInteractive, ""},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v terminal=%v", tc.vars, tc.terminal), func(t *testing.T) {
			context, environment := DetectExecutionContext(envFrom(tc.vars), tc.terminal)
			if context != tc.context || environment != tc.environment {
				t.Errorf("got (%s, %q), want (%s, %q)", context, environment, tc.context, tc.environment)
			}
		})
	}
}

func TestSanitizeEnvironment(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"free form name with symbols", "Agent/With Spaces;and$symbols-" + string(make([]byte, 100))},
		{"cut right after a separator", strings.Repeat("a", maxEnvironmentLen-1) + " b"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := sanitizeEnvironment(tc.input)

			if len(got) > maxEnvironmentLen {
				t.Errorf("environment too long: %d", len(got))
			}
			if strings.HasSuffix(got, "-") {
				t.Errorf("%q ends with a hyphen", got)
			}
			for _, r := range got {
				if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '.' || r == '_') {
					t.Errorf("unexpected char %q in %q", r, got)
				}
			}
		})
	}
}
