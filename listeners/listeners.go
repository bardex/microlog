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
	stor, _ := storage.CreateStorage()

	switch protocol {
	case PROTOCOL_UDP:
		return CreateUdp(addr, ext, stor)
	case PROTOCOL_TCP:
		return CreateTcp(addr, ext, stor)
	case PROTOCOL_HTTP:
		return CreateHttp(addr, ext, stor)
	}

	return nil
}
