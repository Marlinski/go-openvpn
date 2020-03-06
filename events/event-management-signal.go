package events

import (
	"fmt"

	"github.com/Marlinski/go-openvpn/signals"
)

// EventManagementSignal is thrown when the management bootstraps
//
type EventManagementSignal struct {
	InternalEvent
	Sig signals.Signal
}

// Code returns the event code.
func (e EventManagementSignal) Code() EventCode {
	return InternalEventSignal
}

// String returns a human readable version of the event.
func (e EventManagementSignal) String() string {
	return fmt.Sprintf("management has started")
}
