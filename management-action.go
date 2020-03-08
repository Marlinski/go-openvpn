package openvpn

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os/exec"
	"strings"
	"time"

	"github.com/Marlinski/go-openvpn/events"
	"github.com/Marlinski/go-openvpn/messages"
	"github.com/Marlinski/go-openvpn/signals"
)

// ActionManagerBootstrap runs the event listener and send a start signal
func (m *Manager) ActionManagerBootstrap() {

	// start the event looper
	go func() {
		for {
			select {
			case e := <-m.eventChannel:
				m.statemux.Lock() /////////// critical section changes the manager states
				m.log.Debugf("openvpn-mgmt:event> %s", e.String())
				err := m.state.onEvent(e)
				if err != nil {
					m.log.Errorf("%+v", err)
				}
				m.statemux.Unlock() ///////////////
			}
		}
	}()

	// send the start signal
	m.eventChannel <- events.EventManagementSignal{Sig: signals.SigStart}
}

// ActionManagerStartReceiver run the packet listener
func (m *Manager) ActionManagerStartReceiver() {

	// start the receiver looper
	go func() {
		for {
			line, err := m.conn.tp.ReadLine()
			if err != nil {
				readErr := fmt.Errorf("openvpn-mgmt:read> could not read socket %+v", err)
				m.eventChannel <- events.EventManagementReadError{Err: readErr}
				return
			}

			if strings.HasPrefix(line, ">") {
				msg, err := messages.ParseMessage(line, m.conn.tp)
				if err != nil {
					parseErr := fmt.Errorf("openvpn-mgmt:read-cmd> error - %+v", err)
					m.eventChannel <- events.EventManagementReadError{Err: parseErr}
					continue // should this be fatal ?
				}
				m.eventChannel <- events.EventManagementRecvMsg{Msg: msg}
			} else {
				resp, err := messages.ParseResponse(line, m.conn.tp)
				if err != nil {
					parseErr := fmt.Errorf("openvpn-mgmt:read-resp> error - %+v", err)
					m.eventChannel <- events.EventManagementReadError{Err: parseErr}
					continue // should be this fatal ?
				}
				m.conn.respChannel <- resp
			}
		}
	}()
}

// ActionListenManagementSocket listen the unix socket and wait for the openvpn to connect
func (m *Manager) ActionListenManagementSocket() error {
	//m.log.Debugf("openvpn-mgmt:action> listening on socket %s", m.conn.unixSocket)
	socket, err := net.Listen("unix", m.conn.unixSocket)
	if err != nil {
		return fmt.Errorf("openvpn-mgmt:wait-conn> could not listen on the socket %+v", err)
	}

	m.conn.socket = socket
	go func() {
		defer socket.Close()
		fd, err := socket.Accept()
		if err != nil {
			sockErr := fmt.Errorf("openvpn-mgmt:wait-conn> could not accept %+v", err)
			m.eventChannel <- events.EventManagementListenSocketError{Err: sockErr}
			return
		}

		m.conn.fd = fd
		m.conn.reader = bufio.NewReader(m.conn.fd)
		m.conn.tp = textproto.NewReader(m.conn.reader)
		m.conn.respChannel = make(chan messages.Response)
		m.eventChannel <- events.EventManagementConnected{UnixSocket: m.conn.unixSocket}
	}()

	return nil
}

// ActionSendCmd sends a single command to the connection
// and reads the response should be either SUCCESS or ERROR
func (m *Manager) ActionSendCmd(command string) (messages.Response, error) {
	// validate command
	validated, err := messages.ValidateCommand(command)
	if err != nil {
		return messages.Response{}, fmt.Errorf("openvpn-mgmt:send> %+v", err)
	}

	// lock the response channel
	m.conn.respMux.Lock()
	defer m.conn.respMux.Unlock()

	// send the command
	_, err = m.conn.fd.Write([]byte(validated))
	if err != nil {
		return messages.Response{}, fmt.Errorf("openvpn-mgmt:send> %+v", err)
	}

	// listen the response or timeout
	select {
	case resp := <-m.conn.respChannel:
		m.log.Debugf("openvpn-mgmt:send> command succeeded: %s", validated)
		return resp, nil
	case <-time.After(5 * time.Second): // hopefully that should never happen
		return messages.Response{}, fmt.Errorf("openvpn-mgmt:send> command timeout")
	}
}

// ActionStartOpenVPN executes the command
func (m *Manager) ActionStartOpenVPN() error {
	m.cmd = exec.Command("openvpn", m.config.params...)

	// log the standard output/err
	if m.config.logStd {
		stdout, _ := m.cmd.StdoutPipe()
		m.monitorStd(stdout, m.config.id)
		stderr, _ := m.cmd.StderrPipe()
		m.monitorStd(stderr, m.config.id)
	}

	return m.cmd.Start()
}

func (m *Manager) monitorStd(reader io.ReadCloser, id string) {
	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			m.log.Debugf("openvpn-mgmt:std> %s", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			m.log.Debugf("openvpn-mgmt:std> %+v", scanner.Err())
		}
	}()
}
