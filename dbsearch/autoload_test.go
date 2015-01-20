package dbsearch

import (
	"reflect"
	"testing"
)

func Test_AutoLoad(t *testing.T) {
	s := init_test_data(t)
	s.SetDebug(false)
	if s != nil {
		_00_autoload_test(t, s)
		_01_autoload_test(t, s)
		_02_autoload_test(t, s)
	}
	//t.Fatal("Success [no error] test")
}

func autoload_main_f_test_table(s *Searcher, cols string) {

	sql_create := " CREATE TABLE public.test ( " + cols + ")"

	make_t_table(s, sql_create, "", []string{})
}

func _00_autoload_test(t *testing.T, s *Searcher) {

	cols := "boolean_list boolean[], bigint_list bigint[], smallint_list smallint[], " +
		"integer_list integer[],  text_list text[], varchar_list varchar(50)[], " +
		"char_list char(10)[], real_list real[], double_precision_list double precision[], " +
		"numeric_list numeric[], decimal_list decimal[], money_list money[], " +
		"date_s date, time_s time, timestamp_s timestamp, " +
		"timestamp_tz_s timestamp with time zone, " +
		"int_s int, bigint_s bigint, smallint_s smallint, integer_s integer, " +
		"serial_s serial, bigserial_s bigserial, text_s text, varchar_s varchar(50), char_s char(10), " +
		"real_s real, double_s double precision, numeric_s numeric, " +
		"decimal_s decimal, money_s money, boolean_s boolean,  " +
		"json_s json "

	autoload_main_f_test_table(s, cols)

	Table := OneTableInfo{
		Table:  "test",
		Schema: "public",
	}

	if err := s.GetTableData(&Table); err != nil {
		t.Fatal(err)
	}

	if len(Table.Rows) == 0 {
		t.Fatalf("Error read fields\n")
	}

	if Table.Rows["money_s"].Type != "money" || Table.Rows["money_s"].Field != "MoneyS" {
		t.Fatalf("Error read fields money_s\n")
	}

	if Table.Rows["text_s"].Type != "text" || Table.Rows["text_s"].Field != "TextS" {
		t.Fatalf("Error read fields text_s\n")
	}

	if Table.Rows["double_precision_list"].Type != "[]double" || Table.Rows["double_precision_list"].Field != "DoublePrecisionList" {
		t.Fatalf("Error read fields double_precision_list\n")
	}

	if Table.Rows["real_s"].Type != "real" || Table.Rows["real_s"].Field != "RealS" {
		t.Fatalf("Error read fields real_s\n")
	}
	if Table.Rows["json_s"].Type != "json" || Table.Rows["json_s"].Field != "JsonS" {
		t.Fatalf("Error read fields json_s\n")
	}
}

type autoload_5_TestPlace struct {
	Col1 string `db:"col1"`
	Col2 string `db:"col2"`
	Col3 string `db:"col3"`
	Col4 string `db:"col4"`
	Col5 string `db:"col5"`
	Col6 string `db:"col6"`
}

type autoload_12_TestPlace struct {
	Col1  string `db:"col1"`
	Col2  string `db:"col2"`
	Col3  string `db:"col3"`
	Col4  string `db:"col4"`
	Col5  string `db:"col5"`
	Col6  string `db:"col6"`
	Col7  string `db:"col7"`
	Col8  string `db:"col8"`
	Col9  string `db:"col9"`
	Col10 string `db:"col10"`
	Col11 string `db:"col11"`
	Col12 string `db:"col12"`
}

func _01_autoload_test(t *testing.T, s *Searcher) {
	var autoload_mTestType *AllRows = &AllRows{
		Table:  "test",
		Schema: "public",
		SType:  reflect.TypeOf(autoload_5_TestPlace{}),
	}

	cols := "col1 date, col2 time, col3 timestamp, " +
		"col4 timestamp with time zone, col5 timestamp with time zone, " +
		"col6 timestamp with time zone, col7 timestamp with time zone "

	autoload_main_f_test_table(s, cols)
	s.PreInit(autoload_mTestType)
}

type autoload_02_TestPlace struct {
	Col1 int
	Col2 string
	Col3 string
	Col4 []string
	Col5 []string
	Col6 []string
}

func _02_autoload_test(t *testing.T, s *Searcher) {
	var autoload_mTestType *AllRows = &AllRows{
		Table:  "test",
		Schema: "public",
		SType:  reflect.TypeOf(autoload_02_TestPlace{}),
	}

	cols := "col1 date, col2 time, col3 int, " +
		"col4 smallint[], col5 text[], " +
		"col6 bigint[], col7 char(100)[] "

	autoload_main_f_test_table(s, cols)
	s.PreInit(autoload_mTestType)
}
