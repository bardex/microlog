package settings

import (
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	dbname := "test_1.db"
	sqlitePath = dbname
	os.Remove(dbname)
	defer os.Remove(dbname)

	Migrate()

	item := Input{
		Protocol:  "udp",
		Extractor: "JSON",
		Addr:      ":8080",
		Enabled:   1,
	}

	repo := Inputs

	err := repo.Add(&item)
	if err != nil {
		t.Fatalf("%s", err)
	}

	defer closeDb()
}

func TestGetOne(t *testing.T) {
	dbname := "test_2.db"
	sqlitePath = dbname
	os.Remove(dbname)
	defer os.Remove(dbname)

	Migrate()

	item := &Input{
		Protocol:  "udp",
		Extractor: "JSON",
		Addr:      ":8081",
		Enabled:   1,
	}

	Inputs.Add(item)

	newItem, err := Inputs.GetOne(1)
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
	dbname := "test_3.db"
	sqlitePath = dbname
	os.Remove(dbname)
	defer os.Remove(dbname)

	Migrate()

	item1 := Input{
		Protocol:  "udp",
		Extractor: "JSON",
		Addr:      ":8080",
		Enabled:   1,
	}

	item2 := Input{
		Protocol:  "udp",
		Extractor: "JSON",
		Addr:      ":8081",
		Enabled:   0,
	}

	Inputs.Add(&item1)
	Inputs.Add(&item2)

	items, err := Inputs.GetAll()
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
	dbname := "test_4.db"
	sqlitePath = dbname
	os.Remove(dbname)
	defer os.Remove(dbname)

	Migrate()

	item := &Input{
		Protocol:  "udp",
		Extractor: "JSON",
		Addr:      ":8080",
		Enabled:   1,
	}

	Inputs.Add(item)

	item.Protocol = "udp"
	item.Addr = ":8081"
	item.Extractor = "JSON"
	item.Enabled = 0

	Inputs.Update(item)

	newItem, err := Inputs.GetOne(item.Id)
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
	dbname := "test_5.db"
	sqlitePath = dbname
	os.Remove(dbname)
	defer os.Remove(dbname)

	Migrate()

	item := &Input{
		Protocol: "udp",
		Addr:     ":8080",
		Enabled:  1,
	}

	Inputs.Add(item)

	err := Inputs.Delete(item.Id)
	if err != nil {
		t.Fatalf("%s", err)
	}
	delItem, _ := Inputs.GetOne(item.Id)
	if delItem == item {
		t.Fatal("Equal source item and delete item", item, delItem)
	}

	defer closeDb()
}
