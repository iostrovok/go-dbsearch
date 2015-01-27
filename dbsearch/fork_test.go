package dbsearch

import (
	"reflect"
	"runtime"
	"testing"
)

func Test_Fork(t *testing.T) {
	s := init_test_data()
	if s != nil {
		_01_fork_test(t, s)
		_02_fork_test(t, s)
		_03_fork_test(t, s)
		_04_fork_test(t, s)
		_05_fork_test(t, s)
		_21_fork_test(t, s)
		_22_fork_test(t, s)
		_23_fork_test(t, s)
		_24_fork_test(t, s)
		_25_fork_test(t, s)
	}
	//t.Fatal("Success [no error] test")
}

func check_face(t *testing.T, face map[string]interface{}, line string) {
	if _AnyToString(face["col2"]) != "9223372036854775807" {
		t.Fatalf("%s\n", line)
	}
}

func check_face_json(t *testing.T, face map[string]interface{}, line string) {

	switch face["col2"].(type) {
	case map[string]interface{}:
		r := face["col2"].(map[string]interface{})
		if _AnyToString(r["mail"]) != "weq" {
			t.Fatalf("%s\n", line)
		}
	default:
		t.Fatalf("Bad type: %T\n", face["col2"])
	}
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

func main_fork_json_test_table(s *Searcher, count int) {
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

	make_t_table(s, sql_create, sql_cols, sql_vals)
}

type fork_01_TestPlace struct {
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

var fork_01_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(fork_01_TestPlace{}),
}

func _01_fork_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(1)
	main_fork_test_table(s, 10)
	p := []fork_01_TestPlace{}
	s.GetFork(fork_01_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFaceFork(fork_01_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	check_face(t, face[0], "_01_fork_test")
}

func _02_fork_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(8)
	main_fork_test_table(s, 10)
	p := []fork_01_TestPlace{}
	s.GetFork(fork_01_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFaceFork(fork_01_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	check_face(t, face[0], "_02_fork_test")

}

func _03_fork_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(1)
	main_fork_test_table(s, 10)
	p := []fork_01_TestPlace{}
	s.Get(fork_01_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFace(fork_01_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	check_face(t, face[0], "_03_fork_test")

}

func _04_fork_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(8)
	main_fork_test_table(s, 10)
	p := []fork_01_TestPlace{}
	s.Get(fork_01_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFace(fork_01_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	check_face(t, face[0], "_04_fork_test")

}

func _05_fork_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(8)
	main_fork_test_table(s, 10)
	p := []fork_01_TestPlace{}
	s.GetNoFork(fork_01_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFaceNoFork(fork_01_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	check_face(t, face[0], "_05_fork_test")

}

/*----------------------------------------------------------------------------*/
/*
	int64 test
*/
type fork_05_TestPlace struct {
	Col1 int                    `db:"col1" type:"serial"`
	Col2 map[string]interface{} `db:"col2" type:"json"`
	Col3 map[string]interface{} `db:"col3" type:"json"`
	Col4 []int                  `db:"col4" type:"[]int"`
	Col5 []string               `db:"col5" type:"[]text"`
}

var fork_05_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(fork_05_TestPlace{}),
}

func _21_fork_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(1)
	main_fork_json_test_table(s, 10)
	p := []fork_05_TestPlace{}
	s.GetFork(fork_05_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFaceFork(fork_05_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	check_face_json(t, face[0], "_21_fork_test")
}

func _22_fork_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(8)
	main_fork_json_test_table(s, 10)
	p := []fork_05_TestPlace{}
	s.GetFork(fork_05_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFaceFork(fork_05_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	check_face_json(t, face[0], "_22_fork_test")
}

func _23_fork_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(1)
	main_fork_json_test_table(s, 10)
	p := []fork_05_TestPlace{}
	s.Get(fork_05_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFace(fork_05_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	check_face_json(t, face[0], "_23_fork_test")
}

func _24_fork_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(8)
	main_fork_json_test_table(s, 10)
	p := []fork_05_TestPlace{}
	s.Get(fork_05_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFace(fork_05_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	check_face_json(t, face[0], "_24_fork_test")
}

func _25_fork_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(8)
	main_fork_json_test_table(s, 10)
	p := []fork_05_TestPlace{}
	s.GetNoFork(fork_05_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFaceNoFork(fork_05_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	check_face_json(t, face[0], "_25_fork_test")
}
