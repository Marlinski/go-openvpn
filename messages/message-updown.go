package messages

import (
	"fmt"
	"net/textproto"
	"regexp"
)

// env
var (
	envRegexp     = regexp.MustCompile(">UPDOWN:ENV,([^=]+)=(.*)")
	ErrNotAnEnv   = fmt.Errorf("parsing error: expected environment value")
	MessageUpDown = "UPDOWN"
)

func init() {
	messageParsers[MessageUpDown] = parseUpDown
}

// UPDownMessage contains all the environoment variable
// spanned on a multiline reply
type UPDownMessage struct {
	State string
	Env   map[string]string
}

// Cmd implements Message Cmd interface
func (m UPDownMessage) Cmd() string {
	return MessageUpDown
}

// String implements Message String interface
func (m UPDownMessage) String() string {
	return fmt.Sprintf("cmd=UPDOWN state=%s", m.State)
}

func parseUpDown(cmd string, str string, reader *textproto.Reader) (Message, error) {
	ret := UPDownMessage{
		State: str,
		Env:   make(map[string]string),
	}

	// read all env
	for {
		line, err := reader.ReadLine()
		if err != nil {
			return ret, err
		}

		if line == ">UPDOWN:ENV,END" {
			return ret, nil
		}

		match := envRegexp.FindStringSubmatch(line)
		if len(match) != 3 {
			return nil, ErrNotAnEnv
		}

		ret.Env[match[1]] = match[2]
	}
}
