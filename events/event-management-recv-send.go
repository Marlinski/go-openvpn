package events

import (
	"fmt"

	"github.com/Marlinski/go-openvpn/messages"
)

// EventManagementRecvMsg is thrown when a message is received from
// the openvpn management interface
//
// UnixSocket: the unix socket listening interface
type EventManagementRecvMsg struct {
	InternalEvent
	Msg messages.Message
}

// Code returns the event code.
func (e EventManagementRecvMsg) Code() EventCode {
	return InternalEventRecvMsg
}

// String returns a human readable version of the event.
func (e EventManagementRecvMsg) String() string {
	return fmt.Sprintf("recv msg: %s", e.Msg.String())
}

// EventManagementReadError is thrown when a message is received from
// the openvpn management interface
//
// UnixSocket: the unix socket listening interface
type EventManagementReadError struct {
	InternalEvent
	Err error
}

// Code returns the event code.
func (e EventManagementReadError) Code() EventCode {
	return InternalEventReadError
}

// String returns a human readable version of the event.
func (e EventManagementReadError) String() string {
	return fmt.Sprintf("read error: %+v", e.Err)
}
