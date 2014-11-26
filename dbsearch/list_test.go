package dbsearch

import (
	//"github.com/iostrovok/go-iutils/iutils"
	"log"
	"testing"
	//"time"
)

func Test_List(t *testing.T) {
	s := init_test_data(t)
	// col1, col2, col3, col4, col5, col6, col7
	if s != nil {
		mTestType.PreInit(TestPlace{})
		_01_List_Next(t, s)
		_02_List_Last(t, s)
		_03_List_First(t, s)
	}
	//t.Fatal("error test")
}

func _000_Get_List(t *testing.T, s *Searcher) (*List, []interface{}) {
	val := []interface{}{1, 22, 999, 192, 111}

	list, err := s.List(mTestType, "select * from  public.test where col1 IN ($1,$2,$3,$4,$5) ORDER BY col1", val)
	if err != nil {
		log.Fatalln(err)
	}

	if list.Count() != len(val) {
		log.Fatalln("Error. dbsearch func (l *List) Count() int")
	}
	return list, val
}

func _01_List_Next(t *testing.T, s *Searcher) {

	list, val := _000_Get_List(t, s)

	__01_List_Next(t, s, list, val)
	list.Reset()
	__01_List_Next(t, s, list, val)
}

func __01_List_Next(t *testing.T, s *Searcher, list *List, val []interface{}) {

	CountGet := 0
	r := list.Next()
	for r != nil {
		CountGet++
		switch r.Int("col1") {
		case 1, 22, 999, 192, 111:
			// Nothing
		default:
			log.Fatalln("Error. dbsearch func (l *List) Next() *Row")
		}
		r = list.Next()
	}

	if CountGet != len(val) {
		log.Fatalln("Error. dbsearch func (l *List) Next() *Row don't return all rows")
	}
}

func _02_List_Last(t *testing.T, s *Searcher) {

	list, _ := _000_Get_List(t, s)

	if list.Last().Int("col1") != 999 {
		log.Fatalln("Error. dbsearch func (l *List) Last() *Row")
	}
}

func _03_List_First(t *testing.T, s *Searcher) {

	list, _ := _000_Get_List(t, s)

	if list.Fist().Int("col1") != 1 {
		log.Fatalln("Error. dbsearch func (l *List) Last() *Row")
	}
}
