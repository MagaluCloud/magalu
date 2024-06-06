package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type genericIDNameModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// BEGIN CONFIG
type (
	Config interface {
		Region() *string
		Env() *string
		SetRegion(region string)
		SetEnv(env string)
	}

	config struct {
		region *string
		env    *string
	}
)

func (c *config) Region() *string {
	return c.region
}

func (c *config) Env() *string {
	return c.env
}

func (c *config) SetRegion(region string) {
	c.region = &region
}

func (c *config) SetEnv(env string) {
	c.env = &env
}

func NewConfigRaw() Config {
	return &config{}
}

func NewConfig(region, env *string) Config {
	return &config{
		region: region,
		env:    env,
	}
}

func NewConfigDefault(region, env *string) Config {
	if region == nil || *region == "" {
		region = new(string)
		*region = "br-se1"
	}
	if env == nil || *env == "" {
		env = new(string)
		*env = "prod"
	}

	return &config{
		region: region,
		env:    env,
	}
}

// END CONFIG
