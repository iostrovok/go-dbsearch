package dbsearch

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
)

func convert(List map[string]*OneRow, name string, data interface{}, Result interface{}) error {

	val := reflect.ValueOf(Result).Elem()

	if data == nil {
		val.Set(reflect.Zero(val.Type()))
		return nil
	}

	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if !dataVal.IsValid() {
		val.Set(reflect.Zero(val.Type()))
		return nil
	}

	structType := val.Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldName := structType.Field(i).Name

		el := List[fieldName]
		if el == nil {
			return fmt.Errorf("No setup fieldName [%s]\n", fieldName)
		}

		rawMapVal := dataVal.MapIndex(reflect.ValueOf(fieldName))
		if val.Field(i).CanSet() && rawMapVal != reflect.Zero(reflect.TypeOf(rawMapVal)).Interface() {
			if err := el.SetFunc(rawMapVal.Interface(), val.Field(i)); err != nil {
				return err
			}
		}
	}

	return nil
}

func _AnyToString(data interface{}) string {

	dataVal := reflect.ValueOf(data)

	out := ""
	switch dataVal.Kind() {
	case reflect.String:
		out = dataVal.String()
	case reflect.Bool:
		if dataVal.Bool() {
			out = "1"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out = strconv.FormatInt(dataVal.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out = strconv.FormatUint(dataVal.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		out = strconv.FormatFloat(dataVal.Float(), 'f', -1, 64)
	case reflect.Slice:
		// TODO : Bad solution
		if dataVal.Type().Elem().Kind() == reflect.Uint8 {
			out = string(dataVal.Interface().([]uint8))
		}
	}

	return out
}

func (s *Searcher) _initGet(aRows *AllRows, sqlLine string,
	values ...[]interface{}) (*GetRowResultStr, error) {

	s.LastCols = []string{}
	if err := s.PreInit(aRows); err != nil {
		return nil, err
	}

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
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	s.LastCols = cols
	R, err_rr := s.prepare_raw_result(aRows, cols)
	if err_rr != nil {
		return nil, err_rr
	}

	R.Rows = rows

	if s.log {
		log.Printf("GetRowResultStr.Cols: %#v\n", R.Cols)
		log.Printf("GetRowResultStr.Dest: %#v\n", R.Dest)
		log.Printf("GetRowResultStr.RawResult: %#v\n", R.RawResult)
		log.Printf("GetRowResultStr.Rows: %#v\n", R.Rows)
		log.Printf("GetRowResultStr.SkipList: %#v\n", R.SkipList)
	}

	return R, nil
}

type GetRowResultStr struct {
	Rows      *sql.Rows
	Cols      []string
	Dest      []interface{}
	RawResult []interface{}
	SkipList  map[int]bool
}

func (aRows *AllRows) GetRowResultFace(R *GetRowResultStr) (map[string]interface{}, error) {

	mCheckError(R.Rows.Scan(R.Dest...))

	val := reflect.Indirect(reflect.New(aRows.SType).Elem())
	out := map[string]interface{}{}
	for i, raw := range R.RawResult {
		if !R.SkipList[i] {
			continue
		}

		if len(R.Cols) < i-1 {
			return nil, fmt.Errorf("No setup\n")
		}

		DBName := R.Cols[i]
		el := aRows.DBList[DBName]

		fieldName := el.Name
		rawMapVal := reflect.ValueOf(raw)

		field := val.FieldByName(fieldName)

		if rawMapVal == reflect.Zero(reflect.TypeOf(rawMapVal)).Interface() {
			field.Set(reflect.Zero(field.Type()))
		} else {
			if err := el.SetFunc(rawMapVal.Interface(), field); err != nil {
				return nil, err
			}
		}

		out[el.DBName] = field.Interface()
	}

	return out, nil
}

func (aRows *AllRows) GetRowResult(R *GetRowResultStr) interface{} {

	mCheckError(R.Rows.Scan(R.Dest...))

	resultDB := map[string]interface{}{}
	for i, raw := range R.RawResult {
		if R.SkipList[i] {
			resultDB[aRows.DBList[R.Cols[i]].Name] = raw
		}
	}

	resultStr := reflect.New(aRows.SType).Interface()
	mCheckError(convert(aRows.List, "", resultDB, resultStr))
	return resultStr
}

func (s *Searcher) prepare_raw_result(aRows *AllRows, cols []string) (*GetRowResultStr, error) {

	SkipList := map[int]bool{}
	rawResult := make([]interface{}, 0)
	for i := 0; i < len(cols); i++ {
		t, find := aRows.DBList[cols[i]]
		if !find {
			if s.DieOnColsName {
				return nil, fmt.Errorf("dbsearch.Get not found column: %s!", cols[i])
			}

			SkipList[i] = false
			rawResult = append(rawResult, make([]byte, 0))
			continue
		}
		SkipList[i] = true

		switch t.Type {
		case "date", "time", "timestamp":
			datetime := new(*time.Time)
			//datetime := new(*pq.NullTime)
			rawResult = append(rawResult, datetime)
		case "int", "bigint", "smallint", "integer", "serial", "bigserial":
			rawResult = append(rawResult, new(int))
		case "real", "double", "numeric", "decimal", "money":
			rawResult = append(rawResult, new(float64))
		case "text", "varchar", "char", "bool":
			rawResult = append(rawResult, new(string))
		case "[]text", "[]varchar", "[]char", "[]bool",
			"[]real", "[]double", "[]numeric", "[]decimal", "[]money",
			"[]int", "[]bigint", "[]smallint", "[]integer":
			rawResult = append(rawResult, make([]byte, 0))
		case "json", "jsonb":
			rawResult = append(rawResult, make([]byte, 0))
		default:
			rawResult = append(rawResult, make([]byte, 0))
		}
	}

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	R := GetRowResultStr{
		Cols:      cols,
		Dest:      dest,
		RawResult: rawResult,
		SkipList:  SkipList,
	}

	return &R, nil
}
