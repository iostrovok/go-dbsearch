package dbsearch

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
)

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
