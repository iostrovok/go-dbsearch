package dbsearch

import (
	"log"
	"os"
	"testing"
)

const TEST_TIME_ZONE = "Europe/Berlin"

func init_test_data() *Searcher {
	login := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWD")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	dbname := os.Getenv("DBNAME")
	sslmode := os.Getenv("SSLMODE")

	usr := ""
	if login != "" {
		usr = " user=" + login
		if pass != "" {
			usr += " password=" + pass
		}
	}

	if host == "" {
		host = "127.0.0.1"
	}

	if port == "" {
		port = "5432"
	}

	if dbname == "" {
		dbname = "test"
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	//"host=127.0.0.1 port=5432 user=smarttable password=smarttable dbname=smarttable sslmode=disable"
	dsn := "host=" + host + " port=" + port + usr + " dbname=" + dbname + " sslmode=" + sslmode

	dbh, err := DBI(2, dsn, true)
	if err != nil {
		log.Panicf("%s\n", err)
	}

	dbh.Do("DROP TABLE IF EXISTS public.test")

	sql_create := " CREATE TABLE public.test (col1 int, col2 character varying(255), " +
		" col3 text, col4 real, col5 real[], col6 int[], col7 text[], col8 text, col9 json, col10 json, col11 timestamp, col12 date, col13 time ) "

	time_line := "'2013-01-01 23:23:23', '2013-01-01', '23:23:23'" // timestamp, date, time, interval

	json_one := "'{\"one\":1,\"list\":[1,2,3,4,5]}'"
	json_line := json_one + ", " + json_one + "::json, " + json_one + "::json, " + time_line

	sql_cols := "INSERT INTO test(col1, col2, col3, col4, col5, col6, col7, col8, col9, col10, col11, col12, col13) "
	sql_vals := []string{
		"VALUES (1,   'John',  'Lennon',   9.12313,  '{3.56, 3.45}'::real[], '{10,20,30,40,50}'::int[], '{one,two,three,four}'::text[], " + json_line + ")",
		"VALUES (22,  'Telok', 'Macar',    -9.12313, '{-3.56, -3.45}'::real[], '{20}'::int[], '{one}'::text[], " + json_line + ")",
		"VALUES (999, 'Harr',  'Jordjjj',  +9.12313, '{3.56, 3.45}'::real[], '{30}'::int[], '{one}'::text[], " + json_line + ")",
		"VALUES (192, 'Mart',  'Smart',    -9.12313, '{-3.56, -3.45}'::real[], '{40}'::int[], '{one}'::text[], " + json_line + ")",
		"VALUES (111, 'Storm', 'Tropical', 9.12313,  '{3.56, 3.45}'::real[], '{50}'::int[], '{one}'::text[], " + json_line + ")",
	}

	make_t_table(dbh, sql_create, sql_cols, sql_vals)

	return dbh
}

func make_t_table(dbh *Searcher, sql_create, sql_cols string, sql_vals []string) {
	dbh.Do("DROP TABLE IF EXISTS public.test")

	dbh.Do(sql_create)

	for _, v := range sql_vals {
		dbh.Do(sql_cols + v)
	}
}

func Test_(t *testing.T) {
	_01_Array_Int(t)
	_02_Array_Float(t)
	_03_Array_String(t)
	_04_Array_String(t)
}

func _01_Array_Int(t *testing.T) {
	int_line := []byte("{1,23,45,6,0,2,2323, 32432423 }")
	int_list := parseIntArray(int_line)
	if len(int_list) != 8 || int_list[4] != 0 || int_list[7] != 32432423 {
		t.Fatal("error parseIntArray")
	}
}

func _02_Array_Float(t *testing.T) {
	f_line := []byte("{3.56, -3.45,12331.213215367, -9999.04353245,0,3.676,4.56}")
	list := parseFloat64Array(f_line)
	if len(list) != 7 || list[4] != 0 || list[1] != -3.45 || list[0] != 3.56 {
		t.Fatal("error _02_Array_Float")
	}
}

func _03_Array_String(t *testing.T) {

	text_array := []string{
		`Великобритания`, `UK`, `"United ' Kingdom`,
		`UK,United Kingdom of "Great , Britain" and Northern Ireland`,
		`Соединенное Королевство Великобритании и Северной Ирландии`,
		`ВНУТРИ КАВЫКИ ", С ЗАПЯТОЙ`,
		`"`, `1`, `1-2: 1"2"`, "", `single slash: \" and \\"`,
	}

	text_line := `{Великобритания,UK,"\"United ' Kingdom","UK,United Kingdom of \"Great , Britain\" and Northern Ireland","Соединенное Королевство Великобритании и Северной Ирландии","ВНУТРИ КАВЫКИ \", С ЗАПЯТОЙ","\"",1,"1-2: 1\"2\"",NULL,"single slash: \\\" and \\\\\""}`
	text_list := parseArray(text_line)

	for i := range text_array {
		if text_array[i] != text_list[i] {
			log.Printf("Need: %s GET result: %s\n", text_array[i], text_list[i])
			log.Printf("Need: %q GET result: %q\n", text_array[i], text_list[i])
			t.Fatal("error parseArray")
		}
	}

	text_array2 := []string{
		`Соединенное Королевство Великобритании и Северной Ирландии`,
	}

	text_line2 := "{\"Соединенное Королевство Великобритании и Северной Ирландии\"}"
	text_list2 := parseArray(text_line2)

	for i := range text_array2 {
		if text_array2[i] != text_list2[i] {
			log.Printf("Need: %s GET result: %s\n", text_array2[i], text_list2[i])
			log.Printf("Need: %q GET result: %q\n", text_array2[i], text_list2[i])
			t.Fatal("error parseArray")
		}
	}
}

/*  Check array with empty string */
func _04_Array_String(t *testing.T) {

	text_array := []string{
		``, `UK`, `"`, ``,
	}

	text_line := `{"",UK,"\"",""}`
	text_list := parseArray(text_line)

	for i := range text_array {
		if text_array[i] != text_list[i] {
			log.Printf("Need: %s GET result: %s\n", text_array[i], text_list[i])
			log.Printf("Need: %q GET result: %q\n", text_array[i], text_list[i])
			t.Fatal("error parseArray")
		}
	}
}
