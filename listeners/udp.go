package listeners

import (
	"net"
)

type udp struct {
	Addr   string
	Error  string
	Active bool
	conn   *net.UDPConn
	extractor Extractor
}

func CreateUdp(addr string, extractor string) Listener {
	return &udp{
		Addr:   addr,
		Active: false,
		extractor: createExtractor(extractor),
	}
}

// start listen
func (udp *udp) Start() {
	go (func() {
		ServerAddr, err := net.ResolveUDPAddr("udp", udp.Addr)

		if err != nil {
			udp.Error = err.Error()
			return
		}

		ServerConn, err := net.ListenUDP("udp", ServerAddr)

		if err != nil {
			udp.Error = err.Error()
		}

		udp.conn = ServerConn
		udp.Active = true
		udp.Error = ""
		defer udp.Stop()

		buf := make([]byte, 64*1024*1024)

		for {
			_, _, err := ServerConn.ReadFromUDP(buf)

			if err != nil {
				udp.Error = err.Error()
				break
			}

			udp.extractor.Extract(buf)
		}
	})()
}

func (udp *udp) Stop() {
	if udp.Active {
		udp.conn.Close()
	}
	udp.Active = false
}
