package openvpn

import (
	"fmt"

	"github.com/Marlinski/go-openvpn/events"
	"github.com/Marlinski/go-openvpn/messages"
	"github.com/Marlinski/go-openvpn/signals"
)

// MgmtStateCode are the different machine states
type MgmtStateCode string

// al the states
const (
	MgmtStateCodeIdle         MgmtStateCode = "Idle"
	MgmtStateCodeWaitConnect                = "Listen"
	MgmtStateCodeConnected                  = "Connected"
	MgmtStateCodeDisconnected               = "Disconnected"
)

// ManagerState for management
type ManagerState interface {
	state() MgmtStateCode

	// onEnter is called when first entering the state
	onEnter() error

	// onMessage is called when a message was received from the management interface
	onMessage(messages.Message) error

	// onSignal is called when a signal was raised internally
	onSignal(signals.Signal) error

	// onEvent is called for any other event
	onEvent(events.Event) error

	// onExit is called in last when exiting the state
	onExit() error

	// StateMgmtBasic: unhandledMessage called by onMessage when no handler was found
	unhandledMessage(messages.Message) error
	// StateMgmtBasic: unhandledSignal called by onMessage when no handler was found
	unhandledSignal(signals.Signal) error
	// StateMgmtBasic: unhandledEvent called by onMessage when no handler was found
	unhandledEvent(events.Event) error
}

func errUnexpectedMessage(m messages.Message, state string) error {
	return fmt.Errorf("openvpn-mgmt:state> unexpected message %s in this state: %s", m.Cmd(), state)
}

func errUnexpectedSignal(sig signals.Signal, state string) error {
	return fmt.Errorf("openvpn-mgmt:state> unexpected signal %d in this state: %s", sig, state)
}

func errUnexpectedEvent(e events.Event, state string) error {
	return fmt.Errorf("openvpn-mgmt:state> unexpected event %d in this state: %s", e.Code(), state)
}
