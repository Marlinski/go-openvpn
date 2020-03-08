package messages

import (
	"fmt"
	"net/textproto"
	"regexp"
	"strconv"
)

// env
var (
	holdRegexp  = regexp.MustCompile("([^:]+):([0-9]+)$")
	ErrNotAHold = fmt.Errorf("hold parsing error")
	MessageHold = "HOLD"
)

func init() {
	messageParsers[MessageHold] = parseHold
}

// HoldMessage returns the hold time
type HoldMessage struct {
	HoldTime int64
}

// Cmd implements Message Cmd interface
func (m HoldMessage) Cmd() string {
	return MessageHold
}

// String implements Message String interface
func (m HoldMessage) String() string {
	return fmt.Sprintf("cmd=%s time=%d", MessageHold, m.HoldTime)
}

func parseHold(cmd string, str string, reader *textproto.Reader) (Message, error) {
	ret := HoldMessage{}

	// read all env
	match := holdRegexp.FindStringSubmatch(str)
	if len(match) != 3 {
		return nil, ErrNotAHold
	}

	ht, err := strconv.ParseInt(match[2], 0, 64)
	if err != nil {
		return nil, ErrNotAHold
	}

	ret.HoldTime = ht
	return ret, nil
}
