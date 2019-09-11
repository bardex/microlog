package listeners

import (
	"encoding/json"
	"fmt"
	"microlog/storage"
	"net"
	net_http "net/http"
	"strings"
	"testing"
	"time"
)

func TestTcp(t *testing.T) {
	testListener("tcp", ":8091", t)
}

func TestUdp(t *testing.T) {
	testListener("udp", ":8092", t)
}

func TestHttp(t *testing.T) {
	testListener("http", ":8093", t)
}

func testListener(protocol string, addr string, t *testing.T) {
	listener := CreateListenerByParams(protocol, addr, EXTRACTOR_JSON)
	listener.Start()
	defer listener.Stop()
	time.Sleep(2 * time.Second) // wait start

	tests := []string{
		`{"facility":"abc1", "level":1, "percents":11.5, "date":"2018-10-20"}`,
		`{"facility":"abc2", "level":1, "percents":12.5, "date":"2018-11-21"}`,
		`{"facility":"abc3", "level":1, "percents":13.5, "date":"2018-12-23"}`,
	}

	for _, test := range tests {
		switch protocol {
		case "tcp":
			func() {
				err := tcpSend(test, addr)
				if err != nil {
					t.Fatal(err)
				}
			}()
		case "udp":
			func() {
				err := udpSend(test, addr)
				if err != nil {
					t.Fatal(err)
				}
			}()
		case "http":
			func() {
				err := httpSend(test, addr)
				if err != nil {
					t.Fatal(err)
				}
			}()
		default:
			t.Fatalf("Sender for %s is undefined", protocol)
		}

	}

	time.Sleep(1 * time.Second) // wait until all messages will send

	s, _ := storage.GetStorage()
	s.Init()
	defer s.Close()

	qb := storage.QueryBuilder{}

	for _, str := range tests {
		test := storage.Row{}
		json.Unmarshal([]byte(str), &test)
		facility := fmt.Sprintf("%v", test["facility"])
		result, err := s.Find(qb.Equal("facility", facility), 1, 10)

		if err != nil {
			t.Fatal(err, protocol)
		}

		if len(result) != 1 {
			t.Fatalf("Expected len:1, actual:%d Protocol:%s Test:%s", len(result), protocol, str)
		}

		if result[0]["facility"] != facility {
			t.Fatalf("Expected facility:%s, actual:%s Protocol:%s Test:%s", facility, result[0]["facility"], protocol, str)
		}
	}

}

func tcpSend(msg string, addr string) error {
	var err error
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(conn, msg)
	if err != nil {
		return err
	}

	return nil
}

func udpSend(msg string, addr string) error {
	var err error
	conn, err := net.Dial("udp", addr)

	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = fmt.Fprint(conn, msg)

	return err
}

func httpSend(msg string, addr string) error {
	var err error
	client := &net_http.Client{}
	req, err := net_http.NewRequest("post", "http://127.0.0.1"+addr, strings.NewReader(msg))

	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 405 {
		return fmt.Errorf("Method not allowed. Allow methods: %s", resp.Header.Get("Allow"))
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Invalid response code: %d", resp.StatusCode)
	}

	return nil
}
