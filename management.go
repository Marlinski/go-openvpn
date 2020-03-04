package openvpn

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"

	"github.com/rs/xid"
)

// Manager manages an openvpn connection
type Manager struct {
	unixSocket string
	channel    chan *Event
}

// NewManager Creates a new manager ready to accept a connection
func newManager() (*Manager, error) {
	mgr := Manager{}
	mgr.unixSocket = "/tmp/openvpn-" + xid.New().String() + ".sock"
	return &mgr, nil
}

// UpdateConfig to add the management parameters
// https://openvpn.net/community-resources/reference-manual-for-openvpn-2-4/
func (m *Manager) updateConfig(c *Config) {
	c.Set("management", m.unixSocket, "unix") // management unix interface
	c.Flag("management-client")               // Connect to the management interface
	c.Flag("management-hold")                 // Do not connect the tunnel until manager send "hold release"
	c.Flag("management-up-down")              // Report tunnel up/down events to management interface.
}

func (m *Manager) run(channel chan *Event) error {
	socket, err := net.Listen("unix", m.unixSocket)
	if err != nil {
		return fmt.Errorf("%s: could not listen on the socket %+v", m.unixSocket, err)
	}

	go func() {
		defer socket.Close()
		fd, err := socket.Accept()
		if err != nil {
			return
		}

		log.Println("Management: openvpn management interface have connected")
		reader := bufio.NewReader(fd)
		tp := textproto.NewReader(reader)
		for {
			line, err := tp.ReadLine()
			if err != nil {
				return
			}

			log.Println(line)
		}
	}()
	return nil
}
