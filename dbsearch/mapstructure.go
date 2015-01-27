package dbsearch

import (
	//"database/sql"
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

func (s *Searcher) _initGet(aRows *AllRows, p interface{}, sqlLine string,
	values ...[]interface{}) (*GetRowResultStr, error) {

	s.LastCols = []string{}
	if err := s.PreInit(aRows, p); err != nil {
		return nil, err
	}

	value := []interface{}{}
	if len(values) > 0 {
		value = values[0]
	}

	if s.log {
		log.Printf("dbsearch.Get: %s\n%#v\n", sqlLine, values)
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
	R, err_rr := s.prepare_fork_raw_result(aRows, cols)
	if err_rr != nil {
		return nil, err_rr
	}

	R.Rows = rows

	if s.log {
		line := "GetRowResultStr.Cols: %#v\nGetRowResultStr.Dest: %#v\n" +
			"GetRowResultStr.RawResult: %#v\nGetRowResultStr.Rows: %#v\n" +
			"GetRowResultStr.SkipList: %#v\n"
		log.Printf(line, R.Cols, R.Dest, R.RawResult, R.Rows, R.SkipList)
	}

	return R, nil
}

func (aRows *AllRows) GetRowResultFace(R *GetRowResultStr) (map[string]interface{}, error) {

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

func GetRowResultFaceRoutine(Point int, dataCh, resCh chan *EnvelopeRowResult) {

	lastE := &EnvelopeRowResult{Point: Point, IsLast: true}
	can_run := true
	for can_run {
		var E *EnvelopeRowResult
		select {
		case E, can_run = <-dataCh:
			if can_run {
				E.R.RawResult = E.RawResult
				resultStr, err := E.aRows.GetRowResultFace(E.R)
				if err != nil {
					E.Err = err
				} else {
					E.ResM = resultStr
				}
			} else {
				E = lastE
			}
		}
		select {
		case resCh <- E:
		}
		if !can_run {
			return
		}
	}

	resCh <- lastE
}

func GetRowResultRoutine(Point int, dataCh, resCh chan *EnvelopeRowResult) {

	lastE := &EnvelopeRowResult{Point: Point, IsLast: true}
	can_run := true
	for can_run {
		var E *EnvelopeRowResult
		select {
		case E, can_run = <-dataCh:
			if can_run {
				resultDB := map[string]interface{}{}
				for i, raw := range E.RawResult {
					if E.R.SkipList[i] {
						resultDB[E.aRows.DBList[E.R.Cols[i]].Name] = raw
					}
				}
				resultStr := reflect.New(E.aRows.SType).Interface()
				if err := convert(E.aRows.List, "", resultDB, resultStr); err != nil {
					E.Err = err
				} else {
					E.Res = resultStr
				}
			} else {
				E = lastE
			}
		}
		select {
		case resCh <- E:
		}
		if !can_run {
			return
		}
	}

	resCh <- lastE
}

func (aRows *AllRows) GetRowResult(R *GetRowResultStr) interface{} {

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

func (s *Searcher) prepare_fork_raw_result(aRows *AllRows, cols []string) (*GetRowResultStr, error) {

	SkipList := map[int]bool{}
	R := &GetRowResultStr{UseFork: false}

	for i := 0; i < len(cols); i++ {
		var fn ElemConvertFunc

		t, find := aRows.DBList[cols[i]]
		if !find {
			if s.DieOnColsName {
				return nil, fmt.Errorf("dbsearch.Get not found column: %s!", cols[i])
			}

			SkipList[i] = false

			fn = func() interface{} {
				return make([]byte, 0)
			}
			R.DestFunL = append(R.DestFunL, fn)
			continue
		}
		SkipList[i] = true
		switch t.Type {
		case "date", "time", "timestamp":
			fn = func() interface{} {
				return new(*time.Time)
			}
		case "int", "bigint", "smallint", "integer", "serial", "bigserial":
			fn = func() interface{} {
				return new(int)
			}
		case "real", "double", "numeric", "decimal":
			fn = func() interface{} {
				return new(float64)
			}
		case "text", "varchar", "char", "bool", "money":
			fn = func() interface{} {
				return new(string)
			}
		case "[]text", "[]varchar", "[]char", "[]bool",
			"[]real", "[]double", "[]numeric", "[]decimal", "[]money",
			"[]int", "[]bigint", "[]smallint", "[]integer":
			fn = func() interface{} {
				return make([]byte, 0)
			}
			R.UseFork = true
		case "json", "jsonb":
			fn = func() interface{} {
				return make([]byte, 0)
			}
			R.UseFork = true
		default:
			fn = func() interface{} {
				return make([]byte, 0)
			}
			R.UseFork = true
		}
		R.DestFunL = append(R.DestFunL, fn)
	}

	R.Dest, R.RawResult = R.PrepareDestFun()
	R.SkipList = SkipList
	R.Cols = cols

	return R, nil
}
