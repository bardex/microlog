package storage

import (
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	var b QueryBuilder
	var filter Filter

	tester := func(q Query, rtrue Message, rfalse Message) {
		yes, err := filter.Test(q, rtrue)
		if err != nil {
			t.Fatal(err)
		}
		if yes == false {
			t.Fatalf("`True` Test failed %#v, %#v", q, rtrue)
		}
		no, err := filter.Test(q, rfalse)
		if err != nil {
			t.Fatal(err)
		}
		if no == true {
			t.Fatalf("`False` Test failed %#v, %#v", q, rfalse)
		}
	}

	// tester(query, true row, false row)
	tester(b.Exists("id"), Message{"id": 20}, Message{"level": 1})
	tester(b.NotExists("id"), Message{"level": 20}, Message{"id": 1})

	tester(b.Equal("id", "10"), Message{"id": 10}, Message{"id": 20})
	tester(b.Equal("id", "10"), Message{"id": "10"}, Message{"id": "20"})
	tester(b.Equal("id", "10"), Message{"id": "10"}, Message{"level": "20"})

	tester(b.Like("login", "admin"), Message{"login": "_t_admin_10"}, Message{"login": "adm"})
	tester(b.Like("login", "admin"), Message{"login": "_t_admin_10"}, Message{"name": "admin"})
	tester(b.Like("login", "admin"), Message{"login": "_t_admin_10"}, Message{"login": "Admin"})

	tester(b.NotEqual("id", "10"), Message{"id": 20}, Message{"id": 10})
	tester(b.NotEqual("id", "10"), Message{"id": "20"}, Message{"id": "10"})
	tester(b.NotEqual("id", "10"), Message{"level": "10"}, Message{"id": "10"})

	tester(b.Less("id", "10"), Message{"id": 9}, Message{"id": 11})
	tester(b.Less("id", "10"), Message{"id": "9"}, Message{"id": "11"})
	tester(b.Less("id", "10"), Message{"id": "9.5"}, Message{"id": "10.5"})
	tester(b.Less("id", "10"), Message{"id": 9.5}, Message{"id": 10.5})
	tester(b.Less("id", "10"), Message{"id": 9.5}, Message{"level": 10.5})

	tester(b.LessOrEqual("id", "10"), Message{"id": 9}, Message{"id": 11})
	tester(b.LessOrEqual("id", "10"), Message{"id": "10"}, Message{"id": "11"})
	tester(b.LessOrEqual("id", "10"), Message{"id": "9.5"}, Message{"id": "10.5"})
	tester(b.LessOrEqual("id", "10"), Message{"id": 10}, Message{"id": 10.5})
	tester(b.LessOrEqual("id", "10"), Message{"id": 10}, Message{"level": 10.5})

	tester(b.Greater("id", "10"), Message{"id": 11}, Message{"id": 9})
	tester(b.Greater("id", "10"), Message{"id": "11"}, Message{"id": "9"})
	tester(b.Greater("id", "10"), Message{"id": "10.5"}, Message{"id": "9.5"})
	tester(b.Greater("id", "10"), Message{"id": 10.5}, Message{"id": 9.5})
	tester(b.Greater("id", "10"), Message{"id": 10.5}, Message{"level": 9.5})

	tester(b.GreaterOrEqual("id", "11"), Message{"id": 11}, Message{"id": 9})
	tester(b.GreaterOrEqual("id", "11"), Message{"id": "11.5"}, Message{"id": "10.5"})
	tester(b.GreaterOrEqual("id", "10.5"), Message{"id": "10.5"}, Message{"id": 9.5})
	tester(b.GreaterOrEqual("id", "10.5"), Message{"id": 10.5}, Message{"id": "9.4"})
	tester(b.GreaterOrEqual("id", "10.5"), Message{"id": 10.5}, Message{"level": "9.4"})

	tester(
		b.Or(b.Equal("id", "10"), b.Equal("status", "0")),
		Message{"id": "10"},
		Message{"status": 1},
	)

	tester(
		b.And(b.Equal("id", "10"), b.Equal("status", "0")),
		Message{"id": 10, "status": "0"},
		Message{"id": 10, "status": 1},
	)

}

func TestInterToInt64(t *testing.T) {
	var data []interface{}

	data = append(data, int(10))
	data = append(data, int8(10))
	data = append(data, int16(10))
	data = append(data, int32(10))
	data = append(data, int64(10))
	data = append(data, uint(10))
	data = append(data, uint8(10))
	data = append(data, uint16(10))
	data = append(data, uint32(10))
	data = append(data, uint64(10))
	data = append(data, "10")

	var actual int64
	var expected int64 = 10
	var err error

	for _, v := range data {
		actual, err = interToInt64(v)
		if err != nil {
			t.Fatal(err)
		}
		if actual != expected {
			t.Fatalf("Actual %v not equal expected %v", actual, expected)
		}
	}

	_, err = interToInt64(float32(10))
	if err == nil {
		t.Fatal("Expected error")
	}

	_, err = interToInt64("aaghf")
	if err == nil {
		t.Fatal("Expected error")
	}

}

func TestInterToFloat64(t *testing.T) {
	var data []interface{}

	data = append(data, float32(10.15))
	data = append(data, float64(10.15))
	data = append(data, "10.15")

	var actual float64
	var expected float64 = 10.15
	var err error

	for _, v := range data {
		actual, err = interToFloat64(v)
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%.2f", actual) != fmt.Sprintf("%.2f", expected) {
			t.Fatalf("Actual %v not equal expected %v", actual, expected)
		}
	}

	_, err = interToFloat64("aaghf")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestEqualBench(b *testing.T) {
	qb := QueryBuilder{}
	q := qb.Equal("id", "10")
	f := Filter{}
	count := 10000000

	for i := 0; i < count; i++ {
		row := Message{"id": i}
		f.Test(q, row)
	}

}
