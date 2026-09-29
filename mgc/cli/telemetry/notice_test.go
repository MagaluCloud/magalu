package telemetry

import (
	"testing"
	"time"
)

const expectedNotice = "\nA MGC CLI coleta dados de uso anônimos/pseudônimos para nos ajudar a\n" +
	"priorizar melhorias. Para desativar: `mgc telemetry disable` ou\n" +
	"defina MGC_CLI_TELEMETRY_OPTOUT=1. Saiba mais: https://magalu.cloud/termos-legais/politica-de-privacidade/\n"

const expectedNoticeNotCollected = expectedNotice +
	"Nenhum dado desta execução foi coletado. A coleta começa somente a partir da próxima execução.\n"

func TestNotice(t *testing.T) {
	firstLogin := loginAt.Add(-24 * time.Hour)
	seen := State{NoticeShown: true, CredentialsSetAt: &firstLogin}

	testCases := []struct {
		name       string
		setup      testSetup
		executions []execution
		check      func(t *testing.T, r runResult)
	}{
		{
			name:       "first interactive login",
			setup:      testSetup{terminal: true},
			executions: []execution{{path: authLogin, stderr: expectedNotice, events: 1}},
			check: func(t *testing.T, r runResult) {
				if !r.state.NoticeShown {
					t.Error("notice_shown must be persisted")
				}
				if r.state.CredentialsSetAt == nil || !r.state.CredentialsSetAt.Equal(loginAt) {
					t.Errorf("credentials_set_at must be recorded, got %v", r.state.CredentialsSetAt)
				}
			},
		},
		{
			name:       "already logged in user",
			setup:      testSetup{terminal: true},
			executions: []execution{{path: vmList, stderr: expectedNoticeNotCollected}},
			check: func(t *testing.T, r runResult) {
				if !r.state.NoticeShown {
					t.Error("notice_shown must be persisted")
				}
			},
		},
		{
			name:       "second login",
			setup:      testSetup{state: seen, terminal: true},
			executions: []execution{{path: authLogin, events: 1}},
			check: func(t *testing.T, r runResult) {
				if !r.state.CredentialsSetAt.Equal(firstLogin) {
					t.Error("credentials_set_at must keep the first login time")
				}
			},
		},
		{
			name:       "next command after the notice",
			setup:      testSetup{state: seen, terminal: true},
			executions: []execution{{path: vmList, events: 1}},
		},
		{
			name:  "ci",
			setup: testSetup{env: map[string]string{"CI": "true"}, terminal: true},
			executions: []execution{
				{path: authLogin, events: 1},
				{path: vmList, events: 1},
			},
			check: func(t *testing.T, r runResult) {
				if r.state.NoticeShown || r.state.CredentialsSetAt == nil {
					t.Errorf("want no notice but credentials_set_at recorded, got %+v", r.state)
				}
			},
		},
		{
			name:  "agent",
			setup: testSetup{env: map[string]string{"CLAUDECODE": "1"}, terminal: true},
			executions: []execution{
				{path: authLogin, events: 1},
				{path: vmList, events: 1},
			},
			check: func(t *testing.T, r runResult) {
				if r.state.NoticeShown || r.state.CredentialsSetAt == nil {
					t.Errorf("want no notice but credentials_set_at recorded, got %+v", r.state)
				}
			},
		},
		{
			name:       "telemetry disabled",
			setup:      testSetup{env: map[string]string{"DO_NOT_TRACK": "1"}, terminal: true},
			executions: []execution{{path: authLogin}},
			check: func(t *testing.T, r runResult) {
				if r.state.NoticeShown || r.saves != 0 {
					t.Errorf("disabled telemetry must neither show nor persist the notice, got %+v saves=%d", r.state, r.saves)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := runExecutions(t, tc.setup, tc.executions)
			if tc.check != nil {
				tc.check(t, r)
			}
		})
	}
}
