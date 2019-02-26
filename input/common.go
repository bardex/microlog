package input

const SIGNAL_STOP = 1

type AbstractInput interface {
	Start()
	Shutdown()
	IsActive() bool
	IsEnabled() bool
	HasError() bool
	GetError() string
	GetAddr() string
	GetProtocol() string
	GetId() int
}

type dbRecord struct {
	id       int
	protocol string
	addr     string
	enabled  int8
}

var allInputs []AbstractInput



func GetAllInputs() []AbstractInput {
	return allInputs
}

func GetById(id int) AbstractInput {
	for _, input := range GetAllInputs() {
		if input.GetId() == id {
			return input
		}
	}
	return nil
}
