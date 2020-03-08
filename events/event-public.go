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
	OpenvpnEventUp        EventCode = 100
	OpenvpnEventDown                = 101
	OpenvpnEventBytecount           = 102
)
