package xSql

import (
	"fmt"
	"github.com/iostrovok/go-iutils/iutils"
	"log"
	"strings"
)

var MarkList = map[string]bool{
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

var LogicList = map[string]bool{
	"AND": true,
	"OR":  true,
}

type One struct {
	Data   []interface{}
	Field  string
	Type   string
	NoVals bool
}

func IN(field string, data []interface{}) *One {
	one := One{}
	one.Data = data
	one.Type = "IN"
	one.Field = field

	return &one
}

func (one *One) Comp(PointIn ...int) (string, []interface{}) {

	Point := 1
	if len(PointIn) > 0 {
		Point = PointIn[0]
	}

	sqlLine := ""
	values := []interface{}{}
	if one.NoVals {
		return one.Field, values
	}

	switch one.Type {
	case "AND", "OR":
		s := []string{}
		for _, v := range one.Data {
			switch v.(type) {
			case *One:
				sql, vals := v.(*One).Comp(Point)
				Point += len(vals)
				s = append(s, sql)
				if len(vals) > 0 {
					values = append(values, vals...)
				}
			default:
				log.Fatalf("Comp. Not defined %T\n", v)
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

	return &In
}
