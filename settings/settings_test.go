package settings

import (
	"microlog/settings/inputs"
	"testing"
)

func TestAll(t *testing.T) {

	errCon := inputs.Connect()
	if errCon != nil {
		t.Fatalf("%s", errCon)
	}

	t.Log(inputs.Install())

	t.Log(inputs.Add(inputs.Record{
		Protocol: "udp",
		Addr:     ":8081",
		Enabled:  1,
	}))
	t.Log(inputs.Add(inputs.Record{
		Protocol: "udp",
		Addr:     ":8082",
		Enabled:  0,
	}))

	for i := 1; i < 10; i++ {
		inputs.Add(inputs.Record{
			Protocol: "udp",
			Addr:     ":8083",
			Enabled:  1,
		})
	}

	t.Log(inputs.Update(inputs.Record{
		Id:       1,
		Protocol: "tcp",
		Addr:     "127.0.0.1:8080",
		Enabled:  0,
	}))
	records, err := inputs.GetAll()
	t.Log(err)
	for _, r := range records {
		t.Log(r.Id, r.Protocol, r.Addr, r.Enabled)
	}

	t.Log(inputs.GetOne(1))
	t.Log(inputs.Delete(2))

	t.Log(inputs.Disconnect())
}
