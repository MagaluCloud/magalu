package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type accessTokenParameters struct {
	Validate bool `json:",omitempty" jsonschema_description:"Validate the token, refreshing if needed"`
}

// accessTokenResult is the output of access_token. It holds AccessToken, or, when kubectl runs the
// command as an exec credential plugin, the fields of a client.authentication.k8s.io ExecCredential.
// The CLI prints an ExecCredential as JSON whatever the output format is
// (cli/cmd: RunByKubectl).
type accessTokenResult struct {
	AccessToken string                `json:"access_token,omitempty"`
	APIVersion  string                `json:"apiVersion,omitempty"`
	Kind        string                `json:"kind,omitempty"`
	Status      *execCredentialStatus `json:"status,omitempty"`
}

// execCredentialStatus is an alias of an unnamed struct, so the schema reflector puts its fields
// inline. The core schema conversion gives an empty schema for a $ref to a named struct, and
// then the result validation fails.
type execCredentialStatus = struct {
	ExpirationTimestamp string `json:"expirationTimestamp"`
	Token               string `json:"token"`
}

// RunByKubectl reports whether kubectl runs the command as an exec credential plugin.
func RunByKubectl() bool {
	_, ok := os.LookupEnv("KUBERNETES_EXEC_INFO")
	return ok
}

var getAccessToken = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "access_token",
			Description: "Retrieve the access token used in the APIs",
		},
		func(ctx context.Context, parameters accessTokenParameters, _ struct{}) (output *accessTokenResult, err error) {
			auth := mgcAuthPkg.FromContext(ctx)
			if auth == nil {
				return nil, fmt.Errorf("unable to retrieve authentication configuration")
			}

			if parameters.Validate {
				err := auth.ValidateAccessToken(ctx)
				if err != nil {
					return nil, fmt.Errorf("Could not validate the Access Token: %w", err)
				}
			}

			token, err := auth.AccessToken(ctx)
			if err != nil {
				return nil, err
			}

			if !RunByKubectl() {
				return &accessTokenResult{AccessToken: token}, nil
			}

			expiresAt, err := auth.AccessTokenExpiresAt()
			if err != nil {
				return nil, fmt.Errorf("could not read the expiry of the Access Token: %w", err)
			}
			return &accessTokenResult{
				APIVersion: "client.authentication.k8s.io/v1",
				Kind:       "ExecCredential",
				Status: &execCredentialStatus{
					ExpirationTimestamp: expiresAt.UTC().Format(time.RFC3339),
					Token:               token,
				},
			}, nil
		},
	)
})
