package dbsearch

import (
	"reflect"
	"testing"
)

func Test_Fork(t *testing.T) {
	s := init_test_data(t)
	if s != nil {
		_01_fork_test(t, s)
	}
	t.Fatal("Success [no error] test")
}

func main_fork_test_table(s *Searcher, count int) {
	sql_create := " CREATE TABLE public.test " +
		"(col1 int, col2 bigint, col3 smallint, col4 integer, " +
		"col5 serial, col6 bigserial, col7 text, col8 varchar(50), col9 char(10), " +
		" col11 real, col12 double precision, col13 numeric, col14 decimal, col15 money, col16 boolean  " +
		") "

	sql_cols := "INSERT INTO test(col1, col2, col3, col4, col5, col6, col7, col8, col9, col11, col12, col13, col14, col15, col16 ) "
	int_vals := []string{
		"VALUES (1, 9223372036854775807, 883, 884, 885, 886, '123456789', '123456789', '1234567890', 12.13, 14.15, 16.17, 18.19, 20.21, TRUE )",
		"VALUES (2, -9223372036854775807, -883, -884, -885, -886, '-123456789', '-123456789', '-123456789', -12.13, -14.15, -16.17, -18.19, -20.21, FALSE )",
		"VALUES (3, null, null, null, 0, 0, null, null, null, null, null, null, null, null, null )", // check null - nil
	}

	sql_vals := []string{}
	for count > 0 {
		count--
		sql_vals = append(sql_vals, int_vals...)
	}

	make_t_table(s, sql_create, sql_cols, sql_vals)
}

/*
	int32 test
*/
type fork_01_TestPlace struct {
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

var fork_01_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(fork_01_TestPlace{}),
}

func _01_fork_test(t *testing.T, s *Searcher) {
	main_fork_test_table(s, 1000)
	p := []fork_01_mTestType{}
	s.Get(fork_01_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")
}
