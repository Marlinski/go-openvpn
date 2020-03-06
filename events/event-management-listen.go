package events

import "fmt"

// EventManagementListenSocket is thrown when management has set up the unix socket
// and is waiting for openvpn to connect.
//
// UnixSocket: the unix socket listening interface
type EventManagementListenSocket struct {
	InternalEvent
	UnixSocket string
}

// Code returns the event code.
func (e EventManagementListenSocket) Code() EventCode {
	return InternalEventListenSocket
}

// String returns a human readable version of the event.
func (e EventManagementListenSocket) String() string {
	return fmt.Sprintf("waiting for openvpvn to connect on %s", e.UnixSocket)
}

// EventManagementListenSocketError is thrown when an error happened while listening
// the socket
//
// UnixSocket: the unix socket listening interface
type EventManagementListenSocketError struct {
	InternalEvent
	Err error
}

// Code returns the event code.
func (e EventManagementListenSocketError) Code() EventCode {
	return InternalEventListenSocketError
}

// String returns a human readable version of the event.
func (e EventManagementListenSocketError) String() string {
	return fmt.Sprintf("%+v", e.Err)
}
