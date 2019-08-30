package listeners

import (
	"io/ioutil"
	"microlog/storage"
	"net"
	net_http "net/http"
)

const PROTOCOL_HTTP = "http"

type http struct {
	addr      string
	error     string
	active    bool
	conn      net.Listener
	extractor Extractor
	storage   storage.Storage
}

func CreateHttp(addr string, extractor Extractor, storage storage.Storage) Listener {
	return &http{
		addr:      addr,
		active:    false,
		extractor: extractor,
		storage:   storage,
	}
}

// start listen
func (http *http) Start() {
	go (func() {

		net_http.HandleFunc("/", http.handleConn)
		err := net_http.ListenAndServe(http.addr, nil)

		if err != nil {
			http.error = err.Error()
			return
		}

		http.active = true
		http.error = ""
	})()
}

func (http *http) handleConn(w net_http.ResponseWriter, r *net_http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.error = err.Error()
		}

		row, err := http.extractor.Extract(body)
		row["remote_addr"] = r.RemoteAddr

		if err == nil {
			http.storage.Write(row)
		} else {
			http.error = err.Error()
		}
	} else {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
	}
}

// stop listen
func (http *http) Stop() {
	if http.active {
		if http.conn != nil {
			_ = http.conn.Close()
		}
	}
	http.active = false
	http.conn = nil
}

func (http *http) IsActive() bool {
	return http.active
}

func (http *http) GetError() string {
	return http.error
}

func (http *http) GetAddr() string {
	return http.addr
}
