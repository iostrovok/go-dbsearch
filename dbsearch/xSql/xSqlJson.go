package xSql

import (
	"encoding/json"
	"github.com/iostrovok/go-iutils/iutils"
	"log"
)

func PrepareJsonVals(val []interface{}) []interface{} {

	if len(val) == 0 {
		return []interface{}{""}
	}

	b, err_j := json.Marshal(val[0])
	if err_j != nil {
		log.Fatalf("PrepareJsonVals. Bad Json data: %s\n", err_j)
	}

	return []interface{}{string(b)}
}

func (one *One) Json(f ...interface{}) *One {

	if len(f) == 0 {
		one.Type = "JSON"
		return one
	}

	if len(f) < 2 {
		log.Fatalln("There Json(f ...interface{}) needs more then 2 params or nothing")
	}

	field := iutils.AnyToString(f[0])
	mark := iutils.AnyToString(f[1])
	data := f[2:]

	n_one := Mark(field, mark, data...)
	n_one.Type = "JSON"

	one.Append(n_one)

	return one
}
