package input

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestUdp(t *testing.T) {
	udp := CreateUdp(1, ":8080", true)
	go udp.Start()

	time.Sleep(1 * time.Second)

	for i := 0; i < 100; i++ {
		udpSend(strconv.Itoa(i) + "\n")
		if i >= 50 {
			udp.Shutdown()
		}
	}

}

func udpSend(msg string) error {
	conn, err := net.Dial("udp", ":8080")
	if err != nil {
		return err
	}

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write([]byte(msg))
	w.Close()

	fmt.Fprintf(conn, b.String())
	conn.Close()
	return nil
}
