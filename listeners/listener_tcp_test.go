package listeners

import (
	"fmt"
	"net"
	"testing"
	"time"
	"microlog/storage"
)

func TestTcp(t *testing.T) {
	addr := ":8088"
	tests := []string{"test 1", "test 2", "test 3"}
	storage := storage.StorageStub{}
	extractor, _ := GetExtractor(EXTRACTOR_STRING)
	tcp := CreateTcp(addr, extractor, &storage)
	tcp.Start()

	time.Sleep(2 * time.Second)

	for _, test := range tests {
		err := tcpSend(test, addr)
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

	tcp.Stop()

	storage.Clear()

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", addr)

	if conn != nil || err == nil {
		t.Fatal("Unexpected connection")
	}
}

func tcpSend(msg string, addr string) error {
	var err error
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(conn, msg)
	conn.Close()
	return err
}
