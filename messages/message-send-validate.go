package messages

import (
	"fmt"
	"regexp"
	"strings"
)

// todo: validate all the commands
var (
	allCommands = map[string]string{
		"bytecount": "bytecount ([0-9]+)$",                         // bytecount n            : Show bytes in/out, update every n secs (0=off).
		"help":      "help$",                                       // help                   : Print this message.
		"hold":      "hold( (on|off|release))?$",                   // hold [on|off|release]  : Set/show hold flag to on/off state, or
		"status":    "status( [0-9]+)?",                            // status [n]             : Show current daemon status info using format #n./
		"signal":    "signal( (SIGHUP|SIGTERM|SIGUSR1|SIGUSR2))?$", // signal s               : Send signal s to daemon (s = SIGHUP|SIGTERM|SIGUSR1|SIGUSR2)
		"version":   "version$",                                    // version                : Show current version number.
		"username":  "username type ([^\r\n]+)$",                   // username type u        : Enter username u for a queried OpenVPN username.
		"password":  "password type ([^\r\n]+)$",                   // password type p        : Enter password p for a queried OpenVPN password.
		"pid":       "pid",                                         // pid                    : Show process ID of the current OpenVPN process.
	}

// auth-retry t           : Auth failure retry mode (none,interact,nointeract).
// echo [on|off] [N|all]  : Like log, but only show messages in echo buffer.
// exit|quit              : Close management session.
// forget-passwords       : Forget passwords entered so far.
//                          release current hold and start tunnel.
// kill cn                : Kill the client instance(s) having common name cn.
// kill IP:port           : Kill the client instance connecting from IP:port.
// load-stats             : Show global server load stats.
// log [on|off] [N|all]   : Turn on/off realtime log display
//                          + show last N lines or 'all' for entire history.
// mute [n]               : Set log mute level to n, or show level if n is absent.
// needok type action     : Enter confirmation for NEED-OK request of 'type',
//                          where action = 'ok' or 'cancel'.
// needstr type action    : Enter confirmation for NEED-STR request of 'type',
//                          where action is reply string.
// net                    : (Windows only) Show network info and routing table.
// remote type [host port] : Override remote directive, type=ACCEPT|MOD|SKIP.
// proxy type [host port flags] : Enter dynamic proxy server info.
// pkcs11-id-count        : Get number of available PKCS#11 identities.
// pkcs11-id-get index    : Get PKCS#11 identity at index.
// client-auth CID KID    : Authenticate client-id/key-id CID/KID (MULTILINE)
// client-auth-nt CID KID : Authenticate client-id/key-id CID/KID
// client-deny CID KID R [CR] : Deny auth client-id/key-id CID/KID with log reason
//                              text R and optional client reason text CR
// client-kill CID [M]    : Kill client instance CID with message M (def=RESTART)
// env-filter [level]     : Set env-var filter level
// client-pf CID          : Define packet filter for client CID (MULTILINE)
// rsa-sig                : Enter an RSA signature in response to >RSA_SIGN challenge
//                          Enter signature base64 on subsequent lines followed by END
// certificate            : Enter a client certificate in response to >NEED-CERT challenge
//                          Enter certificate base64 on subsequent lines followed by END
// state [on|off] [N|all] : Like log, but show state history.
// test n                 : Produce n lines of output for testing/debugging.
// verb [n]               : Set log verbosity level to n, or show if n is absent.
)

// ValidateCommand before sending to the openvpn mgmt interface
func ValidateCommand(command string) (string, error) {
	args := strings.Split(command, " ")

	cmdRegexp, ok := allCommands[args[0]]
	if !ok {
		return "", fmt.Errorf("not a command")
	}

	// check that command is correct
	if match, _ := regexp.MatchString(cmdRegexp, command); !match {
		return "", fmt.Errorf("not a %s command", args[0])
	}

	// add the semicolon
	return command + "\n", nil
}
