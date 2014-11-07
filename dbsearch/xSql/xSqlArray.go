package xSql

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

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

func (one *One) CompArray(PointIn ...int) (string, []interface{}) {

	Point := 1
	if len(PointIn) > 0 {
		Point = PointIn[0]
	}

	if one.Type != "Array" {
		log.Fatalf("CompArray. It does not have array type: %v\n", one)
	}

	sqlLine, values := PrepareArray(one.AddParam, one.Data, Point)
	sqlLine = fmt.Sprintf(" %s %s %s ", one.Field, one.Marker, sqlLine)

	return sqlLine, values
}

func PrepareArray(TypeArray string, val []interface{}, start_point int) (string, []interface{}) {

	if len(val) == 0 {
		return " '{}'" + TypeArray + " ", []interface{}{}
	}

	line := []string{}
	values := []interface{}{}
	for _, v := range val {
		line = append(line, "$"+strconv.Itoa(start_point))
		values = append(values, v)
		start_point++
	}

	return "ARRAY[ " + strings.Join(line, ", ") + " ]" + TypeArray, values
}

func (one *One) TArray(type_array string, field string, mark string, data ...interface{}) *One {
	one.Append(TArray(type_array, field, mark, data...))
	return one
}

func TArray(type_array string, field string, mark string, data ...interface{}) *One {
	In := Array(field, mark, data...)
	In.AddParam = "::" + type_array + "[]"
	return In
}

func (one *One) Array(field string, mark string, data ...interface{}) *One {
	one.Append(Array(field, mark, data...))
	return one
}

func Array(field string, mark string, data ...interface{}) *One {
	In := One{}

	if _, find := MarkArrayList[mark]; find {
		In.Marker = mark
	} else {
		log.Fatalf("Array. Not defined %s\n", mark)
	}

	In.Data = data
	In.Field = field
	In.Type = "Array"

	return &In
}
