package auth

import (
	"context"
	"fmt"
	"net/http"

	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
)

type authRoundTripper struct {
	parent http.RoundTripper
	a      *Auth
}

func newAuthenticatedTransport(parent http.RoundTripper, a *Auth) http.RoundTripper {
	return &authRoundTripper{parent, a}
}

func (o *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	at, err := o.a.AccessToken(req.Context())
	if err != nil {
		return nil, fmt.Errorf("unable to get access token: %w. Did you forget to log-in?", err)
	}
	req.Header.Set("Authorization", "Bearer "+at)
	return o.parent.RoundTrip(req)
}

func (a *Auth) AuthenticatedHttpClientFromContext(ctx context.Context) *mgcHttpPkg.Client {
	unauthenticatedClient := mgcHttpPkg.ClientFromContext(ctx)
	if unauthenticatedClient == nil {
		return nil
	}

	transport := unauthenticatedClient.Transport
	transport = newAuthenticatedTransport(transport, a)

	return mgcHttpPkg.NewClient(transport)
}

func AuthenticatedHttpClientFromContext(ctx context.Context) *mgcHttpPkg.Client {
	a := FromContext(ctx)
	if a == nil {
		return nil
	}

	return a.AuthenticatedHttpClientFromContext(ctx)
}
