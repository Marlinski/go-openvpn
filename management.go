package openvpn

import (
	"bufio"
	"net"
	"net/textproto"
	"os/exec"
	"sync"

	"github.com/op/go-logging"
	"github.com/rs/xid"

	"github.com/Marlinski/go-openvpn/events"
	"github.com/Marlinski/go-openvpn/util"
)

// Manager manages the management interface for a single openvpn config
type Manager struct {
	config          Config                    // openvpn config
	cmd             *exec.Cmd                 // unix command
	state           ManagerState              // current state
	upstreamChannel chan events.OpenvpnEvent  // upstream event receiver
	eventChannel    chan events.InternalEvent // event channel
	log             *logging.Logger           // debug logger
	conn            ManagerInterface          // interface
	mux             sync.Mutex                // mutex
}

// StateMachine machine
type StateMachine struct {
	idle         *StateMgmtIdle
	listen       *StateMgmtWaitConnect
	connected    *StateMgmtConnected
	disconnected *StateMgmtDisconnected
}

// ManagerInterface holds the parameters for the interface
type ManagerInterface struct {
	unixSocket string            // unix socket path
	socket     net.Listener      // socket to accept connection from openvpn mgmt
	fd         net.Conn          // actual communication socket with openvpn management
	reader     *bufio.Reader     // buffer reader
	tp         *textproto.Reader // text reader
}

// NewManager Creates a new manager ready to accept a connection
func (c Config) NewManager(upstream chan events.OpenvpnEvent) *Manager {
	mgr := Manager{
		config:          c,
		upstreamChannel: upstream,
		eventChannel:    make(chan events.InternalEvent),
		log:             util.CreateLeveledLog(c.id, c.logLevel),
		conn: ManagerInterface{
			unixSocket: "/tmp/openvpn-" + xid.New().String() + ".sock",
		},
	}
	mgr.stateMgmtIdle()

	// update config parameter to enable management
	// https://openvpn.net/community-resources/reference-manual-for-openvpn-2-4/
	mgr.config.Set("management", mgr.conn.unixSocket, "unix") // management unix socket interface
	mgr.config.Flag("management-client")                      // Openvpn must connect to the management interface
	mgr.config.Flag("management-hold")                        // Openvpn must not connect the tunnel until manager send "hold release"
	mgr.config.Flag("management-up-down")                     // Openvpn must report tunnel up/down events to management interface.
	return &mgr
}

func (m *Manager) stateMgmtIdle() error {
	return m.changeState(newStateMgmtIdle(m))
}

func (m *Manager) stateMgmtWaitConnect() error {
	return m.changeState(newStateMgmtWaitConnect(m))
}

func (m *Manager) stateMgmtConnected() error {
	return m.changeState(newStateMgmtConnected(m))
}

func (m *Manager) stateMgmtDisconnected() error {
	return m.changeState(newStateMgmtDisconnected(m))
}

func (m *Manager) changeState(newState ManagerState) error {
	if m.state == nil {
		m.log.Debugf("openvpn-mgmt:state> initialize to state %s", newState.state())
	} else {
		m.log.Debugf("openvpn-mgmt:state> change from state %s to state %s", m.state.state(), newState.state())
		err := m.state.onExit()
		if err != nil {
			return err
		}
	}
	m.state = newState
	return m.state.onEnter()
}
