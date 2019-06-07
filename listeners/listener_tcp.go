package listeners

import (
	"net"
	"microlog/storage"
	"fmt"
)

const PROTOCOL_TCP = "tcp"

type tcp struct {
	addr      string
	error     string
	active    bool
	conn      net.Listener
	extractor Extractor
	storage   storage.Storage
}

func CreateTcp(addr string, extractor Extractor, storage *storage.Storage) Listener {
	return &tcp{
		addr:      addr,
		active:    false,
		extractor: extractor,
		storage:   *storage,
	}
}

// start listen
func (tcp *tcp) Start() {
	go (func() {
		listener, err := net.Listen("tcp", tcp.addr)

		if err != nil {
			tcp.error = err.Error()
			return
		}

		tcp.conn = listener
		tcp.active = true
		tcp.error = ""

		defer tcp.Stop()

		for {
			conn, err := listener.Accept()

			if err != nil {
				tcp.error = err.Error()
				continue
			}

			go func(conn net.Conn) {
				fmt.Println("Open TCP connection")

				defer func() {
					fmt.Println("Close TCP connection")
					conn.Close()
				}()

				for {
					input := make([]byte, 1024*1024)
					n, err := conn.Read(input)
					if err != nil {
						tcp.error = err.Error()
						break
					}

					row, err := tcp.extractor.Extract(input[0:n])
					row["remote_addr"] = conn.RemoteAddr().String()

					if err == nil {
						tcp.storage.Write(row)
					} else {
						tcp.error = err.Error()
					}
				}
			}(conn)
		}
	})()
}

// stop listen
func (tcp *tcp) Stop() {
	if tcp.active {
		if tcp.conn != nil {
			_ = tcp.conn.Close()
		}
	}
	tcp.active = false
	tcp.conn = nil
}

func (tcp *tcp) IsActive() bool {
	return tcp.active
}

func (tcp *tcp) GetError() string {
	return tcp.error
}

func (tcp *tcp) GetAddr() string {
	return tcp.addr
}
