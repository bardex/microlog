package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type sqlite struct {
	db *sql.DB
}

func CreateSqlite(sqlitePath string) (Storage, error) {
	var err error
	var stor *sqlite
	var db *sql.DB

	db, err = sql.Open("sqlite3", "file:"+sqlitePath+"?cache=shared")

	if err == nil {
		db.SetMaxOpenConns(1)

		_, err = db.Exec("PRAGMA synchronous=NORMAL")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("PRAGMA journal_mode=WAL")
		if err != nil {
			return nil, err
		}

		stor = &sqlite{db: db}
	}

	return stor, err
}

func (s *sqlite) Write(row map[string]interface{}) error {
	var err error

	tx, _ := s.db.Begin()

	result, err := tx.Exec("INSERT INTO message (id) values (null)")

	if err != nil {
		return err
	}

	messageId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	for key, value := range row {
		fieldType := getFieldTypeByValueType(getValueType(value))
		_, err := tx.Exec("INSERT INTO message_fields (k, "+fieldType+", 'timestamp', 'message_id') values ($1, $2, CURRENT_TIMESTAMP, $3)", key, value, messageId)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()

	return err
}

func (writer *sqlite) Init() {
	initDb(false)
}

func (writer *sqlite) Close() {
	closeDb()
}

func (writer *sqlite) Clear(ageSec int64) error {
	return removeOld(ageSec)
}

func (writer *sqlite) Find(tsStart int32, tsEnd int32, page int32, limit int32, searchFields []Field, selectFields []string) ([]Message, error) {
	return find(tsStart, tsEnd, page, limit, searchFields, selectFields)
}

// first param: withoutIndexes=false
func initDb(params ...bool) {

	isWithoutIndexes := false
	if len(params) > 0 {
		isWithoutIndexes = params[0]
	}

	os.Remove(sqlitePath)
	openDb()

	var query string
	var err error

	query = `CREATE TABLE message(
      id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE TABLE message_fields(
	  k TEXT NOT NULL,
	  string TEXT,
	  integer INTEGER,
	  float REAL,
	  timestamp INT,
	  message_id INT
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	if isWithoutIndexes == false {
		initDbIndexes()
	}
}

func initDbIndexes() {
	var query string
	var err error

	query = `CREATE INDEX i_message_fields_kstring ON message_fields (k, string COLLATE NOCASE)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE INDEX i_message_fields_kinteger ON message_fields (k, integer)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE INDEX i_message_fields_kfloat ON message_fields (k, float)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE INDEX i_message_fields_message_id ON message_fields (message_id)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `CREATE INDEX i_message_fields_timestamp ON message_fields (timestamp)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
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
		sqlWhere2 = " AND mf.k IN ('" + strings.Join(selectFields, "', '") + "')"
	}

	if page < 1 {
		page = 1
	}
	skip := limit * (page - 1)

	query := "SELECT mf.k, COALESCE (mf.string, mf.integer, mf.float) AS value, mf.timestamp, mf.message_id  " +
		"FROM message_fields mf WHERE mf.message_id IN (" +
		"SELECT m1.message_id FROM message_fields m1 " + sqlJoin + sqlWhere + " " +
		"GROUP BY m1.message_id ORDER BY m1.message_id DESC LIMIT " + fmt.Sprintf("%d", skip) + ", " + fmt.Sprintf("%d", limit) +
		")" + sqlWhere2 + " ORDER BY mf.message_id DESC"
	//log.Fatal(query) // TODO debug
	rows, err := db.Query(query)

	if err != nil {
		return messages, err
	}
	defer rows.Close()

	var messageId int64
	var messageIdPrev int64
	var messageTime string

	messageFields := map[int64][]Field{}
	for rows.Next() {
		var key string
		var value string

		err := rows.Scan(&key, &value, &messageTime, &messageId)
		if err != nil {
			return messages, err
		}

		// save prev message if current field is for new message
		if messageIdPrev != messageId {
			messages = append(messages, Message{
				MessageId: messageId,
				Time:      messageTime,
				Fields:    []Field{},
			})

			messageIdPrev = messageId
		}

		messageFields[messageId] = append(messageFields[messageId], Field{
			Key:   key,
			Value: value,
		})

	}

	for index, mes := range messages {
		messages[index].Fields = messageFields[mes.MessageId]
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

func getValueType(value interface{}) string {
	valueStr := fmt.Sprintf("%v", value)
	valueType := "string"
	_, err := strconv.ParseInt(valueStr, 10, 64)
	if err == nil {
		valueType = "integer"
	} else {
		_, err := strconv.ParseFloat(valueStr, 64)
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
