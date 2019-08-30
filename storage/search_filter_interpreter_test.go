package storage

import (
	"testing"
)

type TestFixture struct {
	row    Row
	answer bool
}

func TestSearchEqual(t *testing.T) {
	b := SearchBuilder{}
	filter := b.Equal("user", "Alice")

	fixtures := []TestFixture{
		{row: Row{"user": "Alice"}, answer: true},
		{row: Row{"user": "Bob"}, answer: false},
		{row: Row{"login": "Alice"}, answer: false},
	}

	testFilter(t, filter, fixtures)
}

func TestSearchNotEqual(t *testing.T) {
	b := SearchBuilder{}
	filter := b.NotEqual("user", "Alice")

	fixtures := []TestFixture{
		{row: Row{"user": "Alice"}, answer: false},
		{row: Row{"user": "Bob"}, answer: true},
		{row: Row{"login": "Alice"}, answer: false},
	}

	testFilter(t, filter, fixtures)
}

func TestSearchExists(t *testing.T) {
	b := SearchBuilder{}
	filter := b.Exists("user")

	fixtures := []TestFixture{
		{row: Row{"user": "Alice"}, answer: true},
		{row: Row{"login": "Alice"}, answer: false},
	}

	testFilter(t, filter, fixtures)
}

func TestSearchNotExists(t *testing.T) {
	b := SearchBuilder{}
	filter := b.NotExists("user")

	fixtures := []TestFixture{
		{row: Row{"user": "Alice"}, answer: false},
		{row: Row{"login": "Alice"}, answer: true},
	}

	testFilter(t, filter, fixtures)
}

func TestSearchInterpreterAnd(t *testing.T) {
	b := SearchBuilder{}
	filter := b.And(
		b.Equal("user", "Alice"),
		b.Equal("level", "10"),
	)

	fixtures := []TestFixture{
		{row: Row{"user": "Alice", "level": "10"}, answer: true},
		{row: Row{"user": "Alice", "level": "10", "login": "Bob"}, answer: true},
		{row: Row{"user": "Alice", "level": "11"}, answer: false},
		{row: Row{"level": "10"}, answer: false},
		{row: Row{"login": "Alice"}, answer: false},
	}

	testFilter(t, filter, fixtures)
}

func testFilter(t *testing.T, filter SearchFilter, fixtures []TestFixture) {
	checker := &SearchFilterInterpreter{}
	checker.Compile(filter)

	for _, fix := range fixtures {
		answer, _ := checker.Test(fix.row)
		if answer != fix.answer {
			t.Fatalf("Fail row: #%v filter:#%v", fix.row, filter)
		}
	}
}
