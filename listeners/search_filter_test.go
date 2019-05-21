package listeners

import (
	"testing"
)

func TestSearchBuilderEqual(t *testing.T) {
	b := SearchBuilder{}
	f := b.Equal("user", "Alice")
	expect := "user=Alice"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}

func TestSearchBuilderNesting(t *testing.T) {
	b := SearchBuilder{}

	f := b.And(
		b.Equal("status", "0"),
		b.Or(
			b.Equal("user", "Alice"),
			b.Equal("user", "Bob"),
		),
	)

	expect := "(status=0 AND (user=Alice OR user=Bob))"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}
