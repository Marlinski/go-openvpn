package events

// Event regarding vpn status
type Event interface {
	Code() EventCode
	String() string
}

// EventCode is an int
type EventCode int

// OpenvpnEvent are used by upstream
type OpenvpnEvent Event

// public events
const (
	OpenvpnEventUp       EventCode = 100
	OpenvpnEventDown               = 101
	OpenvpnEventByteCode           = 102
)

// InternalEvent are used internally
type InternalEvent Event

// internal events
const (
	InternalEventSignal            EventCode = 900
	InternalEventRecvMsg                     = 901
	InternalEventListenSocket                = 902
	InternalEventListenSocketError           = 903
	InternalEventMgmtConnected               = 904
	InternalEventMgmtDisconnected            = 905
)
