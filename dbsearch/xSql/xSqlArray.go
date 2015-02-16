package xSql

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

//MarkArrayList contains actions which are available for array
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

func (one *One) compArray(PointIn ...int) (string, []interface{}) {

	Point := 1
	if len(PointIn) > 0 {
		Point = PointIn[0]
	}

	if one.Type != "Array" {
		log.Fatalf("compArray. It does not have array type: %v\n", one)
	}

	sqlLine, values := prepareArray(one.AddParam, one.Data, Point)
	sqlLine = fmt.Sprintf(" %s %s %s ", one.Field, one.Marker, sqlLine)

	return sqlLine, values
}

func prepareArray(TypeArray string, val []interface{}, startPoint int) (string, []interface{}) {

	if len(val) == 0 {
		return " '{}'" + TypeArray + " ", []interface{}{}
	}

	line := []string{}
	values := []interface{}{}
	for _, v := range val {
		line = append(line, "$"+strconv.Itoa(startPoint))
		values = append(values, v)
		startPoint++
	}

	return "ARRAY[ " + strings.Join(line, ", ") + " ]" + TypeArray, values
}

//TArray provides actions for array with type
func (one *One) TArray(typeArray string, field string, mark string, data ...interface{}) *One {
	one.Append(TArray(typeArray, field, mark, data...))
	return one
}

//TArray provides actions for array with type
func TArray(typeArray string, field string, mark string, data ...interface{}) *One {
	In := Array(field, mark, data...)
	In.AddParam = "::" + typeArray + "[]"
	return In
}

//Array provides actions for array without type
func (one *One) Array(field string, mark string, data ...interface{}) *One {
	one.Append(Array(field, mark, data...))
	return one
}

//Array provides actions for array without type
func Array(field string, mark string, data ...interface{}) *One {
	In := One{}

	if _, find := MarkArrayList[mark]; find {
		In.Marker = mark
	} else {
		log.Fatalf("Array. Not defined %s\n", mark)
	}

	if len(data) > 1 {
		In.Data = data
	} else if len(data) == 1 {
		In.Data = ToFace(data[0])
	}

	In.Field = field
	In.Type = "Array"

	return &In
}
