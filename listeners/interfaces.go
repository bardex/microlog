package listeners

type Listener interface {
	Start()
	Stop()
	IsActive() bool
	GetAddr() string
	GetError() string
}

type Extractor interface {
	Extract([]byte) (map[string]interface{}, error)
	GetName() string
	GetDescription() string
}

type Writer interface {
	Write(map[string]interface{}) error
}