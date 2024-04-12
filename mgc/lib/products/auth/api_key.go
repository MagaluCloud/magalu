package auth

type APIKeyParameters struct {
	Key string
}

func (a APIKeyParameters) GetAPIKey() string {
	return a.Key
}

// APIKeyAuthMethod defines authentication using api-key
var APIKeyAuthMethod = "apikeyauth"
