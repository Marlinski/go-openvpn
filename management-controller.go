package openvpn

import (
	"fmt"

	"github.com/Marlinski/go-openvpn/messages"
)

// Err
var (
	ErrTunnelNotConnected  = fmt.Errorf("the tunnel is not connected")
	ErrEnvironmentNotFound = fmt.Errorf("environment variable not found")
)

// Controller provides API to control the openvpn connection
type Controller struct {
	m *Manager
}

func newController(mgr *Manager) Controller {
	return Controller{
		m: mgr,
	}
}

// GetOpenVpnEnv return the value of a specific environment variable
func (c Controller) GetOpenVpnEnv(vpnEnv string) (string, error) {
	c.m.statemux.Lock()
	defer c.m.statemux.Unlock()

	if c.m.state.state() != MgmtStateCodeConnected {
		return "", ErrTunnelNotConnected
	}

	conn := c.m.state.(*StateMgmtConnected)
	value, ok := conn.env[vpnEnv]
	if !ok {
		return "", ErrEnvironmentNotFound
	}
	return value, nil
}

// SendMgmtCommand send management command
func (c Controller) SendMgmtCommand(cmd string) (messages.Response, error) {
	return c.m.ActionSendCmd(cmd)
}
