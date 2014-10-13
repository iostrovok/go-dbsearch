package xSql

import (
	"encoding/json"
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

func (one *One) Json() *One {
	one.Type = "JSON"
	return one
}
