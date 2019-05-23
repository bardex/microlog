package listeners

import "fmt"

type StorageMemory struct {
	buffer []map[string]interface{}
}

func (writer *StorageMemory) Write(row map[string]interface{}) error {
	writer.buffer = append(writer.buffer, row)
	return nil
}

func (writer *StorageMemory) Search(f SearchFilter) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, 10)
	filter := writer.compileFilter(f)

	for _, row := range writer.buffer {
		if filter(row) {
			results = append(results, row)
		}
	}

	return results
}

func (writer *StorageMemory) compileFilter(f SearchFilter) func(map[string]interface{}) bool {

	if f.ChildsUnion == "OR" {
		if len(f.Childs) > 0 {
			return func(row map[string]interface{}) bool {
				for _, childFilter := range f.Childs {
					childFunc := writer.compileFilter(childFilter)
					if childFunc(row) {
						return true
					}
				}
				return false
			}
		}
	}

	if f.ChildsUnion == "AND" {
		if len(f.Childs) > 0 {
			return func(row map[string]interface{}) bool {
				for _, childFilter := range f.Childs {
					childFunc := writer.compileFilter(childFilter)
					if !childFunc(row) {
						return false
					}
				}
				return true
			}
		}
	}

	if f.Operator == "=" {
		return func(row map[string]interface{}) bool {
			if val, exists := row[f.Field]; exists {
				return fmt.Sprintf("%v", val) == f.Value
			} else {
				return false
			}
		}
	}


	// default function
	return func(row map[string]interface{}) bool {
		return true
	}

}

func (writer *StorageMemory) ReadBuffer() []map[string]interface{} {
	return writer.buffer
}

func (writer *StorageMemory) Clear() {
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
