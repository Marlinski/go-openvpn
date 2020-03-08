package events

// InternalEvent are used internally
type InternalEvent Event

// internal events
const (
	InternalEventSignal            EventCode = 900
	InternalEventRecvMsg                     = 901
	InternalEventReadError                   = 902
	InternalEventListenSocket                = 903
	InternalEventListenSocketError           = 904
	InternalEventMgmtConnected               = 905
	InternalEventMgmtDisconnected            = 906
)
