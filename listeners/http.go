package listeners

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

const ProtocolHttp = "http"

type HTTPHandler struct {
	conn     net.Listener
	listener *Listener
}

func (h *HTTPHandler) Listen(listener *Listener) error {
	h.listener = listener
	err := http.ListenAndServe(listener.Addr, h)
	if err != nil {
		return err
	}
	return nil
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.listener.Error = err.Error()
		}

		row, err := h.listener.Extractor.Extract(body)
		row["remote_addr"] = r.RemoteAddr

		if err == nil {
			h.listener.Storage.Write(row)
		} else {
			h.listener.Error = err.Error()
		}
	} else {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
	}
}

// stop listen
func (h *HTTPHandler) Close() {
	if h.conn != nil {
		h.conn.Close()
	}
	h.conn = nil
}
