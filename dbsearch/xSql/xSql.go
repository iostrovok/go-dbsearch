package xSql

import (
	"fmt"
	"github.com/iostrovok/go-iutils/iutils"
	"log"
	"strings"
)

var ViewDebug bool = false

var MarkList = map[string]bool{
	"&&":    true,
	"<":     true,
	"<=":    true,
	"<>":    true,
	"<@":    true,
	"=":     true,
	">":     true,
	">=":    true,
	"@>":    true,
	"ILIKE": true,
	"IN":    true,
	"IS":    true,
	"LIKE":  true,
	"RET":   true,
	"SQL":   true,
	"||":    true,
}

var LogicList = map[string]bool{
	"AND":    true,
	"OR":     true,
	"INSERT": true,
	"UPDATE": true,
	"SELECT": true,
	"DELETE": true,
}

type One struct {
	Data     []interface{}
	Table    string
	Columns  string
	Field    string
	Marker   string
	AddParam string
	Type     string // NoVals Array Where JSON
}

func Select(table, columns string) *One {

	if ViewDebug {
		log.Printf("exe Select for %s, %s\n", table, columns)
	}

	one := One{}
	one.Marker = "SELECT"
	one.Table = table
	one.Columns = columns

	return &one
}

func Delete(table string) *One {

	if ViewDebug {
		log.Printf("exe Delete for %s\n", table)
	}

	one := One{}
	one.Marker = "DELETE"
	one.Table = table

	return &one
}

func Update(table string) *One {

	if ViewDebug {
		log.Printf("exe Update for %s\n", table)
	}

	one := One{}
	one.Marker = "UPDATE"
	one.Table = table
	return &one
}

func Insert(table string) *One {

	if ViewDebug {
		log.Printf("exe Insert for %s\n", table)
	}

	one := One{}
	one.Marker = "INSERT"
	one.Table = table
	return &one
}

func (one *One) _firstLogical() *One {
	for _, v := range one.Data {
		switch v.(type) {
		case *One:
			log.Printf("_firstLogical: len(one.Data) => %d\n", v.(*One).Marker)
			switch v.(*One).Marker {
			case "AND", "IN", "OR":
				return v.(*One)
			}
		}
	}
	return nil
}

func (one *One) CompSelect() (string, []interface{}) {

	if ViewDebug {
		log.Println("exe CompSelect")
	}

	if one.Columns == "" || "" == one.Table {
		log.Fatalln("Empty Table or Columns")
	}

	sql_s := "SELECT " + one.Columns + " FROM " + one.Table

	if len(one.Data) == 1 {
		sql, values := one.Data[0].(*One).Comp()
		return sql_s + " WHERE " + sql, values
	}

	if v := one._firstLogical(); v != nil {
		sql, values := v.Comp()
		return sql_s + " WHERE " + sql, values
	}

	return sql_s, []interface{}{}
}

func (one *One) CompDelete() (string, []interface{}) {
	sRet := []string{}

	values := []interface{}{}
	sql_where := ""

	if ViewDebug {
		log.Println("exe CompUpdate")
	}

	if len(one.Data) == 1 {
		sql_where, values = one.Data[0].(*One).Comp()
		if sql_where != "" {
			sql_where = " WHERE " + sql_where
		}
	}

	ret := ""
	if len(sRet) > 0 {
		ret = " RETURNING " + strings.Join(sRet, ", ")
	}
	return " DELETE FROM " + one.Table + sql_where + ret, values
}

func (one *One) CompUpdate() (string, []interface{}) {
	sUp := []string{}
	sRet := []string{}

	values := []interface{}{}
	sql_where := ""

	if ViewDebug {
		log.Println("exe CompUpdate")
	}

	for _, v := range one.Data {
		switch v.(type) {
		case *One:
			if v.(*One).Type == "Where" {
				sql, vals := v.(*One).Comp()
				if sql != "" {
					sql_where = " WHERE " + sql
					values = append(values, vals...)
				}
				break
			}
		}
	}

	Point := 1 + len(values)

	for _, v := range one.Data {
		switch v.(type) {
		case *One:
			if v.(*One).Type == "Where" {
				continue
			}

			tp := v.(*One).Marker

			if tp == "RET" {
				sRet = append(sRet, v.(*One).Field)
				continue
			}

			if tp != "=" {
				log.Fatalf("Comp. For update only \"=\" defined %v\n", v)
			}

			sql, vals := v.(*One).Comp(Point)
			Point += len(vals)
			sUp = append(sUp, sql)
			values = append(values, vals...)

		default:
			log.Fatalf("Comp. Not defined %T, %v\n", v, v)
		}
	}
	ret := ""
	if len(sRet) > 0 {
		ret = " RETURNING " + strings.Join(sRet, ", ")
	}
	sql := strings.Join(sUp, ", ")
	return " UPDATE " + one.Table + " SET " + sql + sql_where + ret, values
}

func (one *One) CompInsert() (string, []interface{}) {
	if ViewDebug {
		log.Println("exe CompInsert")
	}

	sIn := []string{}
	sVals := []string{}
	sRet := []string{}
	Point := 1
	values := []interface{}{}
	for _, v := range one.Data {
		switch v.(type) {
		case *One:

			tp := v.(*One).Marker

			if tp == "RET" {
				sRet = append(sRet, v.(*One).Field)
				continue
			}

			sIn = append(sIn, v.(*One).Field)
			vals := v.(*One).Data

			if v.(*One).Type == "Array" {
				s, v := PrepareArray(v.(*One).AddParam, vals, Point)
				sVals = append(sVals, s)
				values = append(values, v...)
				Point += len(v)
				continue
			}

			if len(vals) > 1 {
				log.Fatalf("Comp. You can't INSERT multivalue params %T, %v\n", v, v)
			}

			if v.(*One).Type == "JSON" {
				vals = PrepareJsonVals(vals)
			}

			if tp == "SQL" {
				sVals = append(sVals, iutils.AnyToString(vals[0]))
			} else {
				sVals = append(sVals, fmt.Sprintf("$%d ", Point))
				Point++
				values = append(values, vals...)
			}

		default:
			log.Fatalf("Comp. Not defined %T, %v\n", v, v)
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
	if ViewDebug {
		log.Printf("exe Comp for %s\n", one.Marker)
	}
	Point := 1
	if len(PointIn) > 0 {
		Point = PointIn[0]
	}

	sqlLine := ""
	values := []interface{}{}

	if one.Marker == "INSERT" {
		if Point > 1 {
			log.Fatalf("Comp. You can't combination INSERT into other request\n")
		}
		return one.CompInsert()
	}

	if one.Marker == "DELETE" {
		if Point > 1 {
			log.Fatalf("Comp. You can't combination DELETE into other request\n")
		}
		return one.CompDelete()
	}

	if one.Marker == "UPDATE" {
		if Point > 1 {
			log.Fatalf("Comp. You can't combination INSERT into other request\n")
		}
		return one.CompUpdate()
	}

	if one.Marker == "SELECT" {
		if Point > 1 {
			log.Fatalf("Comp. You can't combination INSERT into other request\n")
		}
		return one.CompSelect()
	}

	switch one.Type {
	case "NoVals":
		return one.Field, values
	case "Array":
		return one.CompArray(Point)
	case "JSON":
		one.Data = PrepareJsonVals(one.Data)
	}

	switch one.Marker {
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
				log.Fatalf("Comp. Not defined %T, %v\n", v, v)
			}
		}
		sqlLine = "( " + strings.Join(s, " "+one.Marker+" ") + ") "
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
		sqlLine = fmt.Sprintf(" %s %s $%d ", one.Field, one.Marker, Point)
		Point++
		values = append(values, one.Data[0])
	}

	return sqlLine, values
}

func (one *One) Append(Nexters ...*One) *One {
	if _, find := LogicList[one.Marker]; !find {
		log.Fatalf("Append. Bad type for append %s\n", one.Marker)
	}

	no_done := true
	if one.Marker == "SELECT" || one.Marker == "DELETE" {
		if v := one._firstLogical(); v != nil {
			no_done = false
			v.Append(Nexters...)
		}
	}

	if no_done {
		for _, v := range Nexters {
			one.Data = append(one.Data, v)
		}
	}

	return one
}

func (one *One) Where(v *One) *One {

	if ViewDebug {
		log.Printf("exe *One.Where\n")
	}

	v.Type = "Where"
	one.Data = append(one.Data, v)
	return one
}

/* ------------------------------ NEW ------------------- */
func NLogic(mark string, Nexters ...*One) *One {

	if ViewDebug {
		log.Printf("exe NLogic for %s\n", mark)
	}

	return Logic(mark, Nexters...)
}

func Logic(mark string, Nexters ...*One) *One {

	one := One{}

	if ViewDebug {
		log.Printf("exe *One.Logic for %s\n", mark)
	}

	if _, find := LogicList[mark]; !find {
		log.Fatalf("Logic. Not defined %s\n", mark)
	}

	one.Data = []interface{}{}
	for _, v := range Nexters {
		one.Data = append(one.Data, v)
	}
	one.Marker = mark
	return &one
}

func (one *One) Or(Nexters ...*One) *One {
	return one.Logic("OR", Nexters...)
}

func (one *One) And(Nexters ...*One) *One {
	return one.Logic("AND", Nexters...)
}

func (one *One) Logic(mark string, Nexters ...*One) *One {
	return one.Append(Logic(mark, Nexters...))
}

/*
	example: where start_date > now()
	Func("start_date", ">", "now()")
*/
func Func(field string) *One {

	if ViewDebug {
		log.Printf("exe Func for %s\n", field)
	}

	In := One{}

	In.Type = "NoVals"
	In.Field = field

	return &In
}

func (one *One) Func(field string) *One {
	if ViewDebug {
		log.Printf("exe *One.Func for %s\n", field)
	}

	return one.Append(Func(field))
}

func (one *One) Mark(field string, mark string, data ...interface{}) *One {
	if ViewDebug {
		log.Printf("exe *One.Mark for %s, %s INTO %s\n", field, mark, one.Marker)
	}

	return one.Append(Mark(field, mark, data...))
}

func Mark(field string, mark string, data ...interface{}) *One {

	if ViewDebug {
		log.Printf("exe Mark for %s, %s\n", field, mark)
	}

	In := One{}

	if _, find := MarkList[mark]; find {
		In.Marker = mark
	} else {
		log.Fatalf("Mark. Not defined %s\n", mark)
	}

	In.Data = data
	In.Field = field
	In.Type = ""

	if mark == "IS" {
		if iutils.AnyToString(In.Data[0]) == "NULL" {
			In.Type = "NoVals"
			In.Field = field + " IS NULL "
		} else if iutils.AnyToString(In.Data[0]) == "NOT NULL" {
			In.Type = "NoVals"
			In.Field = field + " IS NOT NULL "
		} else {
			log.Fatalf("Mark. Not defined %s. You have to use 'IS', 'NULL' or 'IS', 'NOT NULL' \n", mark)
		}
	}

	return &In
}

func (one *One) IN(field string, data []interface{}) *One {
	if ViewDebug {
		log.Printf("exe *One.IN for %s\n", field)
	}

	return one.Append(IN(field, data))
}

func IN(field string, data []interface{}) *One {
	one := One{}
	one.Data = data
	one.Marker = "IN"
	one.Field = field

	return &one
}
