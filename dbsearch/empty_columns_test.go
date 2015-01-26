package dbsearch

import (
	"reflect"
	"testing"
)

func Test_EmptyLoad(t *testing.T) {
	s := init_test_data()
	s.SetDebug(false)
	if s != nil {
		_01_empty_load(t, s)
		_02_empty_load(t, s)
	}
	//t.Fatal("Success [no error] test")
}

func empty_load_f_test_table(s *Searcher) {
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
	empty_load test
	Our struct has less fields than table has their
*/
type empty_load_TestPlace struct {
	Col1 int
	Col2 int
	Col3 int
	Col4 int
	Col5 int
	Col6 int
	Col7 string
}

var empty_load_mTestType *AllRows = &AllRows{
	Table: "test",
	SType: reflect.TypeOf(empty_load_TestPlace{}),
}

func _01_empty_load(t *testing.T, s *Searcher) {
	empty_load_f_test_table(s)
	p := []empty_load_TestPlace{}
	s.SetDieOnColsName(false)
	s.Get(empty_load_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}

type empty_load_02_TestPlace struct {
	Col1 int
	Col2 int
	Col5 int
	Col6 int
	Col7 string
}

var empty_load_02_mTestType *AllRows = &AllRows{
	Table: "test",
	SType: reflect.TypeOf(empty_load_02_TestPlace{}),
}

func _02_empty_load(t *testing.T, s *Searcher) {
	empty_load_f_test_table(s)
	p := []empty_load_02_TestPlace{}
	s.SetDieOnColsName(false)
	s.Get(empty_load_02_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}
