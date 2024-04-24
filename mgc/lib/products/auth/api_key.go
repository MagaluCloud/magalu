package auth

type APIKeyParameters struct {
	Key string
}

func (a APIKeyParameters) GetAPIKey() string {
	return a.Key
}
