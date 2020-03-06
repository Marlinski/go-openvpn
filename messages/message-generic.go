package messages

import (
	"fmt"
	"strings"
)

// GenericMessage doesn't require specific parser
type GenericMessage struct {
	cmd  string
	args []string
}

// Cmd implements Message Cmd interface
func (m GenericMessage) Cmd() string {
	return m.cmd
}

// Args implements Message Args interface
func (m GenericMessage) Args() []string {
	return m.args
}

// String implements Message String interface
func (m GenericMessage) String() string {
	return fmt.Sprintf("command=%s args=%+v", m.Cmd(), m.Args())
}

func parseGeneric(cmd string, str string) (Message, error) {
	args := strings.Split(str, ",")
	return GenericMessage{
		cmd:  cmd,
		args: args,
	}, nil
}
