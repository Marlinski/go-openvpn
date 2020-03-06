package openvpn

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os/exec"

	"github.com/Marlinski/go-openvpn/events"
	"github.com/Marlinski/go-openvpn/messages"
	"github.com/Marlinski/go-openvpn/signals"
)

// ActionManagerBootstrap runs the event listener and send a start signal
func (m *Manager) ActionManagerBootstrap() {

	// start the looper
	go func() {
		for {
			select {
			case e := <-m.eventChannel:
				m.mux.Lock() /////////// critical section
				//m.log.Debugf("openvpn-mgmt:event> %s", e.String())
				err := m.state.onEvent(e)
				if err != nil {
					m.log.Errorf("%+v", err)
				}
				m.mux.Unlock() ///////////////
			}
		}
	}()

	// send the start signal
	m.eventChannel <- events.EventManagementSignal{Sig: signals.SigStart}
}

// ActionListenManagementSocket listen the unix socket and wait for the openvpn to connect
func (m *Manager) ActionListenManagementSocket() {
	//m.log.Debugf("openvpn-mgmt:action> listening on socket %s", m.conn.unixSocket)
	socket, err := net.Listen("unix", m.conn.unixSocket)
	if err != nil {
		sockErr := fmt.Errorf("openvpn-mgmt:wait-conn> could not listen on the socket %+v", err)
		m.eventChannel <- events.EventManagementListenSocketError{Err: sockErr}
		return
	}

	m.conn.socket = socket
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
	m.eventChannel <- events.EventManagementConnected{UnixSocket: m.conn.unixSocket}
}

// ActionSendCmd sends a single command to the connection
// and reads the response should be either SUCCESS or ERROR
func (m *Manager) ActionSendCmd(command string) error {
	// chomp command
	validated, err := messages.ValidateCommand(command)
	if err != nil {
		return err
	}

	// send the command
	m.conn.fd.Write([]byte(validated))

	// read the response
	line, err := m.conn.tp.ReadLine()
	if err != nil {
		return err
	}

	resp, err := messages.ParseResponse(line, m.conn.tp)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("command returned error: %s" + resp.Msg)
	}

	m.log.Debugf("command succeeded: %s", resp.Msg)
	return nil
}

// ActionReadMsg Read a single message line
func (m *Manager) ActionReadMsg() error {
	line, err := m.conn.tp.ReadLine()
	if err != nil {
		return fmt.Errorf("openvpn-mgmt:read-cmd> could not read socket %+v", err)
	}

	msg, err := messages.ParseMessage(line, m.conn.tp)
	if err != nil {
		return fmt.Errorf("openvpn-mgmt:read-cmd> error - %+v", err)
	}

	m.log.Debugf("openvpn-mgmt:recv> %s", msg.String())
	m.eventChannel <- events.EventManagementRecvMsg{Msg: msg}
	return nil
}

// StartOpenVPN executes the command
func (m *Manager) StartOpenVPN() error {
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
