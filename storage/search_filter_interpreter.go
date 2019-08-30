package storage

import (
	"errors"
	"fmt"
)

type Tester func(Row) bool

type SearchFilterInterpreter struct {
	tester Tester
}

func (i *SearchFilterInterpreter) Test(row Row) (bool, error) {
	if i.tester == nil {
		return false, errors.New("You must call `Compile` method before calling `Test`")
	}
	return i.tester(row), nil
}

func (i *SearchFilterInterpreter) Compile(f SearchFilter) {
	i.tester = i.compile(f)
}

func (i *SearchFilterInterpreter) compile(f SearchFilter) Tester {

	if f.ChildsUnion == "OR" {
		if len(f.Childs) > 0 {
			return func(row Row) bool {
				for _, childFilter := range f.Childs {
					childFunc := i.compile(childFilter)
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
			return func(row Row) bool {
				for _, childFilter := range f.Childs {
					childFunc := i.compile(childFilter)
					if !childFunc(row) {
						return false
					}
				}
				return true
			}
		}
	}

	if f.Operator == "=" {
		return func(row Row) bool {
			if val, exists := row[f.Field]; exists {
				return fmt.Sprintf("%v", val) == f.Value
			} else {
				return false
			}
		}
	}

	if f.Operator == "<>" {
		return func(row Row) bool {
			if val, exists := row[f.Field]; exists {
				return fmt.Sprintf("%v", val) != f.Value
			} else {
				return false
			}
		}
	}

	if f.Operator == "EXISTS" {
		return func(row Row) bool {
			_, exists := row[f.Field]
			return exists
		}
	}

	if f.Operator == "NOT_EXISTS" {
		return func(row Row) bool {
			_, exists := row[f.Field]
			return !exists
		}
	}

	// default function
	return func(row Row) bool {
		return true
	}
}
