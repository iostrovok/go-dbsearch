package dbsearch

import (
	"log"
	"reflect"
	"runtime"
	"testing"
)

func Test_Exception(t *testing.T) {
	s := init_test_data()
	s.SetDebug(false)
	if s != nil {
		//_01_exception_load_test(t, s)
		//_02_exception_load_test(t, s)
	}
	//t.Fatal("Success [no error] test")
}

/*
	Check SType field
*/
func _01_exception_load_test(t *testing.T, s *Searcher) {

	var autoload_mTestType *AllRows = &AllRows{}

	var err error = nil
	defer func() {
		if err == nil {
			t.Fatalf("No catch error for empty SType\n")
		}
	}()

	if err := s.PreInit(autoload_mTestType); err == nil {
		t.Fatalf("No catch error for empty SType\n")
	}

}

type exception_05_TestPlace struct {
	Col1 int `db:"col1" type:"serial"`
	Col2 int `db:"col2" type:"int"`
}

var exception_05_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(exception_05_TestPlace{}),
}

func _02_exception_load_test(t *testing.T, s *Searcher) {
	runtime.GOMAXPROCS(1)

	sql_create := " CREATE TABLE public.test " +
		"( col1 serial, col2 text )"

	sql_cols := "INSERT INTO public.test(col2) "
	sql_vals := []string{
		" VALUES ('ee')",
		" VALUES ('ee')",
		" VALUES ('ee')",
	}

	make_t_table(s, sql_create, sql_cols, sql_vals)

	p := []exception_05_TestPlace{}
	s.GetFork(exception_05_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	face, err := s.GetFaceFork(exception_05_mTestType, "SELECT * FROM public.test ORDER BY 1")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("\n%#v\n", face)

	//check_face_json(t, face[0], "_21_exception_test")
}
