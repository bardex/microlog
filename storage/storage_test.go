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

	messages, err := find(0, 2000000000, fields)
	t.Log(err)
	t.Log(messages)
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