package dbsearch

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

const TEST_TIME_ZONE = "Europe/Berlin"

func Test_DateTime(t *testing.T) {
	s := init_test_data(t)
	if s != nil {
		_01_datetime_string(t, s)
		_11_datetime_int(t, s)
		_12_datetime_int64(t, s)
		_21_datetime_time(t, s)
		_31_datetime_float64(t, s)
		_41_datetime_uint(t, s)
		_51_datetime_map(t, s)
		_61_datetime_intlist(t, s)
	}
	//t.Fatal("Success [no error] test")
}

func array_main_f_test_table(s *Searcher) {
	sql_create := " CREATE TABLE public.test ( " +
		"col1 int, col2 date, col3 time, col4 timestamp, " +
		"col5 timestamp with time zone, col6 timestamp with time zone " +
		")"
	sql_cols := "INSERT INTO test( col1, col2, col3, col4, col5, col6 ) "

	loc, _ := time.LoadLocation(TEST_TIME_ZONE)
	t := time.Date(2014, time.January, 1, 1, 1, 0, 0, loc)

	_, n := t.Zone()

	d1 := "'2014-11-12 01:22:12+00'"
	dt := fmt.Sprintf("'2014-11-12 01:22:12 +0%d:00'", n/60/60)

	sql_vals := []string{
		" VALUES (1, " + d1 + ", " + d1 + ", " + d1 + ", " + dt + ", " + dt + " ) ",
		" VALUES (2, " + d1 + ", " + d1 + ", " + d1 + ", " + dt + ", " + dt + " ) ",
		" VALUES (3, null, null, null, null, null ) ",
	}

	s.Do(fmt.Sprintf("SET TIME ZONE '%s'", TEST_TIME_ZONE))
	make_t_table(s, sql_create, sql_cols, sql_vals)
}

/*
	string test
*/

type time_string_TestPlace struct {
	Col1 int    `db:"col1" type:"int"`
	Col2 string `db:"col2" type:"date"`
	Col3 string `db:"col3" type:"time"`
	Col4 string `db:"col4" type:"timestamp"`
	Col5 string `db:"col5" type:"date"`
	Col6 string `db:"col6" type:"timestamp"`
}

var time_string_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(time_string_TestPlace{}),
}

func _01_datetime_string(t *testing.T, s *Searcher) {
	array_main_f_test_table(s)
	p := []time_string_TestPlace{}
	s.Get(time_string_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	if p[0].Col2 != "2014-11-12 00:00:00 +0000" {
		t.Fatal("Error _01_datetime_string.Col2; string <= date [timestamp without time zone]")
	}
	if p[0].Col3 != "0000-01-01 01:22:12 +0000" {
		t.Fatal("Error _01_datetime_string.Col3; string <= datetime [timestamp without time zone]")
	}
	if p[0].Col4 != "2014-11-12 01:22:12 +0000" {
		t.Fatal("Error _01_datetime_string.Col4; string <= time [timestamp without time zone]")
	}
	if p[0].Col5 != "2014-11-12 00:00:00 +0100" {
		t.Fatal("Error _01_datetime_string.Col5; string <= date [timestamp with time zone]")
	}
	if p[0].Col6 != "2014-11-12 01:22:12 +0100" {
		t.Fatal("Error _01_datetime_string.Col6; string <= timestamp [timestamp with time zone]")
	}
}

/*
	int test
*/

type date_int_TestPlace struct {
	Col1 int `db:"col1" type:"int"`
	Col2 int `db:"col2" type:"date"`
	Col3 int `db:"col3" type:"time"`
	Col4 int `db:"col4" type:"timestamp"`
	Col5 int `db:"col5" type:"date"`
	Col6 int `db:"col6" type:"timestamp"`
}

var date_int_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(date_int_TestPlace{}),
}

func _11_datetime_int(t *testing.T, s *Searcher) {
	array_main_f_test_table(s)
	p := []date_int_TestPlace{}
	s.Get(date_int_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	if p[0].Col2 != 1415750400 {
		t.Fatal("Error _11_datetime_int.Col2; int <= date [timestamp without time zone]")
	}
	if p[0].Col3 != -62167214268 {
		t.Fatal("Error _11_datetime_int.Col3; int <= datetime [timestamp without time zone]")
	}
	if p[0].Col4 != 1415755332 {
		t.Fatal("Error _11_datetime_int.Col4; int <= time [timestamp without time zone]")
	}
	if p[0].Col5 != 1415746800 {
		t.Fatal("Error _11_datetime_int.Col5; int <= date [timestamp with time zone]")
	}
	if p[0].Col6 != 1415751732 {
		t.Fatal("Error _11_datetime_int.Col6; int <= timestamp [timestamp with time zone]")
	}
}

type date_int64_TestPlace struct {
	Col1 int   `db:"col1" type:"int"`
	Col2 int64 `db:"col2" type:"date"`
	Col3 int64 `db:"col3" type:"time"`
	Col4 int64 `db:"col4" type:"timestamp"`
	Col5 int64 `db:"col5" type:"date"`
	Col6 int64 `db:"col6" type:"timestamp"`
}

var date_int64_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(date_int64_TestPlace{}),
}

func _12_datetime_int64(t *testing.T, s *Searcher) {
	array_main_f_test_table(s)
	p := []date_int64_TestPlace{}
	s.Get(date_int64_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	if p[0].Col2 != 1415750400 {
		t.Fatal("Error _12_datetime_int64.Col2; int64 <= date [timestamp without time zone]")
	}
	if p[0].Col3 != -62167214268 {
		t.Fatal("Error _12_datetime_int64.Col3; int64 <= datetime [timestamp without time zone]")
	}
	if p[0].Col4 != 1415755332 {
		t.Fatal("Error _12_datetime_int64.Col4; int64 <= time [timestamp without time zone]")
	}
	if p[0].Col5 != 1415746800 {
		t.Fatal("Error _12_datetime_int64.Col5; int64 <= date [timestamp with time zone]")
	}
	if p[0].Col6 != 1415751732 {
		t.Fatal("Error _12_datetime_int64.Col6; int64 <= timestamp [timestamp with time zone]")
	}
}

type date_time_TestPlace struct {
	Col1 int       `db:"col1" type:"int"`
	Col2 time.Time `db:"col2" type:"date"`
	Col3 time.Time `db:"col3" type:"time"`
	Col4 time.Time `db:"col4" type:"timestamp"`
	Col5 time.Time `db:"col5" type:"date"`
	Col6 time.Time `db:"col6" type:"timestamp"`
}

var date_time_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(date_time_TestPlace{}),
}

func _21_datetime_time(t *testing.T, s *Searcher) {
	array_main_f_test_table(s)
	p := []date_time_TestPlace{}

	s.Get(date_time_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	loc, _ := time.LoadLocation(TEST_TIME_ZONE)
	c2 := time.Date(2014, time.November, 12, 0, 0, 0, 0, time.UTC)
	c3 := time.Date(0, time.January, 1, 1, 22, 12, 0, time.UTC)
	c4 := time.Date(2014, time.November, 12, 1, 22, 12, 0, time.UTC)
	c5 := time.Date(2014, time.November, 12, 0, 0, 0, 0, loc)
	c6 := time.Date(2014, time.November, 12, 1, 22, 12, 0, loc)

	if !p[0].Col2.Equal(c2) {
		t.Fatal("Error _21_datetime_time.Col2; time <= date [timestamp without time zone]")
	}
	if !p[0].Col3.Equal(c3) {
		t.Fatal("Error _21_datetime_time.Col3; time <= datetime [timestamp without time zone]")
	}
	if !p[0].Col4.Equal(c4) {
		t.Fatal("Error _21_datetime_time.Col4; time <= time [timestamp without time zone]")
	}
	if !p[0].Col5.Equal(c5) {
		t.Fatal("Error _21_datetime_time.Col5; time <= date [timestamp with time zone]")
	}
	if !p[0].Col6.Equal(c6) {
		t.Fatal("Error _21_datetime_time.Col6; time <= timestamp [timestamp with time zone]")
	}
}

type time_uint_TestPlace struct {
	Col1 int  `db:"col1" type:"int"`
	Col2 uint `db:"col2" type:"date"`
	Col3 uint `db:"col3" type:"time"`
	Col4 uint `db:"col4" type:"timestamp"`
	Col5 uint `db:"col5" type:"date"`
	Col6 uint `db:"col6" type:"timestamp"`
}

var time_uint_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(time_uint_TestPlace{}),
}

func _41_datetime_uint(t *testing.T, s *Searcher) {
	array_main_f_test_table(s)
	p := []time_uint_TestPlace{}
	s.Get(time_uint_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	if p[0].Col2 != 1415750400 {
		t.Fatal("Error _41_datetime_uint.Col2; uint <= date [timestamp without time zone]")
	}
	if p[0].Col3 != 0 {
		t.Fatal("Error _41_datetime_uint.Col3; uint <= datetime [timestamp without time zone]")
	}
	if p[0].Col4 != 1415755332 {
		t.Fatal("Error _41_datetime_uint.Col4; uint <= time [timestamp without time zone]")
	}
	if p[0].Col5 != 1415746800 {
		t.Fatal("Error _41_datetime_uint.Col5; uint <= date [timestamp with time zone]")
	}
	if p[0].Col6 != 1415751732 {
		t.Fatal("Error _41_datetime_uint.Col6; uint <= timestamp [timestamp with time zone]")
	}
}

type time_float64_TestPlace struct {
	Col1 int     `db:"col1" type:"int"`
	Col2 float64 `db:"col2" type:"date"`
	Col3 float64 `db:"col3" type:"time"`
	Col4 float64 `db:"col4" type:"timestamp"`
	Col5 float64 `db:"col5" type:"date"`
	Col6 float64 `db:"col6" type:"timestamp"`
}

var time_float64_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(time_float64_TestPlace{}),
}

func _31_datetime_float64(t *testing.T, s *Searcher) {
	array_main_f_test_table(s)
	p := []time_float64_TestPlace{}
	s.Get(time_float64_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	if fmt.Sprintf("%.0f", p[0].Col2) != "1415750400000000" {
		t.Fatal("Error _01_datetime_float64.Col2; float64 <= date [timestamp without time zone]")
	}
	if fmt.Sprintf("%.0f", p[0].Col3) != "-6826982046871345" {
		t.Fatal("Error _01_datetime_float64.Col3; float64 <= datetime [timestamp without time zone]")
	}
	if fmt.Sprintf("%.0f", p[0].Col4) != "1415755332000000" {
		t.Fatal("Error _01_datetime_float64.Col4; float64 <= time [timestamp without time zone]")
	}
	if fmt.Sprintf("%.0f", p[0].Col5) != "1415746800000000" {
		t.Fatal("Error _01_datetime_float64.Col5; float64 <= date [timestamp with time zone]")
	}
	if fmt.Sprintf("%.0f", p[0].Col6) != "1415751732000000" {
		t.Fatal("Error _01_datetime_float64.Col6; float64 <= timestamp [timestamp with time zone]")
	}
}

type time_map_TestPlace struct {
	Col1 int            `db:"col1" type:"int"`
	Col2 map[string]int `db:"col2" type:"date"`
	Col3 map[string]int `db:"col3" type:"time"`
	Col4 map[string]int `db:"col4" type:"timestamp"`
	Col5 map[string]int `db:"col5" type:"date"`
	Col6 map[string]int `db:"col6" type:"timestamp"`
}

var time_map_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(time_map_TestPlace{}),
}

func _51_datetime_map(t *testing.T, s *Searcher) {
	array_main_f_test_table(s)
	p := []time_map_TestPlace{}
	s.Get(time_map_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	if p[0].Col2["zone"] != 0 || p[0].Col2["year"] != 2014 || p[0].Col2["month"] != 11 || p[0].Col2["day"] != 12 {
		t.Fatal("Error _51_datetime_map.Col2  [a]; map[string]int <= date [timestamp without time zone]")
	}
	if p[0].Col2["hour"] != 0 || p[0].Col2["minute"] != 0 || p[0].Col2["second"] != 0 {
		t.Fatal("Error _51_datetime_map.Col2 [b]; map[string]int <= date [timestamp without time zone]")
	}

	if p[0].Col3["zone"] != 0 || p[0].Col3["year"] != 0 || p[0].Col3["month"] != 1 || p[0].Col3["day"] != 1 {
		t.Fatal("Error _51_datetime_map.Col3 [a]; map[string]int <= date [timestamp without time zone]")
	}
	if p[0].Col3["hour"] != 1 || p[0].Col3["minute"] != 22 || p[0].Col3["second"] != 12 {
		t.Fatal("Error _51_datetime_map.Col3 [b]; map[string]int <= date [timestamp without time zone]")
	}

	if p[0].Col4["zone"] != 0 || p[0].Col4["year"] != 2014 || p[0].Col4["month"] != 11 || p[0].Col4["day"] != 12 {
		t.Fatal("Error _51_datetime_map.Col4 [a]; map[string]int <= date [timestamp without time zone]")
	}
	if p[0].Col4["hour"] != 1 || p[0].Col4["minute"] != 22 || p[0].Col4["second"] != 12 {
		t.Fatal("Error _51_datetime_map.Col4 [b]; map[string]int <= date [timestamp without time zone]")
	}

	if p[0].Col5["zone"] != 3600 || p[0].Col5["year"] != 2014 || p[0].Col5["month"] != 11 || p[0].Col5["day"] != 12 {
		t.Fatal("Error _51_datetime_map.Col5 [a]; map[string]int <= date [timestamp without time zone]")
	}
	if p[0].Col5["hour"] != 0 || p[0].Col5["minute"] != 0 || p[0].Col5["second"] != 0 {
		t.Fatal("Error _51_datetime_map.Col5 [b]; map[string]int <= date [timestamp without time zone]")
	}

	if p[0].Col6["zone"] != 3600 || p[0].Col6["year"] != 2014 || p[0].Col6["month"] != 11 || p[0].Col6["day"] != 12 {
		t.Fatal("Error _51_datetime_map.Col6 [a]; map[string]int <= date [timestamp without time zone]")
	}
	if p[0].Col6["hour"] != 1 || p[0].Col6["minute"] != 22 || p[0].Col6["second"] != 12 {
		t.Fatal("Error _51_datetime_map.Col6 [b]; map[string]int <= date [timestamp without time zone]")
	}
}

type time_intlist_TestPlace struct {
	Col1 int   `db:"col1" type:"int"`
	Col2 []int `db:"col2" type:"date"`
	Col3 []int `db:"col3" type:"time"`
	Col4 []int `db:"col4" type:"timestamp"`
	Col5 []int `db:"col5" type:"date"`
	Col6 []int `db:"col6" type:"timestamp"`
}

var time_intlist_mTestType *AllRows = &AllRows{
	SType: reflect.TypeOf(time_intlist_TestPlace{}),
}

func _61_datetime_intlist(t *testing.T, s *Searcher) {
	array_main_f_test_table(s)
	p := []time_intlist_TestPlace{}
	s.Get(time_intlist_mTestType, &p, "SELECT * FROM public.test ORDER BY 1")

	if p[0].Col2[7] != 0 || p[0].Col2[0] != 2014 || p[0].Col2[1] != 11 || p[0].Col2[2] != 12 {
		t.Fatal("Error _61_datetime_intlist.Col2  [a]; []int <= date [timestamp without time zone]")
	}
	if p[0].Col2[3] != 0 || p[0].Col2[4] != 0 || p[0].Col2[5] != 0 {
		t.Fatal("Error _61_datetime_intlist.Col2 [b]; []int <= date [timestamp without time zone]")
	}

	if p[0].Col3[7] != 0 || p[0].Col3[0] != 0 || p[0].Col3[1] != 1 || p[0].Col3[2] != 1 {
		t.Fatal("Error _61_datetime_intlist.Col3 [a]; []int <= date [timestamp without time zone]")
	}
	if p[0].Col3[3] != 1 || p[0].Col3[4] != 22 || p[0].Col3[5] != 12 {
		t.Fatal("Error _61_datetime_intlist.Col3 [b]; []int <= date [timestamp without time zone]")
	}

	if p[0].Col4[7] != 0 || p[0].Col4[0] != 2014 || p[0].Col4[1] != 11 || p[0].Col4[2] != 12 {
		t.Fatal("Error _61_datetime_intlist.Col4 [a]; []int <= date [timestamp without time zone]")
	}
	if p[0].Col4[3] != 1 || p[0].Col4[4] != 22 || p[0].Col4[5] != 12 {
		t.Fatal("Error _61_datetime_intlist.Col4 [b]; []int <= date [timestamp without time zone]")
	}

	if p[0].Col5[7] != 3600 || p[0].Col5[0] != 2014 || p[0].Col5[1] != 11 || p[0].Col5[2] != 12 {
		t.Fatal("Error _61_datetime_intlist.Col5 [a]; []int <= date [timestamp without time zone]")
	}
	if p[0].Col5[3] != 0 || p[0].Col5[4] != 0 || p[0].Col5[5] != 0 {
		t.Fatal("Error _61_datetime_intlist.Col5 [b]; []int <= date [timestamp without time zone]")
	}

	if p[0].Col6[7] != 3600 || p[0].Col6[0] != 2014 || p[0].Col6[1] != 11 || p[0].Col6[2] != 12 {
		t.Fatal("Error _61_datetime_intlist.Col6 [a]; []int <= date [timestamp without time zone]")
	}
	if p[0].Col6[3] != 1 || p[0].Col6[4] != 22 || p[0].Col6[5] != 12 {
		t.Fatal("Error _61_datetime_intlist.Col6 [b]; []int <= date [timestamp without time zone]")
	}
}
