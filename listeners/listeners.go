package listeners

import (
	"microlog/storage"
)

type Handler interface {
	Listen(*Listener) error
	Close()
}

type Listener struct {
	Protocol  string
	Addr      string
	Error     string
	Active    bool
	Handler   Handler
	Extractor Extractor
	Storage   storage.Storage
}

func (l *Listener) Start() {
	go func() {
		err := l.Handler.Listen(l)

		if err != nil {
			l.Active = false
			l.Error = err.Error()
			return
		}

	}()
}

func (l *Listener) Stop() {
	l.Handler.Close()
}

func NewListenerByParams(protocol string, addr string, extractor string) Listener {
	listener := Listener{}

	listener.Protocol = protocol
	listener.Addr = addr
	listener.Storage, _ = storage.GetStorage()
	listener.Extractor, _ = GetExtractor(extractor)

	switch protocol {
	case ProtocolUdp:
		listener.Handler = &UDPHandler{}
	case PROTOCOL_TCP:
	}

	return listener
}
