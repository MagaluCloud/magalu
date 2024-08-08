package api_key

const (
	scope_PA = "pa:cloud-cli:features"
)

type apiKeysResult struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description,omitempty"`
	StartValidity string  `json:"start_validity"`
	EndValidity   *string `json:"end_validity,omitempty"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	TenantName    *string `json:"tenant_name,omitempty"`
}

type getApiKeyResult struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	ApiKey        string   `json:"api_key"`
	Description   string   `json:"description,omitempty"`
	KeyPairID     string   `json:"key_pair_id"`
	KeyPairSecret string   `json:"key_pair_secret"`
	StartValidity string   `json:"start_validity"`
	EndValidity   *string  `json:"end_validity,omitempty"`
	RevokedAt     *string  `json:"revoked_at,omitempty"`
	TenantName    *string  `json:"tenant_name,omitempty"`
	Scopes        []scopes `json:"scopes"`
}

type scopes struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
}

type apiKeys struct {
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	ApiKey        string  `json:"api_key"`
	Description   string  `json:"description"`
	KeyPairID     string  `json:"key_pair_id"`
	KeyPairSecret string  `json:"key_pair_secret"`
	StartValidity string  `json:"start_validity"`
	EndValidity   *string `json:"end_validity,omitempty"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	TenantName    *string `json:"tenant_name,omitempty"`
	Tenant        struct {
		UUID      string `json:"uuid"`
		LegalName string `json:"legal_name"`
	} `json:"tenant"`
	Scopes []struct {
		UUID        string `json:"uuid"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		ConsentText string `json:"consent_text"`
		Icon        string `json:"icon"`
		APIProduct  struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"api_product"`
	} `json:"scopes"`
}

type scopesCreate struct {
	ID            string `json:"id"`
	RequestReason string `json:"request_reason"`
}

type createApiKey struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	TenantID      string         `json:"tenant_id"`
	ScopesList    []scopesCreate `json:"scopes"`
	StartValidity string         `json:"start_validity"`
	EndValidity   string         `json:"end_validity"`
}

type apiKeyResult struct {
	UUID string `json:"uuid,omitempty"`
	Used bool   `json:"used,omitempty"`
}

func (r *apiKeys) ToResult() *apiKeysResult {
	return &apiKeysResult{
		ID:            r.UUID,
		Name:          r.Name,
		Description:   r.Description,
		StartValidity: r.StartValidity,
		EndValidity:   r.EndValidity,
		RevokedAt:     r.RevokedAt,
		TenantName:    r.TenantName,
	}
}

func (r *apiKeys) ToResultGet() *getApiKeyResult {
	var scopesv []scopes
	for _, s := range r.Scopes {
		scopesv = append(scopesv, scopes{
			ID:    s.UUID,
			Name:  s.Name,
			Title: s.Title,
		})
	}

	return &getApiKeyResult{
		ID:            r.UUID,
		Name:          r.Name,
		ApiKey:        r.ApiKey,
		Description:   r.Description,
		KeyPairID:     r.KeyPairID,
		KeyPairSecret: r.KeyPairSecret,
		StartValidity: r.StartValidity,
		EndValidity:   r.EndValidity,
		RevokedAt:     r.RevokedAt,
		TenantName:    r.TenantName,
		Scopes:        scopesv,
	}
}
