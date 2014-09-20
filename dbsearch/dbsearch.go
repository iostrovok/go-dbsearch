package dbsearch

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/iostrovok/go-dbsearch/dbsearch/sqler"
	"github.com/iostrovok/go-iutils/iutils"
	_ "github.com/lib/pq"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
)

//
type OneRow struct {
	Name    string
	DBName  string
	Type    string
	IsArray bool
}

type AllRows struct {
	DBList map[string]*OneRow
	List   map[string]*OneRow
}

type Searcher struct {
	db *sql.DB
}

func (s *Searcher) Ping() error {

	if s.db == nil {
		return fmt.Errorf("can't connect to DB")
	}

	if err := s.db.Ping(); err != nil {
		return err
	}

	return nil
}

func DBI(poolSize int, dsn string) (*Searcher, error) {

	s := new(Searcher)

	db, _ := sql.Open("postgres", dsn)

	if err := db.Ping(); err != nil {
		log.Fatalf("DB Error: %s\n", err)
	} else {
		s.db = db
	}

	s.db.SetMaxOpenConns(poolSize)

	return s, nil
}

func (s *Searcher) GetCount(sqlLine string, values []interface{}) (int, error) {
	var count int
	err := s.db.QueryRow(sqlLine, values...).Scan(&count)
	return count, err
}

func (s *Searcher) GetOne(mType *AllRows, sqlLine string, values []interface{}) (map[string]interface{}, error) {

	sqlLine += " LIMIT 1 OFFSET 0 "

	list, err := s.Get(mType, sqlLine, values)
	empty := map[string]interface{}{}

	if err != nil {
		return empty, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return empty, nil
}

func (s *Searcher) Get(mType *AllRows, sqlLine string, values []interface{}) ([]map[string]interface{}, error) {

	Out := make([]map[string]interface{}, 0)

	log.Printf("dbsearch.Get: %s\n", sqlLine)
	spew.Dump(values)

	rows, err := s.db.Query(sqlLine, values...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return Out, err
	}

	//log.Println("sqlLine: " + sqlLine + "\nCOLS:")

	rawResult := make([]interface{}, 0)
	for i := 0; i < len(cols); i++ {
		//log.Printf("search for %d => %s", i, cols[i])
		t, find := mType.DBList[cols[i]]
		if !find {
			log.Fatalf("dbsearch.Get not found column: %s!", cols[i])
		}

		switch t.Type {
		case "datetime", "date":
			datetime := new(*time.Time)
			rawResult = append(rawResult, datetime)
		case "int", "numeric":
			rawResult = append(rawResult, new(int))
		case "bigint":
			rawResult = append(rawResult, new(int64))
		default:
			rawResult = append(rawResult, make([]byte, 0))
		}
	}

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			log.Fatal(err)
		}

		result := map[string]interface{}{}
		for i, raw := range rawResult {
			// cols[i] - Column name
			if raw == nil {
				result[cols[i]] = nil
			} else {
				result[cols[i]] = convertType(cols[i], mType, raw)
			}
		}

		Out = append(Out, result)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return Out, nil
}

func convertType(Name string, mType *AllRows, raw_in interface{}) interface{} {
	t, find := mType.DBList[Name]
	if !find {
		log.Fatal("Not found!")
	}

	switch t.Type {
	case "text":
		raw := raw_in.([]byte)
		if t.IsArray {
			return parseArray(string(raw))
		} else {
			return string(raw)
		}
	case "json":
		line := iutils.AnyToString(raw_in)
		if line == "" {
			line = "{}"
		}
		raw := []byte(line)
		var res map[string]interface{}
		err := json.Unmarshal(raw, &res)
		if err != nil {
			log.Fatal("error:", err)
		}
		return res
	case "bigint", "int64", "int":
		return iutils.AnyToInt(raw_in)
	case "date", "datetime":
		return raw_in
	}
	return nil
}

func Prepare(s interface{}) *AllRows {
	st := reflect.TypeOf(s)

	aRows := AllRows{}
	aRows.List = make(map[string]*OneRow, 0)
	aRows.DBList = make(map[string]*OneRow, 0)
	Count := 0
	for true {
		field := st.Field(Count)

		if field.Name == "" {
			break
		}

		Count++

		dbname := field.Tag.Get("db")
		oRow := OneRow{
			Name:    field.Name,
			DBName:  dbname,
			Type:    field.Tag.Get("type"),
			IsArray: false,
		}
		if field.Tag.Get("is_array") == "yes" {
			oRow.IsArray = true
		}
		aRows.List[field.Name] = &oRow
		aRows.DBList[dbname] = &oRow
	}

	return &aRows
}

func (self *Searcher) GetRowsCount(table string) (int, error) {
	return self.GetCount(fmt.Sprintf("SELECT count(*) FROM %s", table), make([]interface{}, 0))
}

/*
*************************** ARRAY PARSER START ******************************
    For more infomation visit page https://gist.github.com/adharris/4163702
*/

// construct a regexp to extract values:
var (
	// unquoted array values must not contain: (" , \ { } whitespace NULL)
	// and must be at least one char
	unquotedChar  = `[^",\\{}\s(NULL)]`
	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

	// quoted array values are surrounded by double quotes, can be any
	// character except " or \, which must be backslash escaped:
	quotedChar  = `[^"\\]|\\"|\\\\`
	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

	// an array value may be either quoted or unquoted:
	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

	// Array values are separated with a comma IF there is more than one value:
	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

	valueIndex int
)

// Find the index of the 'value' named expression
func init() {
	for i, subexp := range arrayExp.SubexpNames() {
		if subexp == "value" {
			valueIndex = i
			break
		}
	}
}

func parseArray(array string) []string {
	results := make([]string, 0)
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		results = append(results, s)
	}
	return results
}

/*
*************************** ARRAY PARSER FINISH **************************
 */

func (s *Searcher) Do(sql string, values ...interface{}) {
	_, err := s.db.Exec(sql, values...)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Searcher) Insert(table string, data map[string]interface{}) {
	sql, values := sqler.InsertLine(table, data)
	s.DoCommit(sql, values)
}

func (s *Searcher) Delete(table string, data_where map[string]interface{}) {
	sql, values := sqler.DeleteLine(table, data_where)
	s.DoCommit(sql, values)
}

func (s *Searcher) Update(table string, data_where map[string]interface{}, data_update map[string]interface{}) {
	sql, values := sqler.UpdateLine(table, data_update, data_where)
	s.DoCommit(sql, values)
}

func (s *Searcher) DoCommit(sql string, values []interface{}) {
	log.Printf("DoCommit: %s\n", sql)
	spew.Dump(values)

	txn, err := s.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = txn.Exec(sql, values...)
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
