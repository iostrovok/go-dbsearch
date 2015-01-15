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

	dataValKind := dataVal.Kind()
	if dataValKind != reflect.Map {
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataValKind)
	}

	dataValType := dataVal.Type()
	if kind := dataValType.Key().Kind(); kind != reflect.String && kind != reflect.Interface {
		return fmt.Errorf(
			"'%s' needs a map with string keys, has '%s' keys",
			name, dataValType.Key().Kind())
	}

	fields := make(map[*reflect.StructField]reflect.Value)
	structType := val.Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)

		if fieldType.Anonymous {
			return fmt.Errorf("%s: unsupported type: %s", fieldType.Name, fieldType.Type.Kind())
		}

		fields[&fieldType] = val.Field(i)
		if !fields[&fieldType].IsValid() {
			panic("field is not valid")
		}
	}

	for fieldType, field := range fields {
		fieldName := fieldType.Name

		el := List[fieldName]
		if el == nil {
			log.Fatalf("No setup fieldName [%s]\n", fieldName)
		}

		rawMapKey := reflect.ValueOf(fieldName)
		rawMapVal := dataVal.MapIndex(rawMapKey)

		if field.CanSet() && rawMapVal != reflect.Zero(reflect.TypeOf(rawMapVal)).Interface() {
			if err := el.SetFunc(rawMapVal.Interface(), field); err != nil {
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
	s.PreInit(aRows)

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
	R, err_rr := aRows.prepare_raw_result(cols)
	if err_rr != nil {
		return nil, err_rr
	}

	R.Rows = rows

	return R, nil
}

type GetRowResultStr struct {
	Rows      *sql.Rows
	Cols      []string
	Dest      []interface{}
	RawResult []interface{}
}

func (aRows *AllRows) GetRowResult(R *GetRowResultStr) interface{} {

	mCheckError(R.Rows.Scan(R.Dest...))

	resultDB := map[string]interface{}{}
	for i, raw := range R.RawResult {
		resultDB[aRows.DBList[R.Cols[i]].Name] = raw
	}

	resultStr := reflect.New(aRows.SType).Interface()
	mCheckError(convert(aRows.List, "", resultDB, resultStr))
	return resultStr
}

func (aRows *AllRows) prepare_raw_result(cols []string) (*GetRowResultStr, error) {
	rawResult := make([]interface{}, 0)
	for i := 0; i < len(cols); i++ {
		t, find := aRows.DBList[cols[i]]
		if !find {
			return nil, fmt.Errorf("dbsearch.Get not found column: %s!", cols[i])
		}

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
	}

	return &R, nil
}