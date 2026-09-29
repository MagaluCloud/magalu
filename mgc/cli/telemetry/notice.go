package telemetry

import (
	"fmt"
	"io"
)

const PrivacyPolicyURL = "https://magalu.cloud/termos-legais/politica-de-privacidade/"

const Notice = "A MGC CLI coleta dados de uso anônimos/pseudônimos para nos ajudar a\n" +
	"priorizar melhorias. Para desativar: `mgc telemetry disable` ou\n" +
	"defina " + EnvOptOut + "=1. Saiba mais: " + PrivacyPolicyURL

const NoticeNotCollectedYet = "Nenhum dado desta execução foi coletado. " +
	"A coleta começa somente a partir da próxima execução."

// noticePending indica se esta execução deve mostrar o aviso: sessão interativa e
// aviso nunca visto.
func (s *Service) noticePending() bool {
	return s.executionContext == ExecutionContextInteractive && !s.state.NoticeShown
}

// printNotice escreve o aviso. Quando a execução não é coletada, acrescenta
// NoticeNotCollectedYet.
func printNotice(w io.Writer, collected bool) {
	if collected {
		fmt.Fprintf(w, "\n%s\n", Notice)
		return
	}
	fmt.Fprintf(w, "\n%s\n%s\n", Notice, NoticeNotCollectedYet)
}
