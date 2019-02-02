package input

import "errors"

type AbstractInput interface {
	Start()
}

func Create(protocol string, addr string) (AbstractInput, error) {
	switch protocol {
	case "udp":
		return &Udp{
			Addr: addr,
		}, nil
	}
	return nil, errors.New("Unknown protocol")
}
