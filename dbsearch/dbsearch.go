package dbsearch

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iostrovok/go-dbsearch/dbsearch/sqler"
	"github.com/iostrovok/go-iutils/iutils"
	_ "github.com/lib/pq"
	"log"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

var m sync.Mutex

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
	Done   bool
}

type Searcher struct {
	db            *sql.DB
	log           bool
	DieOnColsName bool
	LastCols      []string
}

func (s *Searcher) SetDebug(is_debug ...bool) {

	if len(is_debug) > 0 {
		s.log = is_debug[0]
	} else {
		s.log = true
	}
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

func SetDBI(db *sql.DB) (*Searcher, error) {

	s := new(Searcher)
	s.db = db

	return s, nil
}

func DBI(poolSize int, dsn string, stop_error ...bool) (*Searcher, error) {

	s := new(Searcher)

	db, _ := sql.Open("postgres", dsn)

	if err := db.Ping(); err != nil {
		if len(stop_error) > 0 && stop_error[0] {
			log.Fatalf("DB Error: %s\n", err)
		}
		return nil, errors.New(fmt.Sprintf("DB Error: %s\n", err))
	}

	s.db = db
	s.db.SetMaxOpenConns(poolSize)

	return s, nil
}

func (s *Searcher) GetCount(sqlLine string, values []interface{}) (int, error) {
	var count int
	err := s.db.QueryRow(sqlLine, values...).Scan(&count)
	return count, err
}

func (s *Searcher) GetOne(mType *AllRows, sqlLine string, values ...[]interface{}) (map[string]interface{}, error) {

	sqlLine += " LIMIT 1 OFFSET 0 "

	value := []interface{}{}
	if len(values) > 0 {
		value = values[0]
	}

	list, err := s.Get(mType, sqlLine, value)
	empty := map[string]interface{}{}

	if err != nil {
		return empty, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return empty, nil
}

func (s *Searcher) Get(mType *AllRows, sqlLine string, values ...[]interface{}) ([]map[string]interface{}, error) {

	s.LastCols = []string{}

	Out := make([]map[string]interface{}, 0)

	value := []interface{}{}
	if len(values) > 0 {
		value = values[0]
	}

	if s.log {
		log.Printf("dbsearch.Get: %s\n", sqlLine)
		log.Printf("%v\n", values)
	}

	rows, err := s.db.Query(sqlLine, value...)
	if err != nil {
		return Out, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return Out, err
	}

	s.LastCols = cols

	rawResult := make([]interface{}, 0)
	for i := 0; i < len(cols); i++ {
		t, find := mType.DBList[cols[i]]
		if !find {
			log.Fatalf("dbsearch.Get not found column: %s!", cols[i])
		}

		switch t.Type {
		case "datetime", "date", "time":
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
			if s.log {
				log.Printf("parseArray. %s: %s\n", cols[i], raw)
			}

			result[cols[i]] = convertType(cols[i], mType, raw)
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

	if raw_in == nil {
		if t.IsArray {
			switch t.Type {
			case "text", "date", "datetime", "time":
				return []string{}
			case "bigint", "int64", "int":
				return []int{}
			case "real":
				return []float64{}
			}
			return []interface{}{}
		} else {
			return nil
		}
	}

	switch t.Type {
	case "text":
		raw := raw_in.([]byte)
		if t.IsArray {
			return parseArray(string(raw))
		} else {
			return string(raw)
		}
	case "json", "jsonb":
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
		if t.IsArray {
			return parseIntArray(raw_in)
		} else {
			return iutils.AnyToInt(raw_in)
		}
	case "real":
		if t.IsArray {
			return parseFloat64Array(raw_in)
		} else {
			return iutils.AnyToFloat64(raw_in)
		}
	case "date", "datetime", "time":
		return raw_in
	}

	return nil
}

func (mT *AllRows) PreInit(p interface{}) {
	if !mT.Done {
		m.Lock()
		mT.iPrepare(p)
		m.Unlock()
	}
}

func (aRows *AllRows) iPrepare(s interface{}) {
	st := reflect.TypeOf(s)

	aRows.Done = true
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
}

func Prepare(s interface{}) *AllRows {
	aRows := AllRows{}
	aRows.iPrepare(s)
	return &aRows
}

func (self *Searcher) GetRowsCount(table string) (int, error) {
	return self.GetCount(fmt.Sprintf("SELECT count(*) FROM %s", table), make([]interface{}, 0))
}

/*
*************************** ARRAY PARSER START ******************************
 */

// construct a regexp to extract values:
var (
	unquotedRe  = regexp.MustCompile(`([^",\\{}\s]|NULL)+,`)
	_arrayValue = fmt.Sprintf("\"(%s)+\",", `[^"\\]|\\"|\\\\`)
	quotedRe    = regexp.MustCompile(_arrayValue)

	noNumbers      = regexp.MustCompile(`[^-0-9]+`)
	noNumbersStart = regexp.MustCompile(`^[^-0-9]+`)
	noNumbersEnd   = regexp.MustCompile(`[^0-9]+$`)

	noNumberDots      = regexp.MustCompile(`[^-0-9\.,]+`)
	noNumberDotsSplit = regexp.MustCompile(`(,|\s+)+`)
)

func parseIntArray(s interface{}) []int {
	str := strings.TrimSpace(iutils.AnyToString(s))
	str = noNumbersStart.ReplaceAllString(str, "")
	str = noNumbersEnd.ReplaceAllString(str, "")
	return iutils.AnyToIntArray(noNumbers.Split(str, -1))
}

func parseFloat64Array(s interface{}) []float64 {
	out := []float64{}

	str := strings.TrimSpace(iutils.AnyToString(s))
	str = noNumberDots.ReplaceAllString(str, "")
	list := noNumberDotsSplit.Split(str, -1)

	for _, v := range list {
		out = append(out, iutils.AnyToFloat64(v))
	}

	return out
}

func parseArray(line string) []string {

	out := []string{}
	if line == "{}" {
		return out
	}

	if len(line)-1 != strings.LastIndex(line, "}") || strings.Index(line, "{") != 0 {
		return out
	}

	/* Removes lead & last {} and adds "," to end of string */
	line = strings.TrimPrefix(line, "{")
	line = strings.TrimSuffix(line, "}") + ","

	for len(line) > 0 {
		s := ""
		if strings.Index(line, `"`) != 0 {
			s = unquotedRe.FindString(line)
			line = line[strings.Index(line, ",")+1:]
			s = strings.TrimSuffix(s, ",")

			/* counvert NULL to empty string6 however we need nil string */
			if s == "NULL" {
				s = ""
			}
		} else {
			s = quotedRe.FindString(line)
			line = strings.TrimPrefix(line, s)
			s = strings.TrimPrefix(s, "\"")
			s = strings.TrimSuffix(s, "\",")
			s = strings.Join(strings.Split(s, "\\\\"), "\\")
			s = strings.Join(strings.Split(s, "\\\""), "\"")
		}
		out = append(out, s)
	}

	return out
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

	if s.log {
		log.Printf("DoCommit: %s\n", sql)
	}

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
