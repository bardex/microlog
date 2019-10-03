package listeners

import (
	"net"
)

const ProtocolTcp = "tcp"

type TCPHandler struct {
	conn net.Listener
}

// start listen
func (tcp *TCPHandler) Listen(listener *Listener) error {
	TCPConn, err := net.Listen("tcp", listener.Addr)

	if err != nil {
		return err
	}

	tcp.conn = TCPConn

	defer tcp.Close()

	for {
		conn, err := TCPConn.Accept()

		if err != nil {
			listener.Error = err.Error()
			continue
		}

		go tcp.handleConn(conn, listener)
	}
}

func (tcp *TCPHandler) handleConn(conn net.Conn, listener *Listener) {
	defer conn.Close()

	for {
		input := make([]byte, 4*1024)
		n, err := conn.Read(input)

		if err != nil {
			listener.Error = err.Error()
			break
		}

		row, err := listener.Extractor.Extract(input[0:n])
		row["remote_addr"] = conn.RemoteAddr().String()

		if err == nil {
			listener.Storage.Write(row)
		} else {
			listener.Error = err.Error()
		}
	}
}

// stop listen
func (tcp *TCPHandler) Close() {

	if tcp.conn != nil {
		_ = tcp.conn.Close()
	}

	tcp.conn = nil
}
