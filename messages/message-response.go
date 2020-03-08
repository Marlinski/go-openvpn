package messages

import (
	"log"
	"net/textproto"
	"regexp"
	"strings"
)

// Response from openvpn management
type Response struct {
	Success bool
	Msg     string
}

var (
	successOrErrorRegexp = regexp.MustCompile("([^\r\n]+):(.*)$")
)

// ParseResponse a response to a command
// Response is either SUCCESS: or ERROR: or a wall of multiline texts that ends with END
func ParseResponse(reply string, reader *textproto.Reader) (Response, error) {
	match := successOrErrorRegexp.FindStringSubmatch(reply)

	// success or error response
	if match != nil {
		return Response{
			Success: match[1] == "SUCCESS",
			Msg:     match[2],
		}, nil
	}

	// read each line until an "END"
	line := reply
	for {
		log.Printf("<<< %s", line)
		if strings.HasSuffix(line, "END") {
			return Response{
				Success: true,
				Msg:     reply,
			}, nil
		}

		// read next line
		next, err := reader.ReadLine()
		if err != nil {
			return Response{}, err
		}
		line = next

		// add the line to the response
		reply = reply + "\n" + line
	}
}
