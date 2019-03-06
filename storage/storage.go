package storage

import (
	"os"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"regexp"
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

func find(tsStart int32, tsEnd int32, page int32, limit int32, searchFields []Field, selectFields []string) ([]Message, error) {
	messages := []Message{}

	dbErr := openDb()
	if dbErr != nil {
		return messages, dbErr
	}

	var wheres []string
	var joins []string

	var counterFields int
	for _, field := range searchFields {
		key := cleanStr(field.Key)
		value := cleanStr(field.Value)

		counterFields++
		c := fmt.Sprintf("%d", counterFields)

		valueType := getValueType(value)
		fieldType := getFieldTypeByValueType(valueType)

		if counterFields > 1 {
			joins = append(joins, "INNER JOIN message_fields m"+c+" ON (m1.message_id = m"+c+".message_id)")
		}

		if valueType == "string" {
			wheres = append(wheres, "(m"+c+".k = \""+key+"\" AND m"+c+"."+fieldType+" LIKE \""+value+"\")")
		} else if valueType == "integer" || valueType == "float" {
			wheres = append(wheres, "(m"+c+".k = \""+key+"\" AND m"+c+"."+fieldType+" = \""+value+"\")")
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

	var sqlWhere2 string
	if len(selectFields) > 0 {
		sqlWhere2 = " AND mf.k IN ('"+strings.Join(selectFields, "', '")+"')"
	}

	if page < 1 {
		page = 1
	}
	skip := limit * (page - 1)

	sql := "SELECT mf.k, COALESCE (mf.string, mf.integer, mf.float) AS value, mf.timestamp, mf.message_id  " +
		"FROM message_fields mf WHERE mf.message_id IN (" +
		"SELECT m1.message_id FROM message_fields m1 " + sqlJoin + sqlWhere + " " +
		"GROUP BY m1.message_id ORDER BY m1.message_id DESC LIMIT " + fmt.Sprintf("%d", skip) + ", " + fmt.Sprintf("%d", limit) +
		")"+sqlWhere2+" ORDER BY mf.message_id DESC"
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

func cleanStr(str string) string {
	var re = regexp.MustCompile(`[^A-z% 0-9.<=>_-]*`)
	return re.ReplaceAllString(str, "")
}

func removeOld(ageSec int64) error {
	dbErr := openDb()
	if dbErr != nil {
		return dbErr
	}
	maxTimestamp := time.Now().Unix() - ageSec
	_, err := db.Exec("DELETE FROM message_fields WHERE timestamp <= datetime($1, 'unixepoch')", fmt.Sprintf("%d", maxTimestamp))
	return err
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
