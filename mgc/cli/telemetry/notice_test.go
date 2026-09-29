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

func TestNoticeFirstInteractiveLogin(t *testing.T) {
	svc, store, _ := newTestService(testSetup{terminal: true})

	stderr := record(svc, commandAt(authLogin, loginAt), nil)

	if stderr != expectedNotice {
		t.Errorf("got %q\nwant %q", stderr, expectedNotice)
	}
	if !store.state.NoticeShown {
		t.Error("notice_shown must be persisted")
	}
	if store.state.CredentialsSetAt == nil || !store.state.CredentialsSetAt.Equal(loginAt) {
		t.Errorf("credentials_set_at must be recorded, got %v", store.state.CredentialsSetAt)
	}
}

func TestNoticeAlreadyLoggedIn(t *testing.T) {
	svc, store, _ := newTestService(testSetup{terminal: true})

	stderr := record(svc, commandAt(vmList, loginAt), nil)

	if stderr != expectedNoticeNotCollected {
		t.Errorf("got %q\nwant %q", stderr, expectedNoticeNotCollected)
	}
	if !store.state.NoticeShown {
		t.Error("notice_shown must be persisted")
	}
}

func TestNoticeShownOnlyOnce(t *testing.T) {
	firstLogin := loginAt.Add(-24 * time.Hour)
	seen := State{NoticeShown: true, CredentialsSetAt: &firstLogin}

	for name, path := range map[string][]string{"second login": authLogin, "next command": vmList} {
		svc, store, _ := newTestService(testSetup{state: seen, terminal: true})

		if stderr := record(svc, commandAt(path, loginAt), nil); stderr != "" {
			t.Errorf("%s: notice must be shown only once, got %q", name, stderr)
		}
		if !store.state.CredentialsSetAt.Equal(firstLogin) {
			t.Errorf("%s: credentials_set_at must keep the first login time", name)
		}
	}
}

func TestNoticeNotShownOutsideInteractive(t *testing.T) {
	for name, env := range map[string]map[string]string{
		"ci":    {"CI": "true"},
		"agent": {"CLAUDECODE": "1"},
	} {
		svc, store, _ := newTestService(testSetup{env: env, terminal: true})

		stderr := record(svc, commandAt(authLogin, loginAt), nil)
		stderr += record(svc, commandAt(vmList, loginAt), nil)

		if stderr != "" || store.state.NoticeShown {
			t.Errorf("%s: notice must not be shown, got %q", name, stderr)
		}
		if store.state.CredentialsSetAt == nil {
			t.Errorf("%s: credentials_set_at is still recorded after login", name)
		}
	}
}
