package listeners

import (
	"fmt"
	"net"
	"testing"
	"time"
	"microlog/storage"
	"encoding/json"
)

func TestTcp(t *testing.T) {
	addr := ":8080"
	tests := []string{
		`{"facility":"abc1", "level":1, "percents":11.5, "date":"2018-10-20"}`,
		`{"facility":"abc2", "level":1, "percents":12.5, "date":"2018-11-21"}`,
		`{"facility":"abc3", "level":1, "percents":13.5, "date":"2018-12-23"}`,
	}
	stor := storage.Storage{}
	stor.Init()
	defer stor.Close()

	extractor, _ := GetExtractor(EXTRACTOR_STRING)
	tcp := CreateTcp(addr, extractor, &stor)
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
		data := []byte(test)

		fields := make(map[string]interface{})
		json.Unmarshal(data, &fields)

		searchFields := []storage.Field{}

		for key, value := range fields {
			searchFields = append(searchFields, storage.Field{
				Key:   fmt.Sprintf("%v", key),
				Value: fmt.Sprintf("%v", value),
			})
		}

		messages, err := stor.Find(0, 2000000000, 1, 100, searchFields, []string{})

		if err != nil {
			t.Fatal(err.Error())
		}

		if len(messages) == 0 {
			t.Fatalf("Message '%s' not found", test)
		}
	}

	tcp.Stop()

	//stor.Clear(0)

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
