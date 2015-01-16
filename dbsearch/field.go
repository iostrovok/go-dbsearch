package dbsearch

import (
	"fmt"
	"log"
	"strings"
	"unicode"
)

//
type OneRowInfo struct {
	Col   string
	Type  string
	Field string
}

type OneTableInfo struct {
	Rows   map[string]*OneRowInfo
	Done   bool
	Table  string
	Schema string
}

func (aRows *AllRows) GetColInfo(col string) (*OneRowInfo, bool) {
	if aRows.TableInfo != nil {
		if v, f := aRows.TableInfo.Rows[col]; f {
			return v, true
		}
	}
	return nil, false
}

func (aRows *AllRows) GetFieldInfo(field string) (*OneRowInfo, bool) {
	if aRows.TableInfo != nil {
		for _, v := range aRows.TableInfo.Rows {
			if v.Field == field {
				return v, true
			}
		}
	}
	return nil, false
}

func (aRows *AllRows) PreinitTable() {

	if aRows.Table == "" {
		return
	}

	if aRows.Schema == "" {
		aRows.Schema = "public"
	}

	aRows.TableInfo = &OneTableInfo{
		Table:  aRows.Table,
		Schema: aRows.Schema,
	}
}

func _field_name(column_name string) string {
	parts := strings.Split(strings.ToLower(column_name), "_")
	out := ""
	for _, v := range parts {
		if v != "" {
			a := []rune(v)
			a[0] = unicode.ToUpper(a[0])
			out += string(a)
		}
	}

	return out
}

func _field_type(data_type string, udt_name string) (string, error) {

	udt_name = strings.TrimPrefix(strings.ToLower(udt_name), "_")
	out := ""

	switch udt_name {
	case "bool":
		out = "bool"
	case "int8":
		out = "bigint"
	case "int2":
		out = "smallint"
	case "int4":
		out = "int"
	case "text":
		out = "text"
	case "varchar":
		out = "varchar"
	case "bpchar":
		out = "char"
	case "float4":
		out = "real"
	case "float8":
		out = "double"
	case "numeric", "decimal":
		out = "numeric"
	case "money":
		out = "money"
	case "date":
		out = "date"
	case "time":
		out = "time"
	case "timestamptz", "timestamp":
		out = "timestamp"
	case "json", "jsonb":
		out = "json"
	}

	if data_type == "ARRAY" {
		out = "[]" + out
	}

	/* ===> Don't support now */
	/*
	 "bit", "box", "bytea", "cidr", "circle", "inet", "interval", "line",
	 "lseg", "macaddr", "path", "point", "polygon", "tsquery", "tsvector", "txid_snapshot",
	 "uuid", "varbit", "xml",
	*/
	/* <=== Don't support now */

	if out == "" {
		return "", fmt.Errorf("The type '%s' [%s] don't support now\n", udt_name, data_type)
	}
	return out, nil
}

func (s *Searcher) GetTableData(Table *OneTableInfo) error {

	if Table.Done {
		return nil
	}

	Table.Rows = map[string]*OneRowInfo{}

	//sql := " SELECT column_name, column_default, is_nullable, data_type " +
	//" FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2"
	//sql := " SELECT column_name, column_default, is_nullable, data_type, character_maximum_length, udt_name " +
	//	" FROM information_schema.columns " +
	//	"WHERE table_schema = $1 AND table_name = $2"

	vals := []interface{}{strings.ToLower(Table.Schema), strings.ToLower(Table.Table)}
	sql := "SELECT column_name, column_default, is_nullable, " +
		"data_type, character_maximum_length, udt_name " +
		"FROM information_schema.columns " +
		"WHERE table_schema = $1 AND table_name = $2"

	rows, err := s.db.Query(sql, vals...)
	if err != nil {
		return err
	}
	defer rows.Close()

	check := false
	for rows.Next() {
		check = true

		var column_name interface{}
		var column_default interface{}
		var is_nullable interface{}
		var data_type interface{}
		var character_maximum_length interface{}
		var udt_name interface{}

		err := rows.Scan(&column_name, &column_default, &is_nullable, &data_type, &character_maximum_length, &udt_name)
		if err != nil {
			return err
		}

		ft, err_ft := _field_type(_AnyToString(data_type), _AnyToString(udt_name))
		if err_ft != nil {
			return err_ft
		}

		Row := OneRowInfo{
			Field: _field_name(_AnyToString(column_name)),
			Col:   _AnyToString(column_name),
			Type:  ft,
		}

		Table.Rows[Row.Col] = &Row

		if s.log {
			log.Printf("GetTableData field: '%s', col: '%s', type: '%s'\n", Row.Field, Row.Col, Row.Type)
		}
	}

	if !check {
		return fmt.Errorf("Not found table %s.%s\n",
			strings.ToLower(Table.Schema), strings.ToLower(Table.Table))
	}

	return nil
}
