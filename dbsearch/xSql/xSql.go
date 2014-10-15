package xSql

import (
	"fmt"
	"github.com/iostrovok/go-iutils/iutils"
	"log"
	"strings"
)

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
}

type One struct {
	Data     []interface{}
	Table    string
	Field    string
	Mark     string
	AddParam string
	Type     string // NoVals Array Where JSON
}

func Update(table string) *One {
	one := One{}
	one.Mark = "UPDATE"
	one.Table = table
	return &one
}

func Insert(table string) *One {
	one := One{}
	one.Mark = "INSERT"
	one.Table = table
	return &one
}

func IN(field string, data []interface{}) *One {
	one := One{}
	one.Data = data
	one.Mark = "IN"
	one.Field = field

	return &one
}

func (one *One) CompUpdate() (string, []interface{}) {
	sUp := []string{}
	sRet := []string{}

	values := []interface{}{}
	sql_where := ""

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

			tp := v.(*One).Mark

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
	sIn := []string{}
	sVals := []string{}
	sRet := []string{}
	Point := 1
	values := []interface{}{}
	for _, v := range one.Data {
		switch v.(type) {
		case *One:

			tp := v.(*One).Mark

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

	Point := 1
	if len(PointIn) > 0 {
		Point = PointIn[0]
	}

	sqlLine := ""
	values := []interface{}{}

	if one.Mark == "INSERT" {
		if Point > 1 {
			log.Fatalf("Comp. You can't combination INSERT into other request\n")
		}
		return one.CompInsert()
	}

	if one.Mark == "UPDATE" {
		if Point > 1 {
			log.Fatalf("Comp. You can't combination INSERT into other request\n")
		}
		return one.CompUpdate()
	}

	switch one.Type {
	case "NoVals":
		return one.Field, values
	case "Array":
		return one.CompArray(Point)
	case "JSON":
		one.Data = PrepareJsonVals(one.Data)
	}

	switch one.Mark {
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
		sqlLine = "( " + strings.Join(s, " "+one.Mark+" ") + ") "
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
		sqlLine = fmt.Sprintf(" %s %s $%d ", one.Field, one.Mark, Point)
		Point++
		values = append(values, one.Data[0])
	}

	return sqlLine, values
}

func (one *One) Append(Nexters ...*One) *One {
	if _, find := LogicList[one.Mark]; !find {
		log.Fatalf("Append. Bad type for append %s\n", one.Mark)
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
	one.Mark = mark
}

func (one *One) Where(v *One) *One {
	v.Type = "Where"
	one.Data = append(one.Data, v)
	return one
}

/*
	example: where start_date > now()
	Func("start_date", ">", "now()")
*/
func Func(field string) *One {
	In := One{}

	In.Type = "NoVals"
	In.Field = field

	return &In
}

func Mark(field string, mark string, data ...interface{}) *One {
	In := One{}

	if _, find := MarkList[mark]; find {
		In.Mark = mark
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
