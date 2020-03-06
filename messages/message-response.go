package messages

import (
	"fmt"
	"net/textproto"
)

// ParseResponse a response to a command
func ParseResponse(reply string, reader *textproto.Reader) (Response, error) {
	// check that command is correct
	match := responseRegexp.FindStringSubmatch(reply)
	if match == nil {
		return Response{}, fmt.Errorf("not a response")
	}

	var success bool
	switch match[1] {
	case "SUCCESS":
		success = true
	case "ERROR":
		success = false
	default:
		return Response{}, fmt.Errorf("not a response")
	}

	return Response{
		Success: success,
		Msg:     match[2],
	}, nil
}
