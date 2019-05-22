package listeners

import (
	"testing"
)

func TestSearchBuilderEqual(t *testing.T) {
	b := SearchBuilder{}
	f := b.Equal("user", "Alice")
	expect := "user = Alice"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}

func TestSearchBuilderNotEqual(t *testing.T) {
	b := SearchBuilder{}
	f := b.NotEqual("user", "Alice")
	expect := "user <> Alice"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}

func TestSearchBuilderLike(t *testing.T) {
	b := SearchBuilder{}
	f := b.Like("user", "Alice")
	expect := "user LIKE Alice"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}

func TestSearchBuilderExists(t *testing.T) {
	b := SearchBuilder{}
	f := b.Exists("user")
	expect := "EXISTS user"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}

func TestSearchBuilderNotExists(t *testing.T) {
	b := SearchBuilder{}
	f := b.NotExists("user")
	expect := "NOT_EXISTS user"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}

func TestSearchBuilderGreater(t *testing.T) {
	b := SearchBuilder{}
	f := b.Greater("level", "0")
	expect := "level > 0"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}

func TestSearchBuilderGreaterOrEqual(t *testing.T) {
	b := SearchBuilder{}
	f := b.GreaterOrEqual("level", "0")
	expect := "level >= 0"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}

func TestSearchBuilderLess(t *testing.T) {
	b := SearchBuilder{}
	f := b.Less("level", "0")
	expect := "level < 0"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}

func TestSearchBuilderLessOrEqual(t *testing.T) {
	b := SearchBuilder{}
	f := b.LessOrEqual("level", "0")
	expect := "level <= 0"
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
			b.And(
				b.Equal("user", "Bob"),
				b.Greater("age", "18"),
			),
		),
	)

	expect := "(status = 0 AND (user = Alice OR (user = Bob AND age > 18)))"
	actual := f.String()

	if actual != expect {
		t.Fatalf("Expected: '%s' but actual: '%s'", expect, actual)
	}
}
