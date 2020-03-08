package openvpn

import (
	"github.com/Marlinski/go-openvpn/events"
	"github.com/Marlinski/go-openvpn/messages"
	"github.com/Marlinski/go-openvpn/signals"
)

// MsgHandler handles openvpn management message
type MsgHandler func(msg messages.Message) error

// EventHandler handles general event
type EventHandler func(e events.Event) error

// SignalHandler handles general signal
type SignalHandler func(s signals.Signal) error

// StateMgmtBasic basic state separates events in three category for better visibility:
// - message
// - signal
// - events
type StateMgmtBasic struct {
	mgr            *Manager
	stateName      MgmtStateCode
	eventHandlers  map[events.EventCode]EventHandler
	msgHandlers    map[string]MsgHandler
	signalHandlers map[signals.Signal]SignalHandler
}

func newStateMgmtBasic(mgr *Manager, state MgmtStateCode) StateMgmtBasic {
	return StateMgmtBasic{
		mgr:            mgr,
		stateName:      state,
		msgHandlers:    make(map[string]MsgHandler),
		eventHandlers:  make(map[events.EventCode]EventHandler),
		signalHandlers: make(map[signals.Signal]SignalHandler),
	}
}

func (s *StateMgmtBasic) state() MgmtStateCode {
	return s.stateName
}

func (s *StateMgmtBasic) onEnter() error {
	return nil
}

func (s *StateMgmtBasic) onEvent(e events.Event) error {
	// route message to handler
	if e.Code() == events.InternalEventRecvMsg {
		msgEvent := e.(events.EventManagementRecvMsg)
		return s.mgr.state.onMessage(msgEvent.Msg)
	}

	// route signal to handler
	if e.Code() == events.InternalEventSignal {
		sigEvent := e.(events.EventManagementSignal)
		return s.mgr.state.onSignal(sigEvent.Sig)
	}

	// route event to handler
	handler, ok := s.eventHandlers[e.Code()]
	if !ok {
		return s.mgr.state.unhandledEvent(e)
	}
	return handler(e)
}

func (s *StateMgmtBasic) onMessage(msg messages.Message) error {
	handler, ok := s.msgHandlers[msg.Cmd()]
	if !ok {
		return s.mgr.state.unhandledMessage(msg)
	}
	return handler(msg)
}

func (s *StateMgmtBasic) onSignal(sig signals.Signal) error {
	handler, ok := s.signalHandlers[sig]
	if !ok {
		return s.mgr.state.unhandledSignal(sig)
	}
	return handler(sig)
}

func (s *StateMgmtBasic) onExit() error {
	return nil
}

func (s *StateMgmtBasic) unhandledMessage(msg messages.Message) error {
	return errUnexpectedMessage(msg, string(s.mgr.state.state()))
}

func (s *StateMgmtBasic) unhandledSignal(sig signals.Signal) error {
	return errUnexpectedSignal(sig, string(s.mgr.state.state()))
}

func (s *StateMgmtBasic) unhandledEvent(e events.Event) error {
	return errUnexpectedEvent(e, string(s.mgr.state.state()))
}
