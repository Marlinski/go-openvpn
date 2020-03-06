package events

import (
	"fmt"

	"github.com/Marlinski/go-openvpn/messages"
)

// EventByteCount is thrown when a new reading is available on in/out
type EventByteCount struct {
	OpenvpnEvent
	messages.ByteCountMessage
}

// Code returns the event code.
func (e EventByteCount) Code() EventCode {
	return OpenvpnEventByteCode
}

// String implements Message String interface
func (e EventByteCount) String() string {
	return fmt.Sprintf("in=%d out=%d", e.In, e.Out)
}
