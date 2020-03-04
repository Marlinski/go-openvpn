package openvpn

// Config is the openvpn configuration
// in the form of parameters to the command line
type Config struct {
	params []string
}

// NewConfig return an empty config
func NewConfig() *Config {
	return &Config{
		params: make([]string, 0),
	}
}

// LoadConfig sets the openvpn config
func LoadConfig(file string) *Config {
	return &Config{
		params: []string{
			"--config",
			file,
		},
	}
}

// Flag parameter
func (c *Config) Flag(flag string) *Config {
	c.params = append(c.params, flag)
	return c
}

// Set a parameter
func (c *Config) Set(key string, values ...string) *Config {
	c.params = append(c.params, key)
	for _, v := range values {
		c.params = append(c.params, v)
	}
	return c
}
