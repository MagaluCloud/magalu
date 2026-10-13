package config

const ProjectKey = "project"

func (c *Config) Project() string {
	return c.getString(ProjectKey)
}

func (c *Config) SetProject(id string) error {
	return c.Set(ProjectKey, id)
}

func (c *Config) UnsetProject() error {
	return c.unsetKeys(ProjectKey)
}

func (c *Config) unsetKeys(keys ...string) error {
	for _, key := range keys {
		if c.getString(key) == "" {
			continue
		}
		if err := c.Delete(key); err != nil {
			return err
		}
		if err := c.SetTempConfig(key, ""); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) getString(key string) string {
	var s string
	if err := c.Get(key, &s); err != nil {
		return ""
	}
	return s
}
