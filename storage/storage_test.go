package storage

import (
	"fmt"
	"log"
	"testing"
)

func TestAdd(t *testing.T) {
	initDb(false)
	defer closeDb()
	fillDb(1000)
}

func TestFind(t *testing.T) {
	initDb(false)
	defer closeDb()
	fillDb(10000)

	fields := []Field{}
	fields = append(fields, Field{
		Key:   "key1",
		Value: "value1-211%",
	})

	messages, err := find(0, 2000000000, 1, 100, fields, []string{})
	if err != nil {
		log.Fatal(err)
	}
	t.Log(messages)
	if len(messages) == 0 {
		t.Fatal("Count items = 0")
	}

	messages, err = find(0, 2000000000, 1, 2, fields, []string{"key1", "key2"})
	if err != nil {
		log.Fatal(err)
	}
	t.Log(messages)
	if len(messages) != 2 {
		t.Fatal("Count items != 2")
	}
}

func TestRemoveOld(t *testing.T) {
	initDb(false)
	defer closeDb()
	fillDb(1000)

	errRemove := removeOld(0)
	if errRemove != nil {
		log.Fatal(errRemove)
	}

	messages, err := find(0, 2000000000, 1, 10, []Field{}, []string{})
	if err != nil {
		log.Fatal(err)
	}
	if len(messages) > 0 {
		t.Fatal("Count items > 0")
	}
}

func TestCleanStr(t *testing.T)  {
	str := "123 abc def-.%<=>'\""
	if cleanStr(str) != "123 abc def-.%<=>" {
		t.Fatal("Not cleaned string")
	}
}

func fillDb(countRows int64) {
	var messageNum int64
	tx, _ := db.Begin()
	for messageNum = 1; messageNum <= countRows; messageNum++ {

		fields := make(map[string]string)

		var i int8
		// Add string values
		for i = 1; i <= 4; i++ {
			key := fmt.Sprintf("key%d", i)
			value := fmt.Sprintf("value%d-%d", i, messageNum)
			fields[key] = value
		}

		// Add integer values
		for i = 5; i <= 7; i++ {
			fields[fmt.Sprintf("key%d", i)] = fmt.Sprintf("%d00%d", i, messageNum)
		}

		// Add float values
		for i = 8; i <= 10; i++ {
			fields[fmt.Sprintf("key%d", i)] = fmt.Sprintf("%d.00%d", i, messageNum)
		}

		err := add(tx, fields)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}