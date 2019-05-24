package storage

import "fmt"

type StorageStub struct {
	buffer []map[string]interface{}
}

func (writer *StorageStub) Write(row map[string]interface{}) error {
	writer.buffer = append(writer.buffer, row)
	fmt.Printf("%#v", row)
	return nil
}

func (writer *StorageStub) ReadAll() []map[string]interface{} {
	return writer.buffer
}

func (writer *StorageStub) Clear() {
	writer.buffer = make([]map[string]interface{}, 0, 10)
}

func (writer *StorageStub) Find(key string, value string) bool {
	for _, item := range writer.buffer {
		if value == item[key].(string) {
			return true
		}
	}
	return false
}

func (writer *StorageStub) Search(SearchFilter) ([]Message, error) {
	return []Message{}, nil
}
