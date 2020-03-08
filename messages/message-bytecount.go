package messages

import (
	"fmt"
	"net/textproto"
	"regexp"
	"strconv"
)

// env
var (
	inOutRegexp      = regexp.MustCompile("([0-9]+),([0-9]+)$")
	ErrNotAByteCount = fmt.Errorf("bytecount parsing error")
	MessageBytecount = "BYTECOUNT"
)

func init() {
	messageParsers[MessageBytecount] = parseByteCount
}

// BytecountMessage returns bytes in/out
type BytecountMessage struct {
	In  int64
	Out int64
}

// Cmd implements Message Cmd interface
func (m BytecountMessage) Cmd() string {
	return MessageBytecount
}

// String implements Message String interface
func (m BytecountMessage) String() string {
	return fmt.Sprintf("cmd=%s in=%d out=%d", MessageBytecount, m.In, m.Out)
}

func parseByteCount(cmd string, str string, reader *textproto.Reader) (Message, error) {
	ret := BytecountMessage{}

	// read all env
	match := inOutRegexp.FindStringSubmatch(str)
	if len(match) != 3 {
		return nil, ErrNotAByteCount
	}

	in, err := strconv.ParseInt(match[1], 0, 64)
	if err != nil {
		return nil, ErrNotAByteCount
	}

	out, err := strconv.ParseInt(match[2], 0, 64)
	if err != nil {
		return nil, ErrNotAByteCount
	}

	ret.In = in
	ret.Out = out
	return ret, nil
}
