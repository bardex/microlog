package settings

import (
	"testing"
)

func TestAdd(t *testing.T) {
	item := Input{
		Protocol: PROTOCOL_UDP,
		Addr:     ":8080",
		Enabled:  1,
	}

	repo := Inputs{}
	repo.Install()

	err := repo.Add(&item)
	if err != nil {
		t.Fatalf("%s", err)
	}

	defer closeDb()
}

func TestGetOne(t *testing.T) {
	item := Input{
		Protocol: PROTOCOL_TCP,
		Addr:     ":8081",
		Enabled:  1,
	}

	repo := Inputs{}
	repo.Install()
	repo.Add(&item)

	newItem, err := repo.GetOne(1)
	t.Log(newItem)
	if err != nil {
		t.Fatalf("%s", err)
	}

	if newItem != item {
		t.Fatal("Not equal items", newItem, item)
	}

	defer closeDb()
}

func TestGetAll(t *testing.T) {
	item1 := Input{
		Protocol: PROTOCOL_UDP,
		Addr:     ":8080",
		Enabled:  1,
	}

	item2 := Input{
		Protocol: PROTOCOL_TCP,
		Addr:     ":8081",
		Enabled:  0,
	}

	repo := Inputs{}
	repo.Install()
	repo.Add(&item1)
	repo.Add(&item2)

	items, err := repo.GetAll()
	t.Log(items)
	if err != nil {
		t.Fatalf("%s", err)
	}

	if len(items) != 2 {
		t.Fatal("Count items is wrong", len(items))
	}

	defer closeDb()
}

func TestUpdate(t *testing.T) {
	item := Input{
		Protocol: PROTOCOL_UDP,
		Addr:     ":8080",
		Enabled:  1,
	}

	repo := Inputs{}
	repo.Install()
	repo.Add(&item)

	item.Protocol = PROTOCOL_TCP
	item.Addr = ":8081"
	item.Enabled = 0

	repo.Update(&item)

	newItem, err := repo.GetOne(item.Id)
	t.Log(newItem)
	if err != nil {
		t.Fatalf("%s", err)
	}

	if newItem != item {
		t.Fatal("Not equal items", newItem, item)
	}

	defer closeDb()

}

func TestDelete(t *testing.T) {
	item := Input{
		Protocol: PROTOCOL_UDP,
		Addr:     ":8080",
		Enabled:  1,
	}

	repo := Inputs{}
	repo.Install()
	repo.Add(&item)

	err := repo.Delete(item.Id)
	if err != nil {
		t.Fatalf("%s", err)
	}
	delItem, _ := repo.GetOne(item.Id)
	if delItem == item {
		t.Fatal("Equal source item and delete item", item, delItem)
	}

	defer closeDb()
}
