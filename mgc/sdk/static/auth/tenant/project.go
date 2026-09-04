package tenant

import (
	"context"
	"fmt"

	mgcAuthPkg "github.com/MagaluCloud/magalu/mgc/core/auth"
	mgcConfigPkg "github.com/MagaluCloud/magalu/mgc/core/config"
)

// SelectDefaultProject aponta a configuração de projeto para o projeto default
// do tenant atual, que no IAM tem o mesmo id do tenant.
//
// Roda depois de toda troca de tenant — `auth login` e `auth tenant set`. Um
// projeto pertence a um tenant, então mantê-lo depois da troca deixaria a CLI
// mandando o `x-project-id` de um tenant junto do token de outro.
//
// O id vem da claim do token novo, não do que foi pedido: quem decide o tenant
// é o IDP, e é o token que autentica a request seguinte.
func SelectDefaultProject(ctx context.Context, auth *mgcAuthPkg.Auth) error {
	config := mgcConfigPkg.FromContext(ctx)
	if config == nil {
		return fmt.Errorf("programming error: unable to retrieve config from context")
	}

	tenantID, err := auth.CurrentTenantID()
	if err != nil {
		return err
	}

	if tenantID == "" {
		// Token sem claim de tenant não deveria acontecer com o IDP da MGC. Se
		// acontecer, ficar sem projeto é o menos errado: manter o anterior é a
		// própria falha que esta função existe para evitar.
		return config.UnsetProject()
	}

	return config.SetProject(tenantID)
}
