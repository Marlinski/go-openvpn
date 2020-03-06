package openvpn

import (
	"github.com/Marlinski/go-openvpn/signals"
)

// StateMgmtIdle idle state
type StateMgmtIdle struct {
	StateMgmtBasic
}

func newStateMgmtIdle(mgr *Manager) *StateMgmtIdle {
	ret := StateMgmtIdle{
		newStateMgmtBasic(mgr),
	}

	// state event map
	ret.signalHandlers[signals.SigStart] = ret.onSigStart

	return &ret
}

func (s *StateMgmtIdle) state() MgmtStateCode {
	return MgmtStateCodeIdle
}

func (s *StateMgmtIdle) onSigStart(sig signals.Signal) error {
	return s.mgr.stateMgmtWaitConnect()
}
