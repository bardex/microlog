package input

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"net"
	"os"
)

type udp struct {
	id       int
	addr     string
	protocol string
	error    string
	active   bool
	enabled  bool
	ch       chan int
}

func CreateUdp(id int, addr string, enabled bool) AbstractInput {
	return &udp{
		id:       id,
		addr:     addr,
		protocol: "udp",
		active:   false,
		enabled:  enabled,
		ch:       make(chan int),
	}
}

// start listen
func (udp *udp) Start() {
	ServerAddr, err := net.ResolveUDPAddr("udp", udp.GetAddr())

	if err != nil {
		udp.error = err.Error()
		return
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

loop:
	for {
		select {
		case signal := <-udp.ch:
			switch signal {
			case SIGNAL_STOP:
				break loop
			}
		default:
			ServerConn.ReadFromUDP(buf)

			b := bytes.NewReader(buf)
			r, err := zlib.NewReader(b)

			if err == nil {
				io.Copy(os.Stdout, r)
				r.Close()
			}
		}
	}
}

func (udp *udp) Shutdown() {
	if udp.active {
		udp.ch <- SIGNAL_STOP
		udp.active = false
	}
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

func (udp *udp) GetId() int {
	return udp.id
}

func (udp *udp) GetProtocol() string {
	return udp.protocol
}

func (udp *udp) HasError() bool {
	return udp.error != ""
}
