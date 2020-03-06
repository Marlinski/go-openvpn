package openvpn

import (
	"github.com/op/go-logging"
)

// Config is the openvpn configuration
// in the form of parameters to the command line
type Config struct {
	// the ID for this openvpn config
	// it will be used as a token in log
	id string

	// if set to true, will log the stdout and stderr of openvpn
	logLevel logging.Level

	// logStd
	logStd bool

	// params to give to the openvpn command
	params []string
}

// NewConfig return an empty config
func NewConfig(id string) *Config {
	return &Config{
		id:     id,
		params: make([]string, 0),
	}
}

// LoadConfig sets the openvpn config
func LoadConfig(id string, file string) *Config {
	return &Config{
		id: id,
		params: []string{
			"--config",
			file,
		},
	}
}

// SetLogLevel sets the log level for this configuration
func (c *Config) SetLogLevel(level logging.Level) *Config {
	c.logLevel = level
	return c
}

// SetLogStd sets wether we log standard output or not
func (c *Config) SetLogStd(enable bool) *Config {
	c.logStd = enable
	return c
}

// Flag parameter
func (c *Config) Flag(flag string) *Config {
	c.params = append(c.params, "--"+flag)
	return c
}

// Set a parameter
func (c *Config) Set(key string, values ...string) *Config {
	c.params = append(c.params, "--"+key)
	for _, v := range values {
		c.params = append(c.params, v)
	}
	return c
}
