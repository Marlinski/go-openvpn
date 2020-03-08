package openvpn

import (
	"fmt"
	"time"

	"github.com/Marlinski/go-openvpn/events"
	"github.com/Marlinski/go-openvpn/messages"
)

// StateMgmtConnected listen for openvpn command
type StateMgmtConnected struct {
	StateMgmtBasic

	tunnelUp bool              // actual tunnel status flag
	env      map[string]string // cache environment variable when tunnel is up
}

func newStateMgmtConnected(mgr *Manager) *StateMgmtConnected {
	ret := StateMgmtConnected{
		StateMgmtBasic: newStateMgmtBasic(mgr, MgmtStateCodeConnected),
		tunnelUp:       false,
		env:            make(map[string]string),
	}

	// state event map
	ret.msgHandlers["HOLD"] = ret.onMsgHold
	ret.msgHandlers["UPDOWN"] = ret.onMsgTunnelUpDown
	ret.msgHandlers["BYTECOUNT"] = ret.onMsgByteCount
	return &ret
}

// we start listening for the next message
func (s *StateMgmtConnected) onEnter() error {
	s.mgr.ActionManagerStartReceiver()
	return nil
}

// Message Receive >HOLD
func (s *StateMgmtConnected) onMsgHold(msg messages.Message) error {
	h, ok := msg.(messages.HoldMessage)
	if !ok {
		return fmt.Errorf("cast error: expected messages.HoldMessage")
	}

	// openvpn starts up with hold flag up because of the --management-hold option
	// in this case holdtime is 0 and openvpn is waiting for the hold release instruction
	// to keep going.
	if h.HoldTime == 0 {
		_, err := s.mgr.ActionSendCmd("hold release")
		if err != nil {
			return err
		}

		return nil
	}

	// if for some reason the connection with remote is down, the hold flag will also be set.
	// but the holdTime will be > 0
	if s.tunnelUp {
		// so we first send an event upstream to tell that the network is down
		s.tunnelUp = false
		s.mgr.upstreamChannel <- events.EventTunnelDown{Env: s.env}
	}

	// now if we send hold release straight away and it still cannot connect, the hold flag will immediately
	// be send back to us leading to a loop of hold / hold release / hold etc..
	// so we time it to avoid running in circle
	go func() {
		<-time.After(5 * time.Second)
		s.mgr.ActionSendCmd("hold release")
	}()
	return nil
}

// Message Receive  >UPDOWN:UP or >UPDOWN:DOWN with environment variable
func (s *StateMgmtConnected) onMsgTunnelUpDown(msg messages.Message) error {
	updown, ok := msg.(messages.UPDownMessage)
	if !ok {
		return fmt.Errorf("cast error: expected messages.UpDownMessage")
	}

	if updown.State == "UP" {
		// cache environment variable
		s.env = updown.Env

		if !s.tunnelUp {
			s.tunnelUp = true
			s.mgr.upstreamChannel <- events.EventTunnelUp{Env: updown.Env}
			s.mgr.ActionSendCmd("bytecount 5")
		}
	}
	if updown.State == "DOWN" {
		// cache environment variable
		s.env = updown.Env

		if s.tunnelUp {
			s.tunnelUp = false
			s.mgr.upstreamChannel <- events.EventTunnelDown{Env: updown.Env}
		}
	}

	return nil
}

// Message Receive  >BYBTECOUNT,in,out
func (s *StateMgmtConnected) onMsgByteCount(msg messages.Message) error {
	bc, ok := msg.(messages.BytecountMessage)
	if !ok {
		return fmt.Errorf("cast error: expected messages.BytecountMessage")
	}

	s.mgr.upstreamChannel <- events.EventBytecount{BytecountMessage: bc}
	return nil
}

// if no specific handler we ignore the message and read the next one
func (s *StateMgmtConnected) unhandledMessage(msg messages.Message) error {
	return nil
}
