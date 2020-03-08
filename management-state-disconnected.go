package openvpn

// StateMgmtDisconnected connected state
type StateMgmtDisconnected struct {
	StateMgmtBasic
}

func newStateMgmtDisconnected(mgr *Manager) *StateMgmtDisconnected {
	ret := StateMgmtDisconnected{
		newStateMgmtBasic(mgr, MgmtStateCodeDisconnected),
	}
	return &ret
}
