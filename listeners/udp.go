package listeners

import (
	"net"
)

const ProtocolUdp = "udp"

type UDPHandler struct {
	conn *net.UDPConn
}

// start listen
func (udp *UDPHandler) Listen(listener *Listener) error {

	ServerAddr, err := net.ResolveUDPAddr("udp", listener.Addr)

	if err != nil {
		return err
	}

	ServerConn, err := net.ListenUDP("udp", ServerAddr)

	if err != nil {
		return err
	}

	udp.conn = ServerConn
	defer udp.Close()

	buf := make([]byte, 4*1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			listener.Error = err.Error()
			continue
		}

		row, err := listener.Extractor.Extract(buf[:n])
		row["remote_addr"] = addr.String()

		if err == nil {
			listener.Storage.Write(row)
		} else {
			listener.Error = err.Error()
		}
	}

}

// stop listen
func (udp *UDPHandler) Close() {
	if udp.conn != nil {
		_ = udp.conn.Close()
	}
	udp.conn = nil
}
