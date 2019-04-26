package listeners

type Listener interface {
	Start()
	Stop()
	IsActive() bool
	GetAddr() string
	GetError() string
}

func CreateListenerByParams(protocol string, addr string, extractor string) Listener {
	ext, _ := GetExtractor(extractor)
	writer := &WriterStub{}

	switch protocol {
	case PROTOCOL_UDP:
		return CreateUdp(addr, ext, writer)
	}

	return &nilListener{}
}

type nilListener struct {}

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

