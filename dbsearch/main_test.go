package dbsearch

import (
	"reflect"
	"testing"
)

func Test_Row(t *testing.T) {
	s := init_test_data(t)
	if s != nil {
		_01_int32(t, s)
		_02_int64(t, s)
		_03_int(t, s)
		_11_bool(t, s)
		_21_string(t, s)
		_31_float64(t, s)
		_51_json(t, s)
		_61_bytea(t, s)
	}
	//t.Fatal("Success [no error] test")
}

func main_f_test_table(s *Searcher) {
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
	int32 test
*/
type int32_TestPlace struct {
	Col1  int32 `db:"col1" type:"int"`
	Col2  int32 `db:"col2" type:"bigint"`
	Col3  int32 `db:"col3" type:"smallint"`
	Col4  int32 `db:"col4" type:"integer"`
	Col5  int32 `db:"col5" type:"serial"`
	Col6  int32 `db:"col6" type:"bigserial"`
	Col7  int32 `db:"col7" type:"text"`
	Col8  int32 `db:"col8" type:"varchar"`
	Col9  int32 `db:"col9" type:"char"`
	Col11 int32 `db:"col11" type:"real"`
	Col12 int32 `db:"col12" type:"double"`
	Col13 int32 `db:"col13" type:"numeric"`
	Col14 int32 `db:"col14" type:"decimal"`
	Col15 int32 `db:"col15" type:"money"`
	Col16 int32 `db:"col16" type:"bool"`
}

var int32_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(int32_TestPlace{}),
}

func _01_int32(t *testing.T, s *Searcher) {
	main_f_test_table(s)
	p := []int32_TestPlace{}
	s.Get(int32_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}

/*
	int64 test
*/
type int64_TestPlace struct {
	Col1  int64 `db:"col1" type:"int"`
	Col2  int64 `db:"col2" type:"bigint"`
	Col3  int64 `db:"col3" type:"smallint"`
	Col4  int64 `db:"col4" type:"integer"`
	Col5  int64 `db:"col5" type:"serial"`
	Col6  int64 `db:"col6" type:"bigserial"`
	Col7  int64 `db:"col7" type:"text"`
	Col8  int64 `db:"col8" type:"varchar"`
	Col9  int64 `db:"col9" type:"char"`
	Col11 int64 `db:"col11" type:"real"`
	Col12 int64 `db:"col12" type:"double"`
	Col13 int64 `db:"col13" type:"numeric"`
	Col14 int64 `db:"col14" type:"decimal"`
	Col15 int64 `db:"col15" type:"money"`
	Col16 int64 `db:"col16" type:"bool"`
}

var int64_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(int64_TestPlace{}),
}

func _02_int64(t *testing.T, s *Searcher) {
	main_f_test_table(s)
	p := []int64_TestPlace{}
	s.Get(int64_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}

/*
	int test
*/
type int_TestPlace struct {
	Col1  int `db:"col1" type:"int"`
	Col2  int `db:"col2" type:"bigint"`
	Col3  int `db:"col3" type:"smallint"`
	Col4  int `db:"col4" type:"integer"`
	Col5  int `db:"col5" type:"serial"`
	Col6  int `db:"col6" type:"bigserial"`
	Col7  int `db:"col7" type:"text"`
	Col8  int `db:"col8" type:"varchar"`
	Col9  int `db:"col9" type:"char"`
	Col11 int `db:"col11" type:"real"`
	Col12 int `db:"col12" type:"double"`
	Col13 int `db:"col13" type:"numeric"`
	Col14 int `db:"col14" type:"decimal"`
	Col15 int `db:"col15" type:"money"`
	Col16 int `db:"col16" type:"bool"`
}

var int_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(int_TestPlace{}),
}

func _03_int(t *testing.T, s *Searcher) {
	main_f_test_table(s)
	p := []int_TestPlace{}
	s.Get(int_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}

/*
	boolean test
*/
type bool_TestPlace struct {
	Col1 int     `db:"col1" type:"int"`
	Col2 string  `db:"col2" type:"boolean"`
	Col3 int64   `db:"col3" type:"boolean"`
	Col4 float32 `db:"col4" type:"boolean"`
	Col5 float64 `db:"col5" type:"bool"`
	Col6 bool    `db:"col6" type:"bool"`
	Col7 string  `db:"col7" type:"bool"`
	Col8 uint8   `db:"col8" type:"bool"`
	Col9 int     `db:"col9" type:"bool"`
}

var bool_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(bool_TestPlace{}),
}

func _11_bool(t *testing.T, s *Searcher) {

	sql_create := " CREATE TABLE public.test " +
		"(col1 int, col2 boolean, col3 boolean, col4 boolean, " +
		"col5 boolean, col6 boolean, col7 boolean, col8 boolean, col9 boolean ) "
	sql_cols := "INSERT INTO test(col1, col2, col3, col4, col5, col6, col7, col8, col9 ) "
	sql_vals := []string{
		"VALUES (1, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, TRUE )",
		"VALUES (2, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, FALSE)",
		"VALUES (3, null, null, null, null, null, null, null, null )", // check null - nil
	}

	make_t_table(s, sql_create, sql_cols, sql_vals)
	p := []bool_TestPlace{}
	s.Get(bool_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}

/*
	string test
*/
type string_TestPlace struct {
	Col1  string `db:"col1" type:"int"`
	Col2  string `db:"col2" type:"bigint"`
	Col3  string `db:"col3" type:"smallint"`
	Col4  string `db:"col4" type:"integer"`
	Col5  string `db:"col5" type:"serial"`
	Col6  string `db:"col6" type:"bigserial"`
	Col7  string `db:"col7" type:"text"`
	Col8  string `db:"col8" type:"varchar"`
	Col9  string `db:"col9" type:"char"`
	Col11 string `db:"col11" type:"real"`
	Col12 string `db:"col12" type:"double"`
	Col13 string `db:"col13" type:"numeric"`
	Col14 string `db:"col14" type:"decimal"`
	Col15 string `db:"col15" type:"money"`
	Col16 string `db:"col16" type:"bool"`
}

var string_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(string_TestPlace{}),
}

func _21_string(t *testing.T, s *Searcher) {
	main_f_test_table(s)
	p := []string_TestPlace{}
	s.Get(string_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}

/*
	float64 test
*/
type float64_TestPlace struct {
	Col1  float64 `db:"col1" type:"int"`
	Col2  float64 `db:"col2" type:"bigint"`
	Col3  float64 `db:"col3" type:"smallint"`
	Col4  float64 `db:"col4" type:"integer"`
	Col5  float64 `db:"col5" type:"serial"`
	Col6  float64 `db:"col6" type:"bigserial"`
	Col7  float64 `db:"col7" type:"text"`
	Col8  float64 `db:"col8" type:"varchar"`
	Col9  float64 `db:"col9" type:"char"`
	Col11 float64 `db:"col11" type:"real"`
	Col12 float64 `db:"col12" type:"double"`
	Col13 float64 `db:"col13" type:"numeric"`
	Col14 float64 `db:"col14" type:"decimal"`
	Col15 float64 `db:"col15" type:"money"`
	Col16 float64 `db:"col16" type:"bool"`
}

var float64_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(float64_TestPlace{}),
}

func _31_float64(t *testing.T, s *Searcher) {
	main_f_test_table(s)
	p := []float64_TestPlace{}
	s.Get(float64_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}

/*
	json & jsonb test
*/
type json_TestPlace struct {
	Col1 int                    `db:"col1" type:"int"`
	Col2 map[string]interface{} `db:"col2" type:"json"`
	Col3 string                 `db:"col3" type:"json"`
	Col4 map[string]interface{} `db:"col4" type:"json"`
	Col5 string                 `db:"col5" type:"json"`
}

var json_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(json_TestPlace{}),
}

func _51_json(t *testing.T, s *Searcher) {

	str := `'{"array":[{"one":1,"two":"two"},{"next":""}],"null":null,"false":false,"true":true,"ru":"Слова пишем слево на право"}'`

	sql_create := " CREATE TABLE public.test " +
		"(col1 int, col2 text, col3 text, col4 json, col5 json ) "
	sql_cols := "INSERT INTO test(col1, col2, col3, col4, col5 ) "
	sql_vals := []string{
		"VALUES (1, " + str + ", " + str + ", " + str + "::json, " + str + "::json )",
		"VALUES (2, null, null, null, null )", // check null - nil
	}

	make_t_table(s, sql_create, sql_cols, sql_vals)
	p := []json_TestPlace{}
	s.Get(json_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}

/*
	bytea test
*/
type bytea_TestPlace struct {
	Col1 int     `db:"col1" type:"int"`
	Col2 string  `db:"col2" type:"bytea"`
	Col3 []byte  `db:"col3" type:"bytea"`
	Col4 []uint8 `db:"col4" type:"bytea"`
}

var bytea_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(bytea_TestPlace{}),
}

func _61_bytea(t *testing.T, s *Searcher) {

	str := `E'\\xDEADBEEF'`

	sql_create := " CREATE TABLE public.test " +
		"(col1 int, col2 bytea, col3 bytea, col4 bytea ) "
	sql_cols := "INSERT INTO test(col1, col2, col3, col4 ) "
	sql_vals := []string{
		"VALUES (1, " + str + "::bytea, " + str + "::bytea, " + str + "::bytea )",
		"VALUES (2, null, null, null )", // check null - nil
	}

	make_t_table(s, sql_create, sql_cols, sql_vals)
	p := []bytea_TestPlace{}
	s.Get(bytea_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}
