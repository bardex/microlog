package storage

import (
	"errors"
	"log"
)

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

func (w *StorageMemory) Find(query Query, page int32, limit int32) (Rows, error) {
	if page <= 0 {
		return Rows{}, errors.New("Param `page` must be greater 0")
	}

	if limit <= 0 {
		return Rows{}, errors.New("Param `limit` must be greater 0")
	}

	filter := Filter{}
	results := make(Rows, 0, limit)
	offset := (page - 1) * limit
	counter := int32(0)

	for _, item := range w.buffer {
		valid, err := filter.Test(query, item)

		if err != nil {
			log.Println(err)
			continue
		}

		if valid {
			counter++
			if counter > offset {
				results = append(results, item)
				if counter-offset >= limit {
					break
				}
			}
		}
	}

	return results, nil
}

func (w *StorageMemory) Init() error {
	return nil
}

func (w *StorageMemory) Close() error {
	w.buffer = make(Rows, 0, 10)
	return nil
}
