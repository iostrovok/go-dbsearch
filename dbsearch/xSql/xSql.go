package xSql

import (
	"fmt"
	"github.com/iostrovok/go-iutils/iutils"
	"log"
	//"reflect"
	"strings"
)

var MarkList = map[string]bool{
	"RET":   true,
	"SQL":   true,
	"IS":    true,
	"IN":    true,
	"LIKE":  true,
	"ILIKE": true,
	"=":     true,
	">=":    true,
	"<=":    true,
	"<>":    true,
	">":     true,
	"<":     true,
}

var MarkArrayList = map[string]bool{
	"=":  true,
	"<>": true,
	"<":  true,
	">":  true,
	"<=": true,
	">=": true,
	"@>": true,
	"<@": true,
	"&&": true,
	"||": true,
}

var LogicList = map[string]bool{
	"AND":    true,
	"OR":     true,
	"INSERT": true,
}

type One struct {
	Data    []interface{}
	Table   string
	Field   string
	Type    string
	NoVals  bool
	IsArray bool
}

func Insert(table string) *One {
	one := One{}
	one.Type = "INSERT"
	one.Table = table
	return &one
}

func IN(field string, data []interface{}) *One {
	one := One{}
	one.Data = data
	one.Type = "IN"
	one.Field = field

	return &one
}

func (one *One) CompInsert() (string, []interface{}) {
	sIn := []string{}
	sVals := []string{}
	sRet := []string{}
	Point := 1
	values := []interface{}{}
	for _, v := range one.Data {
		switch v.(type) {
		case *One:

			tp := v.(*One).Type

			if tp == "RET" {
				sRet = append(sRet, v.(*One).Field)
				continue
			}

			sIn = append(sIn, v.(*One).Field)
			vals := v.(*One).Data

			if v.(*One).IsArray {
				s, v := PrepaperArray(vals, Point)
				sVals = append(sVals, s)
				values = append(values, v...)
				Point += len(v)
				continue
			}

			if len(vals) > 1 {
				log.Fatalf("Comp. You can't INSERT multivalue params %T, %v\n", v, v)
			}

			if tp == "SQL" {
				sVals = append(sVals, iutils.AnyToString(vals[0]))
			} else {
				sVals = append(sVals, fmt.Sprintf("$%d ", Point))
				Point++
				values = append(values, vals...)
			}

		default:
			log.Printf("Comp. Not defined %T\n", v)
			log.Fatalf("Comp. Not defined %v\n", v)
		}
	}
	ret := ""
	if len(sRet) > 0 {
		ret = " RETURNING " + strings.Join(sRet, ", ")
	}
	sql := "  ( " + strings.Join(sIn, ", ") + ") VALUES (" + strings.Join(sVals, ", ") + ")"
	return " INSERT INTO " + one.Table + sql + ret, values
}

func (one *One) Comp(PointIn ...int) (string, []interface{}) {

	Point := 1
	if len(PointIn) > 0 {
		Point = PointIn[0]
	}

	sqlLine := ""
	values := []interface{}{}

	if one.Type == "INSERT" {
		if Point > 1 {
			log.Fatalf("Comp. You can't combination INSERT into other request\n")
		}
		return one.CompInsert()
	}

	if one.NoVals {
		return one.Field, values
	}

	if one.IsArray {
		return one.CompArray(Point)
	}

	switch one.Type {
	case "AND", "OR":
		s := []string{}
		for _, v := range one.Data {
			switch v.(type) {
			case *One:
				sql, vals := v.(*One).Comp(Point)
				s = append(s, sql)
				if len(vals) > 0 {
					Point += len(vals)
					values = append(values, vals...)
				}
			default:
				log.Printf("Comp. Not defined %T\n", v)
				log.Fatalf("Comp. Not defined %v\n", v)
			}
		}
		sqlLine = "( " + strings.Join(s, " "+one.Type+" ") + ") "
	case "IN":
		s := []string{}
		i := len(one.Data)
		for {
			s = append(s, fmt.Sprintf(" $%d ", Point))
			Point++
			i--
			if i < 1 {
				break
			}
		}
		sqlLine = fmt.Sprintf(" %s IN ( %s ) ", one.Field, strings.Join(s, ", "))
		values = one.Data
	case "IS":
		sqlLine = one.Field + " IS " + iutils.AnyToString(one.Data[0])
	default:
		sqlLine = fmt.Sprintf(" %s %s $%d ", one.Field, one.Type, Point)
		Point++
		values = append(values, one.Data[0])
	}

	return sqlLine, values
}

func (one *One) CompArray(PointIn ...int) (string, []interface{}) {

	Point := 1
	if len(PointIn) > 0 {
		Point = PointIn[0]
	}

	if !one.IsArray {
		log.Fatalf("CompArray. It does not have array type: %v\n", one)
	}

	sqlLine, values := PrepaperArray(one.Data, Point)
	sqlLine = fmt.Sprintf(" %s %s %s ", one.Field, one.Type, sqlLine)

	return sqlLine, values
}

func (one *One) Append(Nexters ...*One) *One {
	if _, find := LogicList[one.Type]; !find {
		log.Fatalf("Append. Bad type for append %s\n", one.Type)
	}
	for _, v := range Nexters {
		one.Data = append(one.Data, v)
	}
	return one
}

func NLogic(mark string, Nexters ...*One) *One {
	one := One{}
	one.Logic(mark, Nexters...)
	return &one
}

func (one *One) Logic(mark string, Nexters ...*One) {

	if _, find := LogicList[mark]; !find {
		log.Fatalf("Logic. Not defined %s\n", mark)
	}

	one.Data = []interface{}{}
	for _, v := range Nexters {
		one.Data = append(one.Data, v)
	}
	one.Type = mark
}

/*
	example: where start_date > now()
	Func("start_date", ">", "now()")
*/
func Func(field string) *One {
	In := One{}

	In.NoVals = true
	In.Field = field

	return &In
}

func Mark(field string, mark string, data ...interface{}) *One {
	In := One{}

	if _, find := MarkList[mark]; find {
		In.Type = mark
	} else {
		log.Fatalf("Mark. Not defined %s\n", mark)
	}

	In.Data = data
	In.Field = field
	In.IsArray = false

	if mark == "IS" {
		if iutils.AnyToString(In.Data[0]) == "NULL" {
			In.NoVals = true
			In.Field = field + " IS NULL "
		} else if iutils.AnyToString(In.Data[0]) == "NOT NULL" {
			In.NoVals = true
			In.Field = field + " IS NOT NULL "
		} else {
			log.Fatalf("Mark. Not defined %s. You have to use 'IS', 'NULL' or 'IS', 'NOT NULL' \n", mark)
		}
	}

	/*
		if mark == "ARRAY" {
			log.Printf("++ %s\n", reflect.TypeOf(data[0]).Kind().String())
			if reflect.TypeOf(data[0]).Kind().String() != "slice" {
				log.Fatalf("Mark. For using \"ARRAY\" type, need xSql.Mark(<string>, \"ARRAY\", <slice>)\n")
			}
		}
	*/
	return &In
}

func Array(field string, mark string, data ...interface{}) *One {
	In := One{}

	if _, find := MarkArrayList[mark]; find {
		In.Type = mark
	} else {
		log.Fatalf("Array. Not defined %s\n", mark)
	}

	In.Data = data
	In.Field = field
	In.IsArray = true

	return &In
}
