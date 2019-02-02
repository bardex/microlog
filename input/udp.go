package input

import (
	"fmt"
	"log"
	"net"
)

type Udp struct {
	Addr string
}

func (udp *Udp) Start() {
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, err := net.ResolveUDPAddr("udp", udp.Addr)

	if err != nil {
		log.Fatal(err)
	}

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)

	if err != nil {
		log.Fatal(err)
	}

	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

}
