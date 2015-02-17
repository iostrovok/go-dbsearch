package dbsearch

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime"
)

// OneRow includes data about single field of structure
type OneRow struct {
	Count   int
	DBName  string
	FType   string
	Name    string
	SetFunc convertData
	Type    string
	Log     bool
	Skip    bool
}

/*
AllRows includes list of *OneRow
	DBList is map[<column name in DB>]*OneRow
	List   is map[<field name in structure>]*OneRow
*/
type AllRows struct {
	TableInfo     *OneTableInfo
	DBList        map[string]*OneRow
	List          map[string]*OneRow
	Done          bool
	SType         reflect.Type
	Table         string
	Schema        string
	DieOnColsName bool
	Log           int
}

/*
PreInit is contains preprocessing of our structure, part 1.
We can call method when we have "db" and "type"
for each field in our structure
*/
func (aRows *AllRows) PreInit() {
	if !aRows.Done {
		m.Lock()
		aRows._iPrepare()
		m.Unlock()
	}
}

/*
PreInit contains preprocessing of our structure, part 1.
Part 2 is in aRows._iPrepare().
Connect to db is necessary.
*/
func (s *Searcher) PreInit(aRows *AllRows, p ...interface{}) error {
	if !aRows.Done {
		m.Lock()

		// Really rare case
		if aRows.Done {
			m.Unlock()
			return nil
		}

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
				return fmt.Errorf("No defined field %s.SType in\n%s line %d\n"+
					"%s line %d\n%s line %d\n", reflect.TypeOf(aRows),
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

		if err := aRows._iPrepare(); err != nil {
			return err
		}
		m.Unlock()
	}
	return nil
}

/*
	This method contains preprocessing of our structure, part 1.
	Connect to db is not necessary.

	AND
	We can call method when we have "db" and "type"
	for each field in our structure (from user or (s *Searcher) PreInit())
*/
func (aRows *AllRows) _iPrepare() error {

	if aRows.SType == nil {
		_, file1, line1, _ := runtime.Caller(2)
		_, file2, line2, _ := runtime.Caller(3)
		_, file3, line3, _ := runtime.Caller(4)
		return fmt.Errorf("No defined field %s.SType in\n%s line %d"+
			"\n%s line %d\n%s line %d\n", reflect.TypeOf(aRows),
			file1, line1, file2, line2, file3, line3)
	}

	st := reflect.TypeOf(reflect.New(aRows.SType).Interface()).Elem()

	aRows.Done = true
	aRows.List = make(map[string]*OneRow, 0)
	aRows.DBList = make(map[string]*OneRow, 0)
	var err error

	Count := 0
	for true {
		field := st.Field(Count)

		if field.Name == "" {
			break
		}

		fName := field.Name

		Count++

		dbType := ""
		isSkip := false
		dbName := field.Tag.Get("db")
		if dbName == "-" {
			isSkip = true
		} else {
			if dbName == "" {
				if a, f := aRows.GetFieldInfo(fName); f {
					dbName = a.Col
				}
			}
			if dbName == "" {
				if aRows.DieOnColsName {
					t := fmt.Sprintf("%s", field.Type)
					return errors.New(aRows.panicInitMessage("field_name", fName, t))
				}
				if aRows.Log > 0 {
					log.Printf("Warning for %s.%s. Not found field for '%s'\n",
						aRows.SType, fName, field.Type)
				}
				continue
			}

			dbType = field.Tag.Get("type")
			if dbType == "" {
				if a, f := aRows.GetColInfo(dbName); f {
					dbType = a.Type
				}
			}

			if dbType == "" {
				if aRows.DieOnColsName {
					return errors.New(aRows.panicInitMessage("db_type", fName, dbName))
				}

				if aRows.Log > 0 {
					log.Printf("Warning for %s.%s. Not found 'db' tag for '%s'\n",
						aRows.SType, fName, field.Type)
				}
				continue
			}
		}
		/*
			Makes OneRow for each field
		*/
		oRow := &OneRow{
			Name:   fName,
			DBName: dbName,
			Type:   dbType,
			Count:  Count,
			FType:  field.Type.String(),
			Skip:   isSkip,
		}

		aRows.List[fName] = oRow
		aRows.DBList[dbName] = oRow

		oRow.SetFunc, err = aRows.convertSelect(*oRow, fName, field.Type)
		if err != nil {
			return err
		}
	}

	return nil
}
