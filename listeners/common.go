package listeners

type Listener interface {
	Start()
	Stop()
	IsActive() bool
	GetAddr() string
	GetError() string
}
