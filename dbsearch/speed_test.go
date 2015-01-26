package dbsearch

import (
	"runtime"
	//"log"
	"reflect"
	"strconv"
	"testing"
)

func Benchmark_TestSpeed_01(b *testing.B) {
	runtime.GOMAXPROCS(8)
	b.StopTimer() //stop the performance timer temporarily while doing initialization
	s := init_test_data()
	if s == nil {
		b.Fatal("Benchmark_TestSpeed_02")
	}
	main_speed_test_table(s, b, 40000)

	b.StartTimer() //restart timer
	for i := 0; i < b.N; i++ {
		_select_TestSpeedGet(true, b, s, i*10)
	}
}

func Benchmark_TestSpeed_02(b *testing.B) {
	runtime.GOMAXPROCS(8)
	b.StopTimer() //stop the performance timer temporarily while doing initialization
	s := init_test_data()
	if s == nil {
		b.Fatal("Benchmark_TestSpeed_02")
	}
	main_speed_test_table(s, b, 40000)

	b.StartTimer() //restart timer
	for i := 0; i < b.N; i++ {
		_select_TestSpeedGet(false, b, s, i*10)
	}
}

func Benchmark_TestSpeed_Json_03(b *testing.B) {
	runtime.GOMAXPROCS(8)
	b.StopTimer() //stop the performance timer temporarily while doing initialization
	s := init_test_data()
	if s == nil {
		b.Fatal("Benchmark_TestSpeed_Json_03")
	}
	main_speed_struct_test_table(s, b, 40000)

	b.StartTimer() //restart timer
	for i := 0; i < b.N; i++ {
		_select_TestSpeedJsonGet(false, b, s, i*10)
	}
}

func Benchmark_TestSpeed_Json_04(b *testing.B) {
	runtime.GOMAXPROCS(8)
	b.StopTimer() //stop the performance timer temporarily while doing initialization
	s := init_test_data()
	if s == nil {
		b.Fatal("Benchmark_TestSpeed_Json_04")
	}
	main_speed_struct_test_table(s, b, 40000)

	b.StartTimer() //restart timer
	for i := 0; i < b.N; i++ {
		_select_TestSpeedJsonGet(true, b, s, i*10)
	}
}

func main_speed_struct_test_table(s *Searcher, b *testing.B, count int) {
	sql_create := " CREATE TABLE public.test " +
		"( col1 serial, col2 json, col3 text, col4 integer[], col5 text[] )"

	js := `'{"mail":"weq","top":"up","list":[1,2,3,5,"assadasd",1233.87],"bool_1":true,"bool_2":true,"inner":{"mail":"weq","top":"up","list":[1,2,3,5,"assadasd",1233.87],"bool_1":true,"bool_2":true}}'`
	ar := `'{10,123,123213,-2323,4345,21232131,466856,123123}'`
	txt := `'{Великобритания,UK,"\"United '' Kingdom","UK,United Kingdom of \"Great , Britain\" ` +
		`and Northern Ireland","Соединенное Королевство Великобритании и Северной Ирландии","ВНУТРИ КАВЫКИ \",` +
		`С ЗАПЯТОЙ","\"",1,"1-2: 1\"2\"",NULL,"single slash: \\\" and \\\\\""}'`

	sql_cols := "INSERT INTO public.test(col2, col3, col4, col5) "
	int_vals := []string{
		" VALUES (" + js + "::json," + js + ", " + ar + "::integer[], " + txt + "::text[])",
		" VALUES (" + js + "::json," + js + ", " + ar + "::integer[], " + txt + "::text[])",
		" VALUES (" + js + "::json," + js + ", " + ar + "::integer[], " + txt + "::text[])",
	}

	sql_vals := []string{}
	for count > 0 {
		count--
		sql_vals = append(sql_vals, int_vals...)
	}

	s.Do("DROP TABLE IF EXISTS public.test")
	s.Do(sql_create)
	for _, v := range sql_vals {
		s.Do(sql_cols + v)
	}

	// Warmimg
	_select_TestSpeedJsonGet(false, b, s, 100)
	_select_TestSpeedJsonGet(true, b, s, 100)
}

func main_speed_test_table(s *Searcher, b *testing.B, count int) {
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

	s.Do("DROP TABLE IF EXISTS public.test")
	s.Do(sql_create)
	for _, v := range sql_vals {
		s.Do(sql_cols + v)
	}

	// Warmimg
	_select_TestSpeedGet(false, b, s, 100)
	_select_TestSpeedGet(true, b, s, 100)
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

func _select_TestSpeedGet(is_fork bool, b *testing.B, s *Searcher, N int) {
	p := []speed_01_TestPlace{}
	sql := "SELECT * FROM public.test ORDER BY 1 LIMIT " + strconv.Itoa(N)
	//log.Println(sql)
	if is_fork {
		s.GetFork(speed_01_mTestType, &p, sql)
	} else {
		s.Get(speed_01_mTestType, &p, sql)
	}
	if len(p) != N {
		b.Fatalf("Bad resault for %d\n", N)
	}
}

type speed_01_Json_TestPlace struct {
	Col1 int                    `db:"col1" type:"serial"`
	Col2 map[string]interface{} `db:"col2" type:"json"`
	Col3 map[string]interface{} `db:"col3" type:"json"`
	Col4 []int                  `db:"col4" type:"[]int"`
	Col5 []string               `db:"col5" type:"[]text"`
}

var speed_01_Json_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(speed_01_Json_TestPlace{}),
}

func _select_TestSpeedJsonGet(is_fork bool, b *testing.B, s *Searcher, N int) {
	p := []speed_01_Json_TestPlace{}
	sql := "SELECT * FROM public.test ORDER BY 1 LIMIT " + strconv.Itoa(N)
	//log.Println(sql)
	if is_fork {
		s.GetFork(speed_01_Json_mTestType, &p, sql)
	} else {
		s.GetNoFork(speed_01_Json_mTestType, &p, sql)
	}
	if len(p) != N {
		b.Fatalf("Bad resault for %d\n", N)
	}
}
