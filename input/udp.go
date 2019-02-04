package input

import (
	"fmt"
	"net"
	"compress/zlib"
	"io"
	"os"
	"bytes"
)

type udp struct {
	addr     string
	protocol string
	error    string
	active   bool
	enabled  bool
}

func CreateUdp(addr string, enabled bool) AbstractInput {
	return &udp{
		addr:     addr,
		protocol: "udp",
		active:   false,
		enabled:  enabled,
	}
}

// start listen
func (udp *udp) Start() bool {
	ServerAddr, err := net.ResolveUDPAddr("udp", udp.GetAddr())

	if err != nil {
		udp.error = err.Error()
		return false
	}

	ServerConn, err := net.ListenUDP("udp", ServerAddr)

	if err != nil {
		udp.error = err.Error()
	}

	defer ServerConn.Close()

	fmt.Println("Listen UDP on ", udp.GetAddr())

	udp.active = true
	udp.error = ""

	buf := make([]byte, 1024)

	for {
		ServerConn.ReadFromUDP(buf)

		b := bytes.NewReader(buf)
		r, err := zlib.NewReader(b)

		if err != nil {
			panic(err)
		}


		io.Copy(os.Stdout, r)
		r.Close()
	}
}

func (udp *udp) Shutdown() bool {
	return true
}

func (udp *udp) IsActive() bool {
	return udp.active
}

func (udp *udp) IsEnabled() bool {
	return udp.enabled
}

func (udp *udp) GetError() string {
	return udp.error
}

func (udp *udp) GetAddr() string {
	return udp.addr
}

func (udp *udp) GetProtocol() string {
	return udp.protocol
}

func (udp *udp) HasError() bool {
	return udp.error != ""
}
