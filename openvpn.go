package openvpn

import (
	"github.com/Marlinski/go-openvpn/events"
)

// Run the openvpn
func (c Config) Run(upstream chan events.OpenvpnEvent) Controller {
	mgr := c.NewManager(upstream)
	mgr.ActionManagerBootstrap()
	return newController(mgr)
}
