package listeners

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestUdp(t *testing.T) {
	addr := ":8080"
	tests := []string{"test 1", "test 2", "test 3"}
	writer := WriterStub{}
	extractor := createExtractor(EXTRACTOR_STRING)
	udp := CreateUdp(addr, extractor, &writer)
	udp.Start()

	time.Sleep(1 * time.Second)

	for _, test := range tests {
		err := udpSend(test, addr)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	time.Sleep(1 * time.Second)

	results := writer.ReadBuffer()
	founded := 0
	for _, test := range tests {
		for _, result := range results {
			if test == result["msg"].(string) {
				founded++
			}
		}
	}

	if founded != len(tests) {
		t.Fatal("Results not match")
	}

	udp.Stop()

	writer.ClearBuffer()

	time.Sleep(1 * time.Second)

	for _, test := range tests {
		err := udpSend(test, addr)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	results = writer.ReadBuffer()

	if len(results) != 0 {
		t.Fatal("Expects empty results")
	}
}


func udpSend(msg string, addr string) error {
	var err error
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(conn, msg)
	conn.Close()
	return err
}
