package sqler

import (
	"github.com/iostrovok/go-iutils/iutils"
	"regexp"
	"strconv"
	"strings"
)

/*
	Using case:
	filters := {
		"id": "10",
		"name": "Join",
	}

	filter_text_like := {
		"family": "ar do",
	}

	sql, params := sqler.GetWhere( filters, filter_text_like )

	// sql => " WHERE is = $1 AND name = $2 AND family ilike $3 "
	// params ["10", "Join", "%ar%do%"]

*/

var NonDigitalRE = regexp.MustCompile(`[^0-9,\.]+`)

func GetWhere(filters ...map[string]interface{}) (string, []interface{}) {

	if len(filters) == 0 {
		return "", make([]interface{}, 0)
	}

	//NextId, where, params := GetWhereEqual(1, filters[0])
	_, where, params := GetWhereEqual(1, filters[0])

	whereOut := ""
	if len(where) > 0 {
		whereOut = " WHERE " + strings.Join(where, " AND ")
	}

	return whereOut, params
}

func GetWhereLike(start_point int, filters map[string]interface{}) (int, []string, []interface{}) {

	params := make([]interface{}, 0)
	where := make([]string, 0)

	for field, value := range filters {
		where = append(where, " "+field+" ilike $"+strconv.Itoa(start_point)+" ")
		words := strings.Fields(iutils.AnyToString(value))
		params = append(params, "%"+strings.Join(words, "%")+"%")
		start_point++
	}

	return start_point, where, params
}

func GetWhereEqual(start_point int, filters map[string]interface{}) (int, []string, []interface{}) {
	params := make([]interface{}, 0)
	where := make([]string, 0)
	for field, value := range filters {
		symbol := " = "
		mvalue := ""
		switch value.(type) {
		case map[string]interface{}:
			for f, v := range value.(map[string]interface{}) {
				symbol = " " + f + " "
				mvalue = iutils.AnyToString(v)
				break
			}
			where = append(where, " "+field+symbol+" $"+strconv.Itoa(start_point)+" ")
			params = append(params, mvalue)
			start_point++
		case []interface{}:
			line, spoint, mvalue := PrepaperArray(value.([]interface{}), start_point)
			if line != "" {
				start_point = spoint
				where = append(where, " "+field+" = "+line+" ")
				params = append(params, mvalue...)
			}
		default:
			mvalue = iutils.AnyToString(value)
			where = append(where, " "+field+symbol+" $"+strconv.Itoa(start_point)+" ")
			params = append(params, mvalue)
			start_point++
		}
	}

	return start_point, where, params
}

func PrepaperArray(val []interface{}, start_point int) (string, int, []interface{}) {

	if len(val) == 0 {
		return " '{}' ", start_point, []interface{}{}
	}

	line := []string{}
	values := []interface{}{}
	for _, v := range val {
		line = append(line, "$"+strconv.Itoa(start_point))
		values = append(values, iutils.AnyToString(v))
		start_point++
	}

	return " ARRAY[ " + strings.Join(line, ", ") + " ] ", start_point, values
}

func DeleteLine(table string, filters ...map[string]interface{}) (string, []interface{}) {

	where, params := GetWhere(filters...)
	sql := "DELETE FROM " + table + " " + where

	return sql, params
}

func UpdateLine(table string, data map[string]interface{}, filters ...map[string]interface{}) (string, []interface{}) {

	where, params := GetWhere(filters...)

	update := make([]string, 0)
	start_point := len(params) + 1
	for field, value := range data {

		switch value.(type) {
		case []interface{}:
			line, spoint, mvalue := PrepaperArray(value.([]interface{}), start_point)
			if line != "" {
				start_point = spoint
				update = append(update, field+" = "+line)
				params = append(params, mvalue...)
			}
		default:
			update = append(update, field+"=$"+strconv.Itoa(start_point))
			params = append(params, iutils.AnyToString(value))
			start_point++
		}

	}

	sql := "UPDATE " + table + " SET " + strings.Join(update, ", ") + where

	return sql, params
}

func InsertLine(table string, filters map[string]interface{}) (string, []interface{}) {
	params := make([]interface{}, 0)
	update := make([]string, 0)
	set := make([]string, 0)

	start_point := 1
	for field, value := range filters {
		switch value.(type) {
		case []interface{}:
			line, spoint, mvalue := PrepaperArray(value.([]interface{}), start_point)
			if line != "" {
				set = append(set, field)
				start_point = spoint
				params = append(params, mvalue...)
				update = append(update, line)
			}
		default:
			set = append(set, field)
			update = append(update, "$"+strconv.Itoa(start_point))
			params = append(params, iutils.AnyToString(value))
			start_point++
		}
	}

	sql := "INSERT INTO " + table + " (" + strings.Join(set, ", ") + ") VALUES (" + strings.Join(update, ", ") + ")"

	return sql, params
}
