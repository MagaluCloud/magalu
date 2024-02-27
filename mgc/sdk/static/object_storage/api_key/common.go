package api_key

const (
	name_ObjectStorage = "Object Storage"
	scope_PA           = "pa:cloud-cli:features"
)

type apiKeysResult struct {
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	KeyPairID     string  `json:"key_pair_id"`
	KeyPairSecret string  `json:"key_pair_secret"`
	StartValidity string  `json:"start_validity"`
	EndValidity   *string `json:"end_validity,omitempty"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	TenantName    *string `json:"tenant_name,omitempty"`
}
type apiKeys struct {
	apiKeysResult
	Tenant struct {
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
			Icon string `json:"icon"`
		} `json:"api_product"`
	} `json:"scopes"`
}

type scopesObjectStorage struct {
	ID            string `json:"id"`
	RequestReason string `json:"request_reason"`
}

type createApiKey struct {
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	TenantID      string                `json:"tenant_id"`
	ScopesList    []scopesObjectStorage `json:"scopes"`
	StartValidity string                `json:"start_validity"`
	EndValidity   string                `json:"end_validity"`
}

type apiKeyResult struct {
	UUID string `json:"uuid,omitempty"`
}

type authSetParams struct {
	AccessKeyId     string `json:"access_key_id" jsonschema_description:"Access key id value" mgc:"positional"`
	SecretAccessKey string `json:"secret_access_key" jsonschema_description:"Secret access key value" mgc:"positional"`
}
