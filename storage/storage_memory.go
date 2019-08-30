package storage

type StorageMemory struct {
	buffer Rows
}

// check interface
var _ Storage = (*StorageMemory)(nil)

func (w *StorageMemory) Write(row Row) error {
	w.buffer = append(w.buffer, row)
	return nil
}

func (w *StorageMemory) Clear() {
	w.buffer = make(Rows, 0, 10)
}

func (w *StorageMemory) Find(filters SearchFilter, page int32, limit int32) (Rows, error) {
	tester := &SearchFilterInterpreter{}
	tester.Compile(filters)
	results := make(Rows, 0, 10)

	for _, item := range w.buffer {
		if valid, _ := tester.Test(item); valid {
			results = append(results, item)
		}
	}

	return results, nil
}
