package auth

type APIKeyParameters struct {
	Key string
}

func (a APIKeyParameters) GetAPIKey() string {
	return a.Key
}

const APIKeyAuthMethod = "api_key"
