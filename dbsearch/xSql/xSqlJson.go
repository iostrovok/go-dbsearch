package xSql

import (
	"encoding/json"
	"github.com/iostrovok/go-iutils/iutils"
	"log"
)

func prepareJSONvals(val []interface{}) []interface{} {

	if len(val) == 0 {
		return []interface{}{""}
	}

	b, errJ := json.Marshal(val[0])
	if errJ != nil {
		log.Fatalf("prepareJSONvals. Bad Json data: %s\n", errJ)
	}

	return []interface{}{string(b)}
}

//Json adds json param
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

	nOne := Mark(field, mark, data...)
	nOne.Type = "JSON"

	one.Append(nOne)

	return one
}
