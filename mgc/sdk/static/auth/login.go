package auth

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/browser"
	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/core/auth"
	mgcAuthPkg "magalu.cloud/core/auth"
)

type authResult struct {
	value string
	err   error
}

type loginParameters struct {
	Show bool `json:"show,omitempty" jsonschema_description:"Show the access token after the login completes"`
}

type loginResult struct {
	AccessToken string `json:"access_token,omitempty"`
}

const serverShutdownTimeout = 500 * time.Millisecond

var (
	//go:embed success.html
	successPage         string
	loginLoggerInstance *zap.SugaredLogger
)

func loginLogger() *zap.SugaredLogger {
	if loginLoggerInstance == nil {
		loginLoggerInstance = initPkgLogger().Named("login")
	}
	return loginLoggerInstance
}

func newLogin() *core.StaticExecute {
	return core.NewStaticExecute(
		"login",
		"",
		"authenticate with magalu cloud",
		func(ctx context.Context, parameters loginParameters, _ struct{}) (output *loginResult, err error) {
			auth := mgcAuthPkg.FromContext(ctx)
			if auth == nil {
				return nil, fmt.Errorf("unable to retrieve authentication configuration")
			}

			resultChan, cancel, err := startCallbackServer(ctx, auth)
			if err != nil {
				return nil, err
			}
			defer cancel()

			codeUrl, err := auth.CodeChallengeToURL()
			if err != nil {
				return nil, err
			}

			loginLogger().Infow("opening browser", "codeUrl", codeUrl)
			if err := browser.OpenURL(codeUrl.String()); err != nil {
				loginLogger().Warnw("could not open browser, please open it manually", "error", err)
			}

			loginLogger().Infow("waiting authentication result", "redirectUri", auth.RedirectUri())
			result := <-resultChan
			if result.err != nil {
				return nil, result.err
			}

			tenants, err := auth.ListTenants(ctx)
			if err != nil || len(tenants) == 0 {
				return nil, fmt.Errorf("error when trying to list tenants for selection: %w", err)
			}

			tenantResult, err := auth.SelectTenant(ctx, tenants[0].UUID)
			if err != nil {
				return nil, fmt.Errorf("error when trying to select default tenant: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Successfully logged in.\n")
			loginLogger().Infow("sucessfully logged in")
			if parameters.Show {
				output = &loginResult{AccessToken: tenantResult.AccessToken}
			}

			return output, nil
		},
	)
}

func startCallbackServer(ctx context.Context, auth *auth.Auth) (resultChan chan *authResult, cancel func(), err error) {
	callbackUrl, err := url.Parse(auth.RedirectUri())
	if err != nil {
		return nil, nil, fmt.Errorf("invalid redirect_uri configuration")
	}

	// Host includes the port, then listen to specific address + port, ex: "localhost:8095"
	addr := callbackUrl.Host

	// Listen so we can fail early on bad address, before starting goroutine
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	resultChan = make(chan *authResult, 1)
	callbackChan := make(chan *authResult, 1)
	cancelChan := make(chan struct{}, 1)

	handler := &callbackHandler{
		auth,
		callbackUrl.Path,
		callbackChan,
		ctx,
	}
	srv := &http.Server{Addr: addr, Handler: handler}

	// serve HTTP until shutdown happened, then report result via channel
	serveAndReportResult := func() {
		serverErrorChan := make(chan error, 1)
		signalChan := make(chan os.Signal, 1)
		serverDoneChan := make(chan *authResult, 1)

		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		waitChannelsAndShutdownServer := func() {
			var result *authResult

			select {
			case err := <-serverErrorChan:
				result = &authResult{err: fmt.Errorf("Could not serve HTTP: %w", err)}

			case sig := <-signalChan:
				result = &authResult{err: fmt.Errorf("Canceled by signal: %v", sig)}

			case <-cancelChan:
				result = &authResult{err: fmt.Errorf("Canceled by user")}

			case result = <-callbackChan:
			}

			signal.Stop(signalChan)

			ctx, cancelShutdown := context.WithTimeout(context.Background(), serverShutdownTimeout)
			defer cancelShutdown()

			// sync: unblocks serveAndReportResult()/srv.Serve()
			if err := srv.Shutdown(ctx); err != nil {
				srv.Close() // aggressively try to close it
			}

			// sync: notify serveAndReportResult() we're done
			serverDoneChan <- result
		}
		go waitChannelsAndShutdownServer()

		if err := srv.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			// sync: unblock waitChannelsAndShutdownServer()
			serverErrorChan <- err
		}

		result := <-serverDoneChan // sync: wait server shutdown by waitChannelsAndShutdownServer()

		close(callbackChan)
		close(cancelChan)

		close(serverErrorChan)
		close(signalChan)
		close(serverDoneChan)

		resultChan <- result
	}

	cancel = func() {
		defer func() {
			// serveAndReportResult() will close channels when done.
			// That means there is nothing to cancel and we should do nothing else, just ignore.
			_ = recover()
		}()

		cancelChan <- struct{}{} // exit waitChannelsAndShutdownServer()
		<-resultChan             // wait serveAndReportResult(), discard as results are not meaningful
	}

	go serveAndReportResult()

	return resultChan, cancel, nil
}

type callbackHandler struct {
	auth *mgcAuthPkg.Auth
	path string
	done chan *authResult
	ctx  context.Context
}

func (h *callbackHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL
	if url.Path != h.path {
		err := fmt.Errorf("Unknown Path: %s", url)
		showErrorPage(w, err, http.StatusNotFound)
		return
	}

	authCode := url.Query().Get("code")
	err := h.auth.RequestAuthTokenWithAuthorizationCode(h.ctx, authCode)
	if err != nil {
		showErrorPage(w, err, http.StatusInternalServerError)
		h.done <- &authResult{err: fmt.Errorf("Could not request auth token with authorization code: %w", err)}
		return
	}

	if err := showSuccessPage(w); err != nil {
		loginLogger().Warnw("could not show whole Succes Page", "error", err)
	}

	token, _ := h.auth.AccessToken(h.ctx) // this is guaranteed if RequestAuthTokeWithAuthorizationCode succeeds
	h.done <- &authResult{value: token}
}

func showSuccessPage(w http.ResponseWriter) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html")
	if _, err := io.WriteString(w, successPage); err != nil {
		return err
	}

	return nil
}

func showErrorPage(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Error: %s", err.Error())
}
