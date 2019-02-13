package input

import "time"

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

func StartAll() {
	// якобы из БД
	records := []dbRecord{
		{
			id:       1,
			protocol: "udp",
			addr:     ":8081",
			enabled:  1,
		},
	}

	for _, row := range records {
		var input AbstractInput
		switch row.protocol {
		case "udp":
			input = CreateUdp(row.id, row.addr, true)
		}

		if input.IsEnabled() {
			go input.Start()
		}

		time.Sleep(2 * time.Second)
		//input.Shutdown();

		allInputs = append(allInputs, input)
	}
}

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
