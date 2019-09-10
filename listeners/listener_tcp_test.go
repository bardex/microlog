package listeners

import (
	"encoding/json"
	"fmt"
	"microlog/storage"
	"net"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestTcp(t *testing.T) {
	addr := ":8080"
	tests := []string{
		`{"facility":"abc1", "level":1, "percents":11.5, "date":"2018-10-20"}`,
		`{"facility":"abc2", "level":1, "percents":12.5, "date":"2018-11-21"}`,
		`{"facility":"abc3", "level":1, "percents":13.5, "date":"2018-12-23"}`,
	}

	tcp := CreateListenerByParams("tcp", addr, EXTRACTOR_JSON)
	tcp.Start()

	time.Sleep(2 * time.Second)

	for _, test := range tests {
		wg.Add(1)
		go tcpSend(test, addr, t)
	}

	wg.Wait()

	s, _ := storage.GetStorage()

	qb := storage.QueryBuilder{}

	for _, str := range tests {
		test := storage.Row{}
		json.Unmarshal([]byte(str), &test)
		facility := fmt.Sprintf("%v", test["facility"])
		result, err := s.Find(qb.Equal("facility", facility), 1, 10)

		if err != nil {
			t.Fatal(err)
		}

		if len(result) != 1 {
			t.Fatalf("Expected len:1, actual:%d", len(result))
		}

		if result[0]["facility"] != facility {
			t.Fatalf("Expected facility:%s, actual:%s", facility, result[0]["facility"])
		}
	}

	tcp.Stop()

	conn, err := net.Dial("tcp", addr)

	if conn != nil || err == nil {
		t.Fatal("Unexpected connection")
	}
}

func tcpSend(msg string, addr string, t *testing.T) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer conn.Close()

	_, err = fmt.Fprint(conn, msg)

	if err != nil {
		t.Fatal(err.Error())
	}

	wg.Done()
}
