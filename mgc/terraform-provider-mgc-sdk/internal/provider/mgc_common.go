package provider

type Config struct {
	region *string
	env    *string
}

func DefaultRegion(region string) string {
	return "br-se1"
}

func DefaultEnv(env string) string {
	return "prod"
}
