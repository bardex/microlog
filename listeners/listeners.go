package listeners

import "microlog/storage"

type Listener interface {
	Start()
	Stop()
	IsActive() bool
	GetAddr() string
	GetError() string
}

func CreateListenerByParams(protocol string, addr string, extractor string) Listener {
	ext, _ := GetExtractor(extractor)
	stor := &storage.Storage{}

	switch protocol {
	case PROTOCOL_UDP:
		return CreateUdp(addr, ext, stor)
	case PROTOCOL_TCP:
		return CreateTcp(addr, ext, stor)
	}

	return nil
}

type nilListener struct{}

func (udp *nilListener) Start() {
}

func (udp *nilListener) Stop() {
}

func (udp *nilListener) IsActive() bool {
	return false
}

func (udp *nilListener) GetError() string {
	return ""
}

func (udp *nilListener) GetAddr() string {
	return ""
}
