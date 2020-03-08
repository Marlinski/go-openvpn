package events

import "fmt"

// EventTunnelUp is thrown when the openvpn tunnel is up
type EventTunnelUp struct {
	OpenvpnEvent
	Env map[string]string // environment variable
}

// Code returns the event code.
func (e EventTunnelUp) Code() EventCode {
	return OpenvpnEventUp
}

// String returns a human readable version of the event.
func (e EventTunnelUp) String() string {
	return fmt.Sprintf("openvpn tunnel is up")

}

// EventTunnelDown is thrown when the openvpn tunnel status is down
type EventTunnelDown struct {
	OpenvpnEvent
	Env map[string]string // environment variable
}

// Code returns the event code.
func (e EventTunnelDown) Code() EventCode {
	return OpenvpnEventDown
}

// String returns a human readable version of the event.
func (e EventTunnelDown) String() string {
	return fmt.Sprintf("openvpn tunnel is down")
}
