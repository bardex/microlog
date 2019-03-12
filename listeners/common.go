package listeners

import "fmt"

const PROTOCOL_UDP = "udp"
const PROTOCOL_TCP = "tcp"

type Listener interface {
	Start()
	Stop()
}

type nilListener struct{}

func (l nilListener) Start() {}
func (l nilListener) Stop()  {}

func CreateListener(protocol string, addr string, extractor string) (Listener, error) {
	switch protocol {
	case PROTOCOL_UDP:
		return CreateUdp(addr, extractor)
	}
	return nilListener{}, fmt.Errorf("listener %s not found", protocol)
}
