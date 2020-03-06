package signals

// Signal is just an int
type Signal int

// all the signals
const (
	SigStart Signal = 0
	SigStop         = 1
	SigTick         = 2
)
