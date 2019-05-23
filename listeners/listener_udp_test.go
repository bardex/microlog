package listeners

import (
	"fmt"
	"net"
	"testing"
	"time"
	"microlog/storage"
)

func TestUdp(t *testing.T) {
	addr := ":8080"
	tests := []string{"test 1", "test 2", "test 3"}
	storage := storage.StorageStub{}
	extractor, _ := GetExtractor(EXTRACTOR_STRING)
	udp := CreateUdp(addr, extractor, &storage)
	udp.Start()

	time.Sleep(1 * time.Second)

	for _, test := range tests {
		err := udpSend(test, addr)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	time.Sleep(1 * time.Second)

	for _, test := range tests {
		if !storage.Find("msg", test) {
			t.Fatalf("Message '%s' not found", test)
		}
	}

	udp.Stop()

	storage.Clear()

	time.Sleep(1 * time.Second)

	for _, test := range tests {
		err := udpSend(test, addr)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	if len(storage.ReadAll()) != 0 {
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
