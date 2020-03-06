package openvpn

import (
	"fmt"
	"time"

	"github.com/Marlinski/go-openvpn/signals"

	"github.com/Marlinski/go-openvpn/events"
	"github.com/Marlinski/go-openvpn/messages"
)

// StateMgmtConnected listen for openvpn command
type StateMgmtConnected struct {
	StateMgmtBasic
	ticker *time.Ticker

	// cache environment variable when tunnel is up
	env map[string]string
}

func newStateMgmtConnected(mgr *Manager) *StateMgmtConnected {
	ret := StateMgmtConnected{
		StateMgmtBasic: newStateMgmtBasic(mgr),
		env:            make(map[string]string),
	}

	// state event map
	ret.msgHandlers["HOLD"] = ret.onMsgHold
	ret.msgHandlers["UPDOWN"] = ret.onMsgTunnelUpDown
	ret.msgHandlers["BYTECOUNT"] = ret.onMsgByteCount
	ret.signalHandlers[signals.SigTick] = ret.onSigTick
	return &ret
}

func (s *StateMgmtConnected) state() MgmtStateCode {
	return MgmtStateCodeConnected
}

// we start listening for the next message
func (s *StateMgmtConnected) onEnter() error {
	go s.mgr.ActionReadMsg()
	return nil
}

// we start listening for the next message
func (s *StateMgmtConnected) onExit() error {
	return nil
}

// Message Receive >HOLD
func (s *StateMgmtConnected) onMsgHold(msg messages.Message) error {
	err := s.mgr.ActionSendCmd("hold release")
	if err != nil {
		return err
	}

	go s.mgr.ActionReadMsg()
	return nil
}

// Message Receive  >UPDOWN
func (s *StateMgmtConnected) onMsgTunnelUpDown(msg messages.Message) error {
	updown, ok := msg.(messages.UPDownMessage)
	if !ok {
		return fmt.Errorf("cast error: expected messages.UpDownMessage")
	}

	if updown.State == "UP" {
		// cache environment variable
		s.env = updown.Env

		s.mgr.upstreamChannel <- events.EventTunnelUp{Env: updown.Env}
		s.mgr.ActionSendCmd("bytecount 5")
	}
	if updown.State == "DOWN" {
		// cache environment variable
		s.env = updown.Env

		s.mgr.upstreamChannel <- events.EventTunnelDown{Env: updown.Env}
	}

	go s.mgr.ActionReadMsg()
	return nil
}

// Message Receive  >BYBTECOUNT
func (s *StateMgmtConnected) onMsgByteCount(msg messages.Message) error {
	bc, ok := msg.(messages.ByteCountMessage)
	if !ok {
		return fmt.Errorf("cast error: expected messages.BytecountMessage")
	}

	s.mgr.upstreamChannel <- events.EventByteCount{ByteCountMessage: bc}
	go s.mgr.ActionReadMsg()
	return nil
}

// if no specific handler we ignore the message and read the next one
func (s *StateMgmtConnected) unhandledMessage(msg messages.Message) error {
	go s.mgr.ActionReadMsg()
	return nil
}

// send status command every 5 seconds
func (s *StateMgmtConnected) onSigTick(sig signals.Signal) error {
	err := s.mgr.ActionSendCmd("status")
	if err != nil {
		return err
	}

	go s.mgr.ActionReadMsg()
	return nil
}
