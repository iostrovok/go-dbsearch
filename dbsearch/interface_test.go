package dbsearch

import (
	"reflect"
	"testing"
)

func Test_Interface(t *testing.T) {
	s := init_test_data(t)
	s.SetDebug(false)
	if s != nil {
		_01_interface_load(t, s, false)
		_01_interface_load(t, s, true)
		_02_interface_json(t, s, false)
		_02_interface_json(t, s, true)
	}
	//t.Fatal("Success [no error] test")
}

func interface_load_f_test_table(s *Searcher) {
	sql_create := " CREATE TABLE public.test " +
		"(col1 int, col2 bigint, col3 smallint, col4 integer, " +
		"col5 serial, col6 bigserial, col7 text, col8 varchar(50), col9 char(10), " +
		" col11 real, col12 double precision, col13 numeric, col14 decimal, col15 money, col16 boolean  " +
		") "
	sql_cols := "INSERT INTO test(col1, col2, col3, col4, col5, col6, col7, col8, col9, col11, col12, col13, col14, col15, col16 ) "
	sql_vals := []string{
		"VALUES (1, 9223372036854775807, 883, 884, 885, 886, '123456789', '123456789', '1234567890', 12.13, 14.15, 16.17, 18.19, 20.21, TRUE )",
		"VALUES (2, -9223372036854775807, -883, -884, -885, -886, '-123456789', '-123456789', '-123456789', -12.13, -14.15, -16.17, -18.19, -20.21, FALSE )",
		"VALUES (3, null, null, null, 0, 0, null, null, null, null, null, null, null, null, null )", // check null - nil
	}

	make_t_table(s, sql_create, sql_cols, sql_vals)
}

/*
	interface_load test
	Our struct has less fields than table has their
*/
type interface_load_TestPlace struct {
	Col1  int     `db:"col1" type:"int"`
	Col2  int     `db:"col2" type:"bigint"`
	Col3  int32   `db:"col3" type:"smallint"`
	Col4  int64   `db:"col4" type:"integer"`
	Col5  int32   `db:"col5" type:"serial"`
	Col6  int32   `db:"col6" type:"bigserial"`
	Col7  string  `db:"col7" type:"text"`
	Col8  string  `db:"col8" type:"varchar"`
	Col9  string  `db:"col9" type:"char"`
	Col11 float64 `db:"col11" type:"real"`
	Col12 float64 `db:"col12" type:"double"`
	Col13 float64 `db:"col13" type:"numeric"`
	Col14 float64 `db:"col14" type:"decimal"`
	Col15 float64 `db:"col15" type:"money"`
	Col16 bool    `db:"col16" type:"bool"`
}

var interface_load_mTestType *AllRows = &AllRows{
	Table: "test",
	SType: reflect.TypeOf(interface_load_TestPlace{}),
}

func _01_interface_load(t *testing.T, s *Searcher, asList bool) {
	interface_load_f_test_table(s)
	s.SetDieOnColsName(false)

	p := map[string]interface{}{}
	var err error = nil
	if asList {
		list, err := s.GetFace(interface_load_mTestType, "SELECT * FROM public.test ORDER BY 1")
		if err != nil {
			t.Fatal(err)
		}
		if len(list) > 0 {
			p = list[0]
		} else {
			t.Fatal("Error _01_interface_load. Point 1.")
		}
	} else {
		p, err = s.GetFaceOne(interface_load_mTestType, "SELECT * FROM public.test ORDER BY 1")
		if err != nil {
			t.Fatal(err)
		}
	}

	switch p["col2"].(type) {
	case int:
		if p["col2"].(int) != 9223372036854775807 {
			t.Fatal("Error _01_interface_load. Point 2. ")
		}
	default:
		t.Fatal("Error _01_interface_load. Point 3.")
	}
}

type interface_json_TestPlace struct {
	Col1 int
	Col2 map[string]interface{} `type:"json"`
	Col3 string                 `type:"json"`
	Col4 map[string]interface{}
	Col5 string
}

var interface_json_mTestType *AllRows = &AllRows{
	Table: "test",
	SType: reflect.TypeOf(interface_json_TestPlace{}),
}

func _02_interface_json(t *testing.T, s *Searcher, asList bool) {

	str := `'{"array":[{"one":1,"two":"two"},{"next":""}],"null":null,"false":false,"true":true,"ru":"Слова пишем слево на право"}'`

	sql_create := " CREATE TABLE public.test " +
		"(col1 int, col2 text, col3 text, col4 json, col5 json ) "
	sql_cols := "INSERT INTO test(col1, col2, col3, col4, col5 ) "
	sql_vals := []string{
		"VALUES (1, " + str + ", " + str + ", " + str + "::json, " + str + "::json )",
		"VALUES (2, null, null, null, null )", // check null - nil
	}

	make_t_table(s, sql_create, sql_cols, sql_vals)

	p := map[string]interface{}{}
	var err error = nil
	if asList {
		list, err := s.GetFace(interface_json_mTestType, "SELECT * FROM public.test ORDER BY 1")
		if err != nil {
			t.Fatal(err)
		}
		if len(list) > 0 {
			p = list[0]
		} else {
			t.Fatal("Error _02_interface_json. Point 1.")
		}
	} else {
		p, err = s.GetFaceOne(interface_json_mTestType, "SELECT * FROM public.test ORDER BY 1")
		if err != nil {
			t.Fatal(err)
		}
	}

	switch p["col2"].(type) {
	case map[string]interface{}:
		s := p["col2"].(map[string]interface{})["true"]
		switch s.(type) {
		case bool:
			if !s.(bool) {
				t.Fatal("Error _02_interface_json. Point 2. ")
			}
		default:
			t.Fatal("Error _02_interface_json. Point 3.")
		}
	default:
		t.Fatal("Error _02_interface_json. Point 4.")
	}

}
