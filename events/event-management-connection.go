package events

import "fmt"

// EventManagementConnected is thrown when the vpn has connected
// to the management interface
//
// UnixSocket: the unix socket listening interface
type EventManagementConnected struct {
	InternalEvent
	UnixSocket string
}

// Code returns the event code.
func (e EventManagementConnected) Code() EventCode {
	return InternalEventMgmtConnected
}

// String returns a human readable version of the event.
func (e EventManagementConnected) String() string {
	return fmt.Sprintf("openvpn has connected on %s", e.UnixSocket)
}

// EventManagementDisconnected is thrown when the vpn has connected
// to the management interface
//
// UnixSocket: the unix socket listening interface
type EventManagementDisconnected struct {
	InternalEvent
	UnixSocket string
}

// Code returns the event code.
func (e EventManagementDisconnected) Code() EventCode {
	return InternalEventMgmtDisconnected
}

// String returns a human readable version of the event.
func (e EventManagementDisconnected) String() string {
	return fmt.Sprintf("openvpn has disconnected from %s", e.UnixSocket)
}
