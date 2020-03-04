package openvpn

import "os/exec"

// Run the openvpn
func (c *Config) Run(channel chan *Event) error {
	// the communication channel with openvpn
	mgr, err := newManager()
	if err != nil {
		return err
	}

	mgr.run(channel)
	mgr.updateConfig(c)

	// run openvpn, from now on everything will be managed by the manager
	exec.Command("openvpn", c.params...)
	return nil
}
