package listeners

type StorageMemory struct {
	buffer []map[string]interface{}
}

func (writer *StorageMemory) Write(row map[string]interface{}) error {
	writer.buffer = append(writer.buffer, row)
	return nil
}

func (writer *StorageMemory) ReadBuffer() []map[string]interface{} {
	return writer.buffer
}

func (writer *StorageMemory) ClearBuffer() {
	writer.buffer = make([]map[string]interface{}, 0, 10)
}

func (writer *StorageMemory) Find(key string, value string) bool {
	for _, item := range writer.buffer {
		if value == item[key].(string) {
			return true
		}
	}
	return false
}
