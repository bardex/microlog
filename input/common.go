package input

type AbstractInput interface {
	Start() bool
	Shutdown() bool
	IsActive() bool
	IsEnabled() bool
	HasError() bool
	GetError() string
	GetAddr() string
	GetProtocol() string
}

type dbRecord struct {
	protocol string
	addr     string
	enabled  int8
}

var allInputs []AbstractInput

func StartAll() {
	// якобы из БД
	records := []dbRecord{
		{
			protocol: "udp",
			addr:     ":8081",
			enabled:  1,
		},
	}

	for _, row := range records {
		var input AbstractInput
		switch row.protocol {
		case "udp":
			input = CreateUdp(row.addr, true)
		}

		if input.IsEnabled() {
			go input.Start()
		}
		allInputs = append(allInputs, input)
	}
}

func GetAllInputs() []AbstractInput {
	return allInputs
}
