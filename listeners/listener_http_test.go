package listeners

import (
	"microlog/storage"
	net_http "net/http"
	"strings"
	"testing"
	"time"
)

func TestHttp(t *testing.T) {
	addr := ":8082"
	url := "http://127.0.0.1" + addr

	// test data
	tests := []string{"test 1", "test 2", "test 3"}

	// http listener
	storage := storage.StorageMemory{}
	extractor, _ := GetExtractor(EXTRACTOR_STRING)
	http := CreateHttp(addr, extractor, &storage)
	http.Start()

	// wait start
	time.Sleep(2 * time.Second)

	// send http data
	for _, test := range tests {
		go httpSend(test, url, "POST", t)
	}

	time.Sleep(1 * time.Second)

	// search sended data
	for _, test := range tests {
		if !storage.Find("msg", test) {
			t.Fatalf("Message '%s' not found", test)
		}
	}

	// stop listener
	http.Stop()

	storage.Clear()
}

// send http data
func httpSend(msg string, addr string, method string, t *testing.T) {
	var err error
	client := &net_http.Client{}
	req, err := net_http.NewRequest(method, addr, strings.NewReader(msg))

	if err != nil {
		t.Fatal(err.Error())
	}

	resp, err := client.Do(req)

	if err != nil {
		t.Fatal(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode == 405 {
		t.Fatalf("Method not allowed. Allow methods: %s", resp.Header.Get("Allow"))
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Invalid response code: %d", resp.StatusCode)
	}
}
