package listeners

import (
	"net"
)

type udp struct {
	addr      string
	error     string
	active    bool
	conn      *net.UDPConn
	extractor Extractor
	writer    Writer
}

func CreateUdp(addr string, extractor Extractor, writer Writer) Listener {
	return &udp{
		addr:      addr,
		active:    false,
		extractor: extractor,
		writer:    writer,
	}
}

// start listen
func (udp *udp) Start() {
	go (func() {
		ServerAddr, err := net.ResolveUDPAddr("udp", udp.addr)

		if err != nil {
			udp.error = err.Error()
			return
		}

		ServerConn, err := net.ListenUDP("udp", ServerAddr)

		if err != nil {
			udp.error = err.Error()
		}

		udp.conn = ServerConn
		udp.active = true
		udp.error = ""
		defer udp.Stop()

		buf := make([]byte, 1024*1024)

		for {
			n, addr, err := ServerConn.ReadFromUDP(buf)

			if err != nil {
				udp.error = err.Error()
				break
			}
			row, err := udp.extractor.Extract(buf[:n])
			row["remote_addr"] = addr.String()

			if err == nil {
				udp.writer.Write(row)
			} else {
				udp.error = err.Error()
			}
		}
	})()
}

// stop listen
func (udp *udp) Stop() {
	if udp.Active {
		if udp.conn != nil {
			_ = udp.conn.Close()
		}
	}
	udp.Active = false
	udp.conn = nil
	udp.active = false
}

func (udp *udp) IsActive() bool {
	return udp.active
}

func (udp *udp) GetError() string {
	return udp.error
}

func (udp *udp) GetAddr() string {
	return udp.addr
}
