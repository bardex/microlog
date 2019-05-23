package storage

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"strings"
	"testing"
)

//
// RUN this:  go test -bench . -benchmem
//

const DEFAULT_COUNT_ROWS = 10000

func fillTable(countRows int64) {
	var messageNum int64
	tx, _ := db.Begin()
	for messageNum = 1; messageNum <= countRows; messageNum++ {
		fields := make(map[string]interface{})

		var i int8
		// Add string values
		for i = 1; i <= 4; i++ {
			key := fmt.Sprintf("key%d", i)
			value := fmt.Sprintf("value%d-%d", i, messageNum)
			fields[key] = value
		}

		// Add integer values
		for i = 5; i <= 7; i++ {
			fields[fmt.Sprintf("key%d", i)] = fmt.Sprintf("%d00%d", i, messageNum)
		}

		// Add float values
		for i = 8; i <= 10; i++ {
			fields[fmt.Sprintf("key%d", i)] = fmt.Sprintf("%d.00%d", i, messageNum)
		}

		err := add(tx, fields)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func BenchmarkSearchByStringIndexes(b *testing.B) {
	sqls := []string{
		// single index
		`CREATE INDEX i_message_fields_kstring ON message_fields (k, string COLLATE NOCASE);`,

		`CREATE INDEX i_message_fields_kstring2 ON message_fields (k, string);`,

		`CREATE INDEX i_message_fields_stringk ON message_fields (string COLLATE NOCASE, k);`,

		`CREATE INDEX i_message_fields_message_id ON message_fields (message_id);`,

		`CREATE INDEX i_message_fields_timestamp ON message_fields (timestamp);`,

		`CREATE INDEX i_message_fields_kstringm ON message_fields (k, string COLLATE NOCASE, message_id);`,

		// few indexes
		`CREATE INDEX i_message_fields_kstring ON message_fields (k, string COLLATE NOCASE);
         CREATE INDEX i_message_fields_timestamp ON message_fields (timestamp);`,

		`CREATE INDEX i_message_fields_stringk ON message_fields (string COLLATE NOCASE, k);
		 CREATE INDEX i_message_fields_message_id ON message_fields (message_id);`,

		`CREATE INDEX i_message_fields_stringk ON message_fields (string COLLATE NOCASE, k);
		 CREATE INDEX i_message_fields_message_idtimestamp ON message_fields (message_id, timestamp);`,

		`CREATE INDEX i_message_fields_kstring ON message_fields (k, string COLLATE NOCASE);
		 CREATE INDEX i_message_fields_message_id ON message_fields (message_id);`,

		`CREATE INDEX i_message_fields_kstring ON message_fields (k, string COLLATE NOCASE);
		 CREATE INDEX i_message_fields_message_id ON message_fields (message_id);
         CREATE INDEX i_message_fields_timestamp ON message_fields (timestamp);`,

		`CREATE INDEX i_message_fields_kstring ON message_fields (k, string COLLATE NOCASE);
		 CREATE INDEX i_message_fields_message_idtimestamp ON message_fields (message_id, timestamp);`,
	}

	for _, sql := range sqls {

		initDb(true)
		fillTable(DEFAULT_COUNT_ROWS)

		_, err := db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}

		b.Run(sql, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				searchFields := []Field{}
				searchFields = append(searchFields, Field{
					Key:   "key1",
					Value: "value1-5000%",
				})
				searchFields = append(searchFields, Field{
					Key:   "key2",
					Value: "value2-5000%",
				})
				searchFields = append(searchFields, Field{
					Key:   "key3",
					Value: fmt.Sprintf("value3-%d", rand.Intn(1000000)) + "%",
				})

				b.StartTimer()
				find(rand.Int31n(1000000000), 2000000000+rand.Int31n(1000000), 1, 100, searchFields, []string{})
			}
		})

		closeDb()
	}
}

func BenchmarkSearchByCountRecords(b *testing.B) {

	counts := []int64{
		10000, 100000, //1000000,
	}

	for _, count := range counts {

		initDb(true)
		fillTable(count)
		initDbIndexes()

		b.Run("Count records="+fmt.Sprintf("%d", count), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				searchFields := []Field{}
				searchFields = append(searchFields,
					Field{
						Key:   "key1",
						Value: "value1-5000%",
					},
					Field{
						Key:   "key2",
						Value: "value2-5000%",
					},
					Field{
						Key:   "key3",
						Value: fmt.Sprintf("value3-%d", rand.Intn(1000000)) + "%",
					})

				b.StartTimer()
				find(rand.Int31n(1000000000), 2000000000+rand.Int31n(1000000), 1, 100, searchFields, []string{})
			}
		})

		closeDb()
	}
}

func BenchmarkFindByStringWillcard(b *testing.B) {

	queries := []string{
		"_*", "%*%", "*%", "*",
	}

	for _, query := range queries {
		initDb()
		fillTable(DEFAULT_COUNT_ROWS)

		b.Run("Query="+query, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				searchFields := []Field{}
				searchFields = append(searchFields, Field{
					Key:   "key1",
					Value: strings.Replace(query, "*", fmt.Sprintf("value%d", rand.Intn(1000000)), 1),
				})
				searchFields = append(searchFields, Field{
					Key:   "key2",
					Value: strings.Replace(query, "*", fmt.Sprintf("value%d", rand.Intn(1000000)), 1),
				})
				searchFields = append(searchFields, Field{
					Key:   "key3",
					Value: strings.Replace(query, "*", fmt.Sprintf("value%d", rand.Intn(1000000)), 1),
				})

				b.StartTimer()
				find(rand.Int31n(1000000000), 2000000000+rand.Int31n(1000000), 1, 100, searchFields, []string{})
			}
		})

		closeDb()
	}
}

func BenchmarkSearchByCountParams(b *testing.B) {

	params := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	}

	for _, count := range params {
		initDb()
		fillTable(DEFAULT_COUNT_ROWS)

		b.Run("Count_params="+fmt.Sprintf("%d", count), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				searchFields := []Field{}
				for k := 1; k <= count; k++ {
					searchFields = append(searchFields, Field{
						Key:   fmt.Sprintf("key%d", k),
						Value: fmt.Sprintf("value%d", rand.Intn(1000000)) + "%",
					})
				}

				b.StartTimer()
				find(rand.Int31n(1000000000), 2000000000+rand.Int31n(1000000), 1, 100, searchFields, []string{})
			}
		})
		closeDb()
	}
}

func BenchmarkSearchByTypeValues(b *testing.B) {

	tables := []string{
		"string",
		"integer",
		"float",
		"string,integer,float",
	}

	for _, table := range tables {
		initDb()
		fillTable(DEFAULT_COUNT_ROWS)

		b.Run("Values_type="+table, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				searchFields := []Field{}
				if table == "string" {
					searchFields = append(searchFields,
						Field{
							Key:   "key1",
							Value: fmt.Sprintf("value%d", rand.Intn(1000000)),
						},
						Field{
							Key:   "key2",
							Value: fmt.Sprintf("value%d", rand.Intn(1000000)),
						},
						Field{
							Key:   "key3",
							Value: fmt.Sprintf("value%d", rand.Intn(1000000)),
						})
				} else if table == "integer" {
					searchFields = append(searchFields,
						Field{
							Key:   "key5",
							Value: fmt.Sprintf("%d", rand.Intn(1000000)),
						},
						Field{
							Key:   "key6",
							Value: fmt.Sprintf("%d", rand.Intn(1000000)),
						},
						Field{
							Key:   "key7",
							Value: fmt.Sprintf("%d", rand.Intn(1000000)),
						})
				} else if table == "float" {
					searchFields = append(searchFields,
						Field{
							Key:   "key8",
							Value: fmt.Sprintf("%d.01", rand.Intn(1000000)),
						},
						Field{
							Key:   "key9",
							Value: fmt.Sprintf("%d.01", rand.Intn(1000000)),
						},
						Field{
							Key:   "key10",
							Value: fmt.Sprintf("%d.01", rand.Intn(1000000)),
						})
				} else {
					searchFields = append(searchFields,
						Field{
							Key:   "key1",
							Value: fmt.Sprintf("value%d", rand.Intn(1000000)),
						},
						Field{
							Key:   "key5",
							Value: fmt.Sprintf("%d", rand.Intn(1000000)),
						},
						Field{
							Key:   "key8",
							Value: fmt.Sprintf("%d.01", rand.Intn(1000000)),
						})
				}

				b.StartTimer()
				find(rand.Int31n(1000000000), 2000000000+rand.Int31n(1000000), 1, 100, searchFields, []string{})
			}
		})

		closeDb()
	}
}
