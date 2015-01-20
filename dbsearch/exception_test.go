package dbsearch

import (
	"testing"
)

func Test_Exception(t *testing.T) {
	s := init_test_data(t)
	s.SetDebug(false)
	if s != nil {
		//_01_exception_load_test(t, s)
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
