package openvpn

import (
	"github.com/Marlinski/go-openvpn/events"
)

// StateMgmtWaitConnect wait for a connection from openvpn management
type StateMgmtWaitConnect struct {
	StateMgmtBasic
}

func newStateMgmtWaitConnect(mgr *Manager) *StateMgmtWaitConnect {
	ret := StateMgmtWaitConnect{
		newStateMgmtBasic(mgr, MgmtStateCodeWaitConnect),
	}

	// state event map
	ret.eventHandlers[events.InternalEventListenSocketError] = ret.onListenError
	ret.eventHandlers[events.InternalEventMgmtConnected] = ret.onMgmtConnected

	return &ret
}

func (s *StateMgmtWaitConnect) onEnter() error {
	err := s.mgr.ActionListenManagementSocket()
	if err != nil {
		return err
	}

	err = s.mgr.ActionStartOpenVPN()
	if err != nil {
		// this will trigger onListenError
		s.mgr.conn.socket.Close()
		return err
	}
	return nil
}

func (s *StateMgmtWaitConnect) onListenError(e events.Event) error {
	return nil
}

func (s *StateMgmtWaitConnect) onMgmtConnected(e events.Event) error {
	return s.mgr.stateMgmtConnected()
}
