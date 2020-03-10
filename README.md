go-openvpvn
===========

go-openvpn is a lightweight golang library to start and monitor an openvpn tunnel. It essentially just wraps around the openvpn process and interacts with it using the [management interface](https://openvpn.net/community-resources/management-interface/). Right now only openvpn client configuration has been implemented and tested.

It supports the following features:

* simple API
* use an existing openvpn configuration
* override the config file settings programmatically
* go channel to receive events:
    * tunnel status (up/down)
    * openvpn statistics
* send management command

# Import 

```
import (
	"github.com/Marlinski/go-openvpn"
	"github.com/Marlinski/go-openvpn/events"
)
```

# Basic example

The following code shows how to run a simple openvpn:

```
cfg := openvpn.LoadConfig(m.iface, "test.ovpn") // openvpn config file
cfg.SetLogStd(true)                             // log openvpn standard output and error
cfg.Set("dev", "tun-test)                       // passed to the openvpn command line parameter --dev tun-test

channel = make(chan events.OpenvpnEvent)        // channel to process the events 
go processEvents()                              // event loop

ctrl = cfg.Run(m.channel)                       // run openvpvn, it returns a controller 
```

And the event loop processor:

```
func processEvents() {
    for {
        e := <-m.channel
        if e.Code() == events.OpenvpnEventUp {
            // tunnel is up
        }
        if e.Code() == events.OpenvpnEventDown {
            // tunnel is down
        }
    }
}
```


