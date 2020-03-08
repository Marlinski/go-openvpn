package events

import (
	"fmt"

	"github.com/Marlinski/go-openvpn/messages"
)

// EventBytecount is thrown when a new reading is available on in/out
type EventBytecount struct {
	OpenvpnEvent
	messages.BytecountMessage
}

// Code returns the event code.
func (e EventBytecount) Code() EventCode {
	return OpenvpnEventBytecount
}

// String implements Message String interface
func (e EventBytecount) String() string {
	return fmt.Sprintf("in=%d out=%d", e.In, e.Out)
}
