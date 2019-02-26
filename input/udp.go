package input

import (
	"bytes"
	"compress/zlib"
	"io"
	"net"
	"os"
)

type udp struct {
	Addr   string
	Error  string
	Active bool
	Conn   *net.UDPConn
}

func CreateUdp(addr string) *udp {
	return &udp{
		Addr:   addr,
		Active: false,
	}
}

// start listen
func (udp *udp) Start() {
	ServerAddr, err := net.ResolveUDPAddr("udp", udp.Addr)

	if err != nil {
		udp.Error = err.Error()
		return
	}

	ServerConn, err := net.ListenUDP("udp", ServerAddr)

	if err != nil {
		udp.Error = err.Error()
	}

	udp.Conn = ServerConn
	udp.Active = true
	udp.Error = ""

	defer udp.Stop()

	buf := make([]byte, 1024)

	for {
		_, _, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			udp.Error = err.Error()
			break
		}

		b := bytes.NewReader(buf)
		r, err := zlib.NewReader(b)

		if err == nil {
			io.Copy(os.Stdout, r)
			r.Close()
		}
	}
}

func (udp *udp) Stop() {
	if udp.Active {
		udp.Conn.Close()
	}
	udp.Active = false
}
