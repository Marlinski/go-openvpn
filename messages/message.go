package messages

import (
	"fmt"
	"net/textproto"
	"regexp"
)

// Message from openvpn management
type Message interface {
	Cmd() string
	String() string
}

// Parser takes a string and returns a message
type Parser func(string, string, *textproto.Reader) (Message, error)

// regexp
var (
	MsgRegexp      = regexp.MustCompile(">([\\w]+):([^\r\n]*)$")
	ErrNotAMessage = fmt.Errorf("not a message")
)

// specific parsers
var (
	messageParsers = map[string]Parser{}
)

// ParseMessage received from an openvpn management interface
func ParseMessage(str string, reader *textproto.Reader) (Message, error) {
	match := MsgRegexp.FindStringSubmatch(str)
	if len(match) != 3 {
		return nil, ErrNotAMessage
	}

	cmd := match[1]
	args := match[2]
	parseSpecific, ok := messageParsers[cmd]
	if !ok {
		return parseGeneric(cmd, args)
	}

	return parseSpecific(cmd, args, reader)
}
