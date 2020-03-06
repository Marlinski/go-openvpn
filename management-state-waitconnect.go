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
		newStateMgmtBasic(mgr),
	}

	// state event map
	ret.eventHandlers[events.InternalEventListenSocketError] = ret.onListenError
	ret.eventHandlers[events.InternalEventMgmtConnected] = ret.onMgmtConnected

	return &ret
}

func (s *StateMgmtWaitConnect) state() MgmtStateCode {
	return MgmtStateCodeWaitConnect
}

func (s *StateMgmtWaitConnect) onEnter() error {
	go s.mgr.ActionListenManagementSocket()
	err := s.mgr.StartOpenVPN()
	if err != nil {
		// this will trigger onListenError
		s.mgr.conn.socket.Close()
	}
	return nil
}

func (s *StateMgmtWaitConnect) onListenError(e events.Event) error {
	return nil
}

func (s *StateMgmtWaitConnect) onMgmtConnected(e events.Event) error {
	return s.mgr.stateMgmtConnected()
}
