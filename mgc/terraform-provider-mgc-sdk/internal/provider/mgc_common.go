package provider

type ApiKey struct {
	keyId     string
	keySecret string
}

type Config struct {
	region *string
	env    *string
}
