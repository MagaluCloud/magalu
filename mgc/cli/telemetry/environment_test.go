package telemetry

import (
	"strings"
	"testing"
)

func TestEnvironmentAgents(t *testing.T) {
	cases := []struct {
		vars map[string]string
		want string
	}{
		{map[string]string{"AI_AGENT": "github-copilot"}, "github-copilot"},
		{map[string]string{"AI_AGENT": " Claude Code "}, "claude-code"},
		{map[string]string{"CLAUDECODE": "1"}, "claude-code"},
		{map[string]string{"GEMINI_CLI": "1"}, "gemini-cli"},
		{map[string]string{"QWEN_CODE": "1"}, "qwen-code"},
		{map[string]string{"OPENCODE": "1"}, "opencode"},
		{map[string]string{"AGENT": "1"}, "opencode"},
		{map[string]string{"AGENT": "amp"}, "amp"},
		{map[string]string{"KILO": "1"}, "kilo"},
		{map[string]string{"CODEX_SANDBOX": "seatbelt"}, "codex"},
		{map[string]string{"GOOSE_TERMINAL": "1"}, "goose"},
		{map[string]string{"CLINE_ACTIVE": "true"}, "cline"},
		{map[string]string{"ROO_ACTIVE": "true"}, "roo"},
		{map[string]string{"OR_SITE_URL": "https://aider.chat"}, "aider"},
		{map[string]string{"COPILOT_CLI": "1"}, "github-copilot"},
		{map[string]string{"CURSOR_AGENT": "1"}, "cursor"},
		{map[string]string{"ANTIGRAVITY_AGENT": "1"}, "antigravity"},
		{map[string]string{"AUGMENT_AGENT": "1"}, "augment"},
		{map[string]string{"CLAUDE_CODE_IS_COWORK": "1", "CLAUDECODE": "1"}, "cowork"},
		{map[string]string{"JUNIE_DATA": "/tmp/x"}, "junie"},
		{map[string]string{"TERM_PROGRAM": "kiro"}, "kiro"},
		{map[string]string{"REPL_ID": "abc"}, "replit"},
	}

	for _, c := range cases {
		ctx, env := DetectExecutionContext(envFrom(c.vars), true)
		if ctx != ExecutionContextAgent || env != c.want {
			t.Errorf("%v: got (%s, %q), want (agent, %q)", c.vars, ctx, env, c.want)
		}
	}
}

func TestEnvironmentCI(t *testing.T) {
	cases := []struct {
		vars map[string]string
		want string
	}{
		{map[string]string{"CI": "true"}, ""},
		{map[string]string{"CI": "true", "GITHUB_ACTIONS": "true"}, "github_actions"},
		{map[string]string{"GITLAB_CI": "true"}, "gitlab_ci"},
		{map[string]string{"BUILDKITE": "true"}, "buildkite"},
		{map[string]string{"CIRCLECI": "true"}, "circleci"},
		{map[string]string{"TRAVIS": "true"}, "travis_ci"},
		{map[string]string{"JENKINS_URL": "http://jenkins"}, "jenkins"},
		{map[string]string{"TEAMCITY_VERSION": "2024.1"}, "teamcity"},
		{map[string]string{"TF_BUILD": "True"}, "azure_pipelines"},
		{map[string]string{"BITBUCKET_BUILD_NUMBER": "42"}, "bitbucket_pipelines"},
	}

	for _, c := range cases {
		ctx, env := DetectExecutionContext(envFrom(c.vars), true)
		if ctx != ExecutionContextCI || env != c.want {
			t.Errorf("%v: got (%s, %q), want (ci, %q)", c.vars, ctx, env, c.want)
		}
	}
}

func TestEnvironmentPrecedence(t *testing.T) {
	ctx, env := DetectExecutionContext(envFrom(map[string]string{"CLAUDECODE": "1", "CI": "true"}), false)
	if ctx != ExecutionContextAgent || env != "claude-code" {
		t.Errorf("agent inside CI: got (%s, %q)", ctx, env)
	}

	ctx, env = DetectExecutionContext(envFrom(nil), false)
	if ctx != ExecutionContextCI || env != "" {
		t.Errorf("script without TTY: got (%s, %q)", ctx, env)
	}

	ctx, env = DetectExecutionContext(envFrom(nil), true)
	if ctx != ExecutionContextInteractive || env != "" {
		t.Errorf("human terminal: got (%s, %q)", ctx, env)
	}

	ctx, _ = DetectExecutionContext(envFrom(map[string]string{"TERM_PROGRAM": "vscode"}), true)
	if ctx != ExecutionContextInteractive {
		t.Errorf("TERM_PROGRAM other than kiro must not mark an agent, got %s", ctx)
	}
}

func TestEnvironmentSanitizesFreeFormAgentName(t *testing.T) {
	long := "Agent/With Spaces;and$symbols-" + string(make([]byte, 100))
	_, env := DetectExecutionContext(envFrom(map[string]string{"AI_AGENT": long}), true)
	if len(env) > maxEnvironmentLen {
		t.Errorf("environment too long: %d", len(env))
	}
	for _, r := range env {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '.' || r == '_') {
			t.Fatalf("unexpected char %q in %q", r, env)
		}
	}
}

func TestEnvironmentAIAgentWithOnlySymbols(t *testing.T) {
	ctx, env := DetectExecutionContext(envFrom(map[string]string{"AI_AGENT": "!!!"}), true)
	if ctx != ExecutionContextAgent || env != "" {
		t.Errorf("got (%q, %q), want (agent, \"\")", ctx, env)
	}

	ctx, env = DetectExecutionContext(envFrom(map[string]string{"AI_AGENT": "!!!", "CLAUDECODE": "1"}), true)
	if ctx != ExecutionContextAgent || env != "claude-code" {
		t.Errorf("got (%q, %q), want (agent, claude-code)", ctx, env)
	}
}

func TestEnvironmentSanitizeDoesNotEndWithHyphenAfterCut(t *testing.T) {
	v := strings.Repeat("a", maxEnvironmentLen-1) + " b"
	got := sanitizeEnvironment(v)
	if strings.HasSuffix(got, "-") {
		t.Errorf("sanitizeEnvironment(%q) = %q, ends with hyphen", v, got)
	}
}
