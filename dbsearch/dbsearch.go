package dbsearch

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/iostrovok/go-dbsearch/dbsearch/sqler"
	"github.com/iostrovok/go-iutils/iutils"
	_ "github.com/lib/pq"
	"log"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// http://play.golang.org/p/zmwyIDpIPN
var m sync.Mutex

//
type OneRow struct {
	Count   int
	DBName  string
	FType   string
	Name    string
	SetFunc ConvertData
	Type    string
	Log     bool
}

type AllRows struct {
	TableInfo     *OneTableInfo
	DBList        map[string]*OneRow
	List          map[string]*OneRow
	SkipList      map[int]bool
	Done          bool
	SType         reflect.Type
	Table         string
	Schema        string
	DieOnColsName bool
	Log           int
}

type Searcher struct {
	db            *sql.DB
	poolSize      int
	dsn           string
	log           bool
	logFull       bool
	DieOnColsName bool
	LastCols      []string
	IsOneRec      bool
}

func (s *Searcher) Close() error {
	return s.db.Close()
}

func (s *Searcher) DBH() *sql.DB {
	return s.db
}

func (s *Searcher) SetDieOnColsName(is_debug ...bool) {
	if len(is_debug) > 0 {
		s.DieOnColsName = is_debug[0]
	} else {
		s.DieOnColsName = true
	}
}

func (s *Searcher) SetDebug(is_debug ...bool) {
	if len(is_debug) > 0 {
		s.log = is_debug[0]
	} else {
		s.log = true
	}
}

func (s *Searcher) SetStrongDebug(is_debug ...bool) {
	if len(is_debug) > 0 {
		s.logFull = is_debug[0]
	} else {
		s.logFull = true
	}
}

func (s *Searcher) Ping() error {
	return s.db.Ping()
}

func SetDBI(db *sql.DB) (*Searcher, error) {

	s := new(Searcher)
	s.db = db

	return s, nil
}

func (s *Searcher) StartReConnect(rto_in ...int) {

	rto := 5
	if len(rto_in) > 0 && rto_in[0] > 0 {
		rto = rto_in[0]
	}

	go func() {
		// Pings our DB each rto_in seconds
		for {
			time.Sleep(time.Duration(rto) * time.Second)
			if err := s.Ping(); err != nil {
				s.Close()

				if db, err := sql.Open("postgres", s.dsn); err != nil {
					log.Println(err)
				} else {
					s.db = db
					s.db.SetMaxOpenConns(s.poolSize)
				}
			}
		}
	}()
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

	s.poolSize = poolSize
	s.dsn = dsn

	return s, nil
}

func (s *Searcher) GetCount(sqlLine string, values []interface{}) (int, error) {
	var count int
	err := s.db.QueryRow(sqlLine, values...).Scan(&count)
	return count, err
}

func (s *Searcher) GetOne(mType *AllRows, p interface{}, sqlLine string, values ...[]interface{}) error {
	defer func() { s.IsOneRec = false }()

	R, err := s._initGet(mType, p, sqlLine, values...)
	if err != nil {
		return err
	}
	defer R.Rows.Close()
	for R.Rows.Next() {
		resultStr := mType.GetRowResult(R)
		reflect.Indirect(reflect.ValueOf(p)).Set(reflect.Indirect(reflect.ValueOf(resultStr)))
		break
	}

	mCheckError(R.Rows.Err())
	return nil
}

func (s *Searcher) Get(mType *AllRows, p interface{}, sqlLine string, values ...[]interface{}) error {
	defer func() { s.IsOneRec = false }()

	R, err := s._initGet(mType, p, sqlLine, values...)
	if err != nil {
		return err
	}
	defer R.Rows.Close()
	var sliceValue = reflect.Indirect(reflect.ValueOf(p))
	for R.Rows.Next() {
		resultStr := mType.GetRowResult(R)
		sliceValue.Set(reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(resultStr))))
		if s.IsOneRec {
			break
		}
	}

	mCheckError(R.Rows.Err())
	return nil
}

func (s *Searcher) GetFace(mType *AllRows, sqlLine string,
	values ...[]interface{}) ([]map[string]interface{}, error) {
	defer func() { s.IsOneRec = false }()

	out := []map[string]interface{}{}

	R, err := s._initGet(mType, nil, sqlLine, values...)
	if err != nil {
		return out, err
	}
	defer R.Rows.Close()

	for R.Rows.Next() {
		resultStr, err := mType.GetRowResultFace(R)
		if err != nil {
			return nil, err
		}

		out = append(out, resultStr)
		if s.IsOneRec {
			break
		}
	}

	mCheckError(R.Rows.Err())
	return out, nil
}

func (s *Searcher) GetFaceOne(mType *AllRows, sqlLine string,
	values ...[]interface{}) (map[string]interface{}, error) {

	s.IsOneRec = true
	out := map[string]interface{}{}

	list, err := s.GetFace(mType, sqlLine, values...)

	if err == nil {
		if len(list) > 0 {
			out = list[0]
		}
	}

	return out, err
}

func mCheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (aRows *AllRows) PreInit() {
	if !aRows.Done {
		m.Lock()
		aRows.iPrepare()
		m.Unlock()
	}
}

func (s *Searcher) PreInit(aRows *AllRows, p ...interface{}) error {
	if !aRows.Done {
		m.Lock()

		if s.logFull {
			aRows.Log = 2
		} else if s.log {
			aRows.Log = 1
		}

		if err := aRows.PreinitTable(); err != nil {
			return err
		}

		if aRows.SType == nil {
			if len(p) == 0 || p[0] == nil {
				_, file1, line1, _ := runtime.Caller(2)
				_, file2, line2, _ := runtime.Caller(3)
				_, file3, line3, _ := runtime.Caller(4)
				return fmt.Errorf("No defined field %s.SType in\n%s line %d\n%s line %d\n%s line %d\n", reflect.TypeOf(aRows),
					file1, line1, file2, line2, file3, line3)
			}

			aRows.SType = reflect.Indirect(reflect.ValueOf(p[0])).Type()
			if aRows.SType.Kind() == reflect.Slice {
				aRows.SType = aRows.SType.Elem()
			}
			if aRows.Log > 0 {
				log.Printf("Auto set type %s\n", aRows.SType)
			}
		}
		if aRows.TableInfo != nil {
			if err := s.GetTableData(aRows.TableInfo); err != nil {
				return err
			}
		}
		aRows.DieOnColsName = s.DieOnColsName

		if err := aRows.iPrepare(); err != nil {
			return err
		}
		m.Unlock()
	}
	return nil
}

func (aRows *AllRows) iPrepare() error {

	if aRows.SType == nil {
		_, file1, line1, _ := runtime.Caller(2)
		_, file2, line2, _ := runtime.Caller(3)
		_, file3, line3, _ := runtime.Caller(4)
		return fmt.Errorf("No defined field %s.SType in\n%s line %d\n%s line %d\n%s line %d\n", reflect.TypeOf(aRows),
			file1, line1, file2, line2, file3, line3)
	}

	st := reflect.TypeOf(reflect.New(aRows.SType).Interface()).Elem()

	aRows.Done = true
	aRows.List = make(map[string]*OneRow, 0)
	aRows.DBList = make(map[string]*OneRow, 0)
	var err error = nil

	Count := 0
	for true {
		field := st.Field(Count)

		if field.Name == "" {
			break
		}

		fieldName := field.Name

		fieldTypeType := field.Type
		fieldTypeTypeStr := fmt.Sprintf("%s", fieldTypeType)

		Count++

		dbname := field.Tag.Get("db")
		if dbname == "" {
			if a, f := aRows.GetFieldInfo(fieldName); f {
				dbname = a.Col
			}
		}
		if dbname == "" {
			if aRows.DieOnColsName {
				return errors.New(aRows.PanicInitMessage("field_name", fieldName, fieldTypeTypeStr))
			} else {
				if aRows.Log > 0 {
					log.Printf("Warning for %s.%s. Not found field for '%s'\n", aRows.SType, fieldName, fieldTypeTypeStr)
				}
				continue
			}
		}

		dbtype := field.Tag.Get("type")
		if dbtype == "" {
			if a, f := aRows.GetColInfo(dbname); f {
				dbtype = a.Type
			}
		}

		if dbtype == "" {
			if aRows.DieOnColsName {
				return errors.New(aRows.PanicInitMessage("db_type", fieldName, dbname))
			} else {
				if aRows.Log > 0 {
					log.Printf("Warning for %s.%s. Not found 'db' tag for '%s'\n", aRows.SType, fieldName, fieldTypeTypeStr)
				}
				continue
			}
		}

		oRow := OneRow{
			Name:   fieldName,
			DBName: dbname,
			Type:   dbtype,
			Count:  Count,
			FType:  field.Type.String(),
		}

		aRows.List[fieldName] = &oRow
		aRows.DBList[dbname] = &oRow

		oRow.SetFunc, err = aRows.convert_select(oRow, fieldTypeTypeStr, fieldName, fieldTypeType)
		if err != nil {
			return err
		}
	}

	return nil
}

func Prepare(s interface{}) *AllRows {
	aRows := AllRows{}
	aRows.iPrepare()
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

	intArrayBrace = regexp.MustCompile(`[^-0-9\.\,]+`)
	intArraySplit = regexp.MustCompile(`,`)
	intArrayTail  = regexp.MustCompile(`\.[0-9]*`)

	noNumberDots      = regexp.MustCompile(`[^-0-9\.,]+`)
	noNumberDotsSplit = regexp.MustCompile(`(,|\s+)+`)

	parseBoolArrayRe     = regexp.MustCompile(`[^FTft]+`)
	parseBoolArrayReTail = regexp.MustCompile(`^[^FTft]+|[^FTft]+$`)
)

func parseBoolArrayForBool(s interface{}) []bool {
	line := parseBoolArrayReTail.ReplaceAllString(_AnyToString(s), "")
	r := parseBoolArrayRe.Split(line, -1)
	out := make([]bool, len(r))
	for i, v := range r {
		if v == "t" || v == "T" {
			out[i] = true
		} else {
			out[i] = false
		}
	}
	return out
}

func parseBoolArrayForString(s interface{}) []bool {
	r := parseArray(_AnyToString(s))
	out := make([]bool, len(r))
	for i, v := range r {
		v := strings.TrimSpace(v)
		if v == "" {
			out[i] = false
		} else {
			out[i] = true
		}
	}
	return out
}

func parseBoolArrayForReal(s interface{}) []bool {
	r := parseFloat64Array(s)
	out := make([]bool, len(r))
	for i, v := range r {
		if v == 0.0 {
			out[i] = false
		} else {
			out[i] = true
		}
	}
	return out
}

func parseBoolArrayForNumber(s interface{}) []bool {
	r := parseInt64Array(s)
	out := make([]bool, len(r))
	for i, v := range r {
		if v == 0 {
			out[i] = false
		} else {
			out[i] = true
		}
	}
	return out
}

func parseUint64Array(s interface{}) []uint64 {
	r := parseInt64Array(s)
	out := make([]uint64, len(r))
	for i, v := range r {
		if v < 0 {
			out[i] = 0
		} else {
			out[i] = uint64(v)
		}
	}
	return out
}

func parseUint8Array(s interface{}) []uint8 {
	r := parseIntArray(s)
	out := make([]uint8, len(r))
	for i, v := range r {
		if v < 0 {
			out[i] = 0
		} else {
			out[i] = uint8(v)
		}
	}
	return out
}

func parseUintArray(s interface{}) []uint {
	r := parseIntArray(s)
	out := make([]uint, len(r))
	for i, v := range r {
		if v < 0 {
			out[i] = 0
		} else {
			out[i] = uint(v)
		}
	}
	return out
}

func parseInt64Array(s interface{}) []int64 {
	str := strings.TrimSpace(_AnyToString(s))
	str = intArrayBrace.ReplaceAllString(str, "")
	str = intArrayTail.ReplaceAllString(str, "")
	k := intArraySplit.Split(str, -1)

	out := make([]int64, len(k))

	for i, v := range k {
		v := intArrayTail.ReplaceAllString(v, "")

		if v == "" {
			out[i] = 0
			continue
		}

		j, err := strconv.Atoi(v)
		if err != nil {
			log.Println(err)
			out[i] = 0
			continue
		}
		out[i] = int64(j)
	}
	return out
}

func parseIntArray(s interface{}) []int {
	k := parseInt64Array(s)
	out := make([]int, len(k))
	for i, v := range k {
		out[i] = int(v)
	}
	return out
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

func parseFloat32Array(s interface{}) []float32 {
	out := []float32{}

	str := strings.TrimSpace(iutils.AnyToString(s))
	str = noNumberDots.ReplaceAllString(str, "")
	list := noNumberDotsSplit.Split(str, -1)

	for _, v := range list {
		out = append(out, float32(iutils.AnyToFloat64(v)))
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
		if strings.Index(line, `""`) == 0 {
			/* Empty line */
			s = ""
			line = line[3:]
		} else if strings.Index(line, `"`) != 0 {
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

func (s *Searcher) Update(table string, data_where, data_update map[string]interface{}) {
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
