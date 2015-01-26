package dbsearch

import (
	"log"
	"reflect"
	"strconv"
	"testing"
)

func Benchmark_TestSpeed(b *testing.B) {
	s := init_test_data()
	if s != nil {
		_01_TestSpeedGet(b, s)
	}
	//b.Fatal("Success [no error] test")
}

func main_speed_test_table(s *Searcher, count int) {
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
type speed_01_TestPlace struct {
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

var speed_01_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(speed_01_TestPlace{}),
}

func _select_TestSpeedGet(b *testing.B, s *Searcher, N int) { //benchmark function starts with "Benchmark" and takes a pointer to type testing.B
	// Warmimg
	p := []speed_01_TestPlace{}
	sql := "SELECT * FROM public.test ORDER BY 1 LIMIT " + strconv.Itoa(N)
	log.Println(sql)
	s.Get(speed_01_mTestType, &p, sql)
	if len(p) != N {
		b.Fatalf("Bad resault for %d\n", N)
	}
}

func _01_TestSpeedGet(b *testing.B, s *Searcher) { //benchmark function starts with "Benchmark" and takes a pointer to type testing.B
	b.StopTimer() //stop the performance timer temporarily while doing initialization
	log.Printf("_01_TestSpeedGet:: %d\n", 1)

	main_speed_test_table(s, 10000)
	log.Printf("_01_TestSpeedGet:: %d\n", 2)

	// Warmimg
	_select_TestSpeedGet(b, s, 4)
	log.Printf("_01_TestSpeedGet:: %d\n", 3)
	_select_TestSpeedGet(b, s, 5)
	log.Printf("_01_TestSpeedGet:: %d\n", 1000)

	b.StartTimer() //restart timer
	log.Printf("_01_TestSpeedGet:: %d\n", 5)
	for i := 0; i < b.N; i++ {
		log.Printf("i:: %d\n", i)
		_select_TestSpeedGet(b, s, i*10)
	}
}
