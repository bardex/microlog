package storage

import (
	"os"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Entity Field
type Field struct {
	Key   string
	Value string
}

// Entity Message
type Message struct {
	MessageId int64
	Time      string
	Fields    []Field
}

// first param: withoutIndexes=false
func initDb(params ...bool) {

	isWithoutIndexes := false
	if len(params) > 0 {
		isWithoutIndexes = params[0]
	}

	os.Remove(sqlitePath)
	openDb()

	var sql string
	var err error

	sql = `CREATE TABLE message(
      id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE
	)`
	_, err = db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}

	sql = `CREATE TABLE message_fields(
	  k TEXT NOT NULL,
	  string TEXT,
	  integer INTEGER,
	  float REAL,
	  timestamp INT,
	  message_id INT
	)`
	_, err = db.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}

	if isWithoutIndexes == false {
		sql = `CREATE INDEX i_message_fields_kstring ON message_fields (k, string COLLATE NOCASE)`
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}

		sql = `CREATE INDEX i_message_fields_kinteger ON message_fields (k, integer)`
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}

		sql = `CREATE INDEX i_message_fields_kfloat ON message_fields (k, float)`
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}

		sql = `CREATE INDEX i_message_fields_message_id ON message_fields (message_id)`
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}

		sql = `CREATE INDEX i_message_fields_timestamp ON message_fields (timestamp)`
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func add(tx *sql.Tx, fields map[string]string) error {
	dbErr := openDb()
	if dbErr != nil {
		return dbErr
	}

	result, err := tx.Exec("INSERT INTO message (id) values (null)")

	if err != nil {
		return err
	}

	messageId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	for key, value := range fields {
		fieldType := getFieldTypeByValueType(getValueType(value))
		_, err := tx.Exec("INSERT INTO message_fields (k, "+fieldType+", 'timestamp', 'message_id') values ($1, $2, CURRENT_TIMESTAMP, $3)", key, value, messageId)
		if err != nil {
			return err
		}
	}

	return nil
}

func find(tsStart int32, tsEnd int32, searchFields []Field) ([]Message, error) {
	messages := []Message{}

	dbErr := openDb()
	if dbErr != nil {
		return messages, dbErr
	}

	var wheres []string
	var joins []string
	// TODO clean key and value
	var counterFields int
	for _, field := range searchFields {
		counterFields++
		c := fmt.Sprintf("%d", counterFields)

		valueType := getValueType(field.Value)
		fieldType := getFieldTypeByValueType(valueType)

		if counterFields > 1 {
			joins = append(joins, "INNER JOIN message_fields m"+c+" ON (m1.message_id = m"+c+".message_id)")
		}

		if valueType == "string" {
			wheres = append(wheres, "(m"+c+".k = \""+field.Key+"\" AND m"+c+"."+fieldType+" LIKE \""+field.Value+"\")")
		} else if valueType == "integer" || valueType == "float" {
			wheres = append(wheres, "(m"+c+".k = \""+field.Key+"\" AND m"+c+"."+fieldType+" = \""+field.Value+"\")")
		}
	}

	if tsStart > 0 {
		wheres = append(wheres, "m1.timestamp > datetime("+fmt.Sprintf("%d", tsStart)+", 'unixepoch')")
	}
	if tsEnd > 0 {
		wheres = append(wheres, "m1.timestamp <= datetime("+fmt.Sprintf("%d", tsEnd)+", 'unixepoch')")
	}

	var sqlJoin string
	if len(joins) > 0 {
		sqlJoin = strings.Join(joins, " ")
	}

	var sqlWhere string
	if len(wheres) > 0 {
		sqlWhere = " WHERE " + strings.Join(wheres, " AND ")
	}

	sql := "SELECT mf.k, COALESCE (mf.string, mf.integer, mf.float) AS value, mf.timestamp, mf.message_id  " +
		"FROM message_fields mf WHERE mf.message_id IN (" +
		"SELECT m1.message_id FROM message_fields m1 " + sqlJoin + sqlWhere + " " +
		"GROUP BY m1.message_id ORDER BY m1.message_id DESC LIMIT 100" +
		") ORDER BY mf.message_id DESC"
	//log.Fatal(sql) // TODO debug
	rows, err := db.Query(sql)

	if err != nil {
		return messages, err
	}
	defer rows.Close()

	fields := []Field{}
	var messageIdPrev int64
	var timePrev string
	for rows.Next() {
		var key string
		var value string
		var time string
		var messageId int64
		err := rows.Scan(&key, &value, &time, &messageId)
		if err != nil {
			return messages, err
		}

		// save prev message if current field is for new message
		if messageIdPrev != messageId {
			messages = append(messages, Message{
				MessageId: messageIdPrev,
				Time:      timePrev,
				Fields:    fields,
			})

			fields = []Field{}
			messageIdPrev = messageId
			timePrev = time
		}

		fields = append(fields, Field{
			Key:   key,
			Value: value,
		})
	}

	return messages, nil
}

func getValueType(value string) string {
	valueType := "string"

	_, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		valueType = "integer"
	} else {
		_, err := strconv.ParseFloat(value, 64)
		if err == nil {
			valueType = "float"
		}
	}

	return valueType
}

func getFieldTypeByValueType(valueType string) string {
	typeFieldsMap := map[string]string{
		"string":  "string",
		"integer": "integer",
		"float":   "float",
	}
	return typeFieldsMap[valueType]
}
