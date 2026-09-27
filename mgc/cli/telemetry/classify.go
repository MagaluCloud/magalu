package telemetry

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"os"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/erikgeiser/promptkit"
)

// TODO: Implementar
type FailureReasoner interface {
	TelemetryFailureReason() FailureReason
}

func isUserCancelled(err error) bool {
	var deniedErr core.UserDeniedConfirmationError

	return errors.As(err, &deniedErr) ||
		errors.Is(err, promptkit.ErrAborted)
}

func isTimeout(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func classifyHttpError(e *mgcHttpPkg.HttpError) FailureReason {
	switch {
	case e.Code == 401:
		return FailureAuthentication
	case e.Code == 403:
		return FailureAuthorization
	case e.Code == 404:
		return FailureNotFound
	case e.Code == 409:
		return FailureConflict
	case e.Code == 429:
		return FailureQuota
	case e.Code >= 400 && e.Code < 500:
		if isQuotaSlug(e.Slug) {
			return FailureQuota
		}
		return FailureValidation
	case e.Code >= 500:
		return FailureAPIServer
	}
	return FailureUnknown
}

func isQuotaSlug(slug string) bool {
	s := strings.NewReplacer("_", "", "-", "", ".", "").Replace(strings.ToLower(slug))
	return strings.Contains(s, "quota") ||
		strings.Contains(s, "limitexceeded") ||
		strings.Contains(s, "limitreached")
}

func isNetworkError(err error) bool {
	var (
		opErr     *net.OpError                      // falha em operação de rede (dial, read, write)
		dnsErr    *net.DNSError                     // falha ao resolver o nome do host
		addrErr   *net.AddrError                    // endereço de rede inválido
		certErr   *tls.CertificateVerificationError // certificado do servidor rejeitado na verificação
		recordErr tls.RecordHeaderError             // resposta não é TLS válido (ex.: HTTP em porta HTTPS)
		authErr   x509.UnknownAuthorityError        // certificado assinado por CA desconhecida
		hostErr   x509.HostnameError                // certificado não corresponde ao hostname
	)
	return errors.As(err, &opErr) ||
		errors.As(err, &dnsErr) ||
		errors.As(err, &addrErr) ||
		errors.As(err, &certErr) ||
		errors.As(err, &recordErr) ||
		errors.As(err, &authErr) ||
		errors.As(err, &hostErr)
}

func ClassifyError(err error) FailureReason {
	if err == nil {
		return ""
	}

	switch {
	case isUserCancelled(err):
		return FailureUserCancelled
	case isTimeout(err):
		return FailureTimeout
	}

	var reasoner FailureReasoner
	if errors.As(err, &reasoner) {
		if r := reasoner.TelemetryFailureReason(); r.Valid() {
			return r
		}
	}

	var httpErr *mgcHttpPkg.HttpError
	if errors.As(err, &httpErr) {
		return classifyHttpError(httpErr)
	}

	if isNetworkError(err) {
		return FailureNetwork
	}

	if errors.As(err, new(core.UsageError)) {
		return FailureValidation
	}

	return FailureUnknown
}
