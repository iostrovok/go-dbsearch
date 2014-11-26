package dbsearch

import (
	"github.com/iostrovok/go-iutils/iutils"
	"log"
	"testing"
	"time"
)

var date map[string]interface{} = map[string]interface{}{
	"col1": 1,
	"col2": "text",
	"col3": 3.0,
	"col4": "golova",
	"col5": "golova",
	"col6": "golova",
	"col7": "golova",
	"col8": []float64{6.66, 45.6},
}

var cols []string = []string{"col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8"}

var rows map[string]bool = map[string]bool{
	"col1": true,
	"col2": true,
	"col3": true,
	"col4": true,
	"col5": true,
	"col6": true,
	"col7": true,
	"col8": true,
}

func Test_Row(t *testing.T) {
	s := init_test_data(t)
	// col1, col2, col3, col4, col5, col6, col7
	if s != nil {
		mTestType.PreInit(TestPlace{})
		_01_Set(t, s)
		_02_Str(t, s)
		_03_Int(t, s)
		_04_Float(t, s)
		_05_Date(t, s)
		_06_DateTime(t, s)
		_07_Time(t, s)
	}
	//t.Fatal("error test")
}

func _01_Set(t *testing.T, s *Searcher) {
	r := s.SetRow(date, cols, rows)

	if r.IsEmpty() {
		log.Fatalln("Error. dbsearch func (r *Row) IsEmpty() bool")
	}
	c := r.Cols()
	if 8 != len(c) {
		log.Fatalln("Error. dbsearch func (r *Row) Cols() []string")
	}

	i := r.Interface()
	if 8 != len(i) {
		log.Fatalln("Error. dbsearch func (r *Row) Interface() map[string]interface{}")
	}
}

func _02_Str(t *testing.T, s *Searcher) {
	r, err := s.One(mTestType, "select * from  public.test where col1 = $1", []interface{}{1})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("_02_Str: %s, %s\n", r.Str("col2"), r.Str("col3"))

	if r.Str("col2") != "John" || r.Str("col3") != "Lennon" {
		log.Fatalln("Error. dbsearch func (r *Row) Str(name string) string")
	}
	//log.Fatalln("Error. Test error.")
}

func _03_Int(t *testing.T, s *Searcher) {
	r, err := s.One(mTestType, "select * from  public.test where col1 = $1", []interface{}{1})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("_03_Int: %d\n", r.Int("col1"))

	if r.Int("col1") != 1 {
		log.Fatalln("Error. dbsearch func (r *Row) Int(name string) int")
	}
	//log.Fatalln("Error. Test error.")
}

func _04_Float(t *testing.T, s *Searcher) {
	r, err := s.One(mTestType, "select * from  public.test where col1 = $1", []interface{}{1})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("_04_Float: %f\n", r.Float("col4"))

	if !iutils.EqFloat64(r.Float("col4"), 9.12313, 5) {
		log.Fatalln("Error. dbsearch func (r *Row) Float(name string) float64")
	}
	//log.Fatalln("Error. Test error.")
}

func _05_Date(t *testing.T, s *Searcher) {
	r, err := s.One(mTestType, "select * from  public.test where col1 = $1", []interface{}{111})
	if err != nil {
		log.Fatalln(err)
	}

	//out := *r.Date("col6")
	it := time.Date(1001, time.December, 31, 0, 0, 0, 0, time.UTC)
	if !it.Equal(*r.Date("col6")) {
		log.Fatalln("Error. dbsearch func (r *Row) DateTime(name string) time.Time")
	}
	//log.Fatalln("Error. Test error.")
}

func _06_DateTime(t *testing.T, s *Searcher) {

	r, err := s.One(mTestType, "select * from  public.test where col1 = $1", []interface{}{111})
	if err != nil {
		log.Fatalln(err)
	}

	it := time.Date(1001, time.November, 02, 23, 59, 59, 0, time.UTC)
	if !it.Equal(*r.DateTime("col5")) {
		log.Fatalln("Error. dbsearch func (r *Row) DateTime(name string) time.Time")
	}
	//log.Fatalln("Error. Test error.")
}

func _07_Time(t *testing.T, s *Searcher) {

	r, err := s.One(mTestType, "select * from  public.test where col1 = $1", []interface{}{111})
	if err != nil {
		log.Fatalln(err)
	}

	it := time.Date(0, 1, 1, 18, 29, 30, 0, time.UTC)
	if !it.Equal(*r.Time("col7")) {
		log.Fatalln("Error. dbsearch func (r *Row) DateTime(name string) time.Time")
	}
	//log.Fatalln("Error. Test error.")
}
