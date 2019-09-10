package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

var testfile string = "logs"

func init() {
	go func() {
		m := &runtime.MemStats{}
		c := time.Tick(time.Second * 1)
		for t := range c {
			runtime.ReadMemStats(m)
			fmt.Printf("Alloc:%d Stack:%d [%d]\n", m.TotalAlloc, m.StackInuse, t.Unix())
		}
	}()
}

func createTestFile() {
	if _, err := os.Stat(testfile); os.IsNotExist(err) {
		f, err := os.Create(testfile)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		count := 10000000
		for i := 0; i < count; i++ {
			str := fmt.Sprintf(`{"id":%v,"price":%v.46,"level":"%v","login":"login_%v"}`, i, i, i%2, i) + fmt.Sprintln()
			f.WriteString(str)
		}
	}
}

func TestSearchByFile(b *testing.T) {
	createTestFile()

	builder := QueryBuilder{}
	query := builder.And(
		builder.Greater("id", "9.5"),
		builder.Less("id", "10.5"),
	)

	var results []Row

	// фильтр
	filter := Filter{}

	// read file
	file, err := os.Open(testfile)
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var row Row
		err := json.Unmarshal(scanner.Bytes(), &row)

		if err != nil {
			b.Fatal(err)
		}

		valid, err := filter.Test(query, row)

		if err != nil {
			b.Fatal(err)
		}
		if valid {
			results = append(results, row)
		}
	}

	if len(results) != 1 {
		b.Fatalf("Invalid length of results: %d", len(results))
	}

	if fmt.Sprintf("%v", results[0]["id"]) != "10" {
		b.Fatalf("Invalid result: %v", results[0]["id"])
	}

	fmt.Println(results[0]["id"])

}
