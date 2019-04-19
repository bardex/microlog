package listeners

type Writer interface {
	Write(map[string]interface{}) error
}

type WriterStub struct {
	buffer []map[string]interface{}
}

func (writer *WriterStub) Write(row map[string]interface{}) error {
	writer.buffer = append(writer.buffer, row)
	return nil
}

func (writer *WriterStub) ReadBuffer() []map[string]interface{} {
	return writer.buffer
}

func (writer *WriterStub) ClearBuffer() {
	writer.buffer = make([]map[string]interface{}, 0, 10)
}

