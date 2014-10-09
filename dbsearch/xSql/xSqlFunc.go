package xSql

import (
	"strconv"
	"strings"
)

func PrepaperArray(TypeArray string, val []interface{}, start_point int) (string, []interface{}) {

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
