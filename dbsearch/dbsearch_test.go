package dbsearch

import (
	"log"
	"os"
	"testing"
)

/*
type PlaceTest struct {
	Id         int    `db:"id" type:"int"`
	ParentId   int    `db:"parent_id" type:"int"`
	ParentsIds int    `db:"parents_ids" type:"int"`
	Name       string `db:"name" type:"text"`
}
*/
type TestPlace struct {
	col1 int     `db:"col1" type:"int"`
	col2 string  `db:"col2" type:"text"`
	col3 string  `db:"col3" type:"text"`
	col4 float64 `db:"col4" type:"real"`
	col5 string  `db:"col5" type:"datetime"`
	col6 string  `db:"col6" type:"date"`
	col7 string  `db:"col7" type:"date"`
	col8 string  `db:"col8" type:"real" is_array:"yes"`
}

var mTestType *AllRows = &AllRows{}

func init_test_data(t *testing.T) *Searcher {
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

	log.Println(dsn)

	dbh, err := DBI(2, dsn, true)
	if err != nil {
		t.Fatal(err)
	}

	dbh.Do("DROP TABLE IF EXISTS public.test")

	sql_create := " CREATE TABLE public.test (col1 int, col2 character varying(255), " +
		" col3 text, col4 real, col5 timestamp, " +
		" col6 date, col7 time, col8 real[] ) "
	dbh.Do(sql_create)

	sql_cols := "INSERT INTO test(col1, col2, col3, col4, col5, col6, col7, col8) "
	sql_vals := []string{
		"VALUES (1,   'John',  'Lennon',   9.12313,  '1945-01-01 00:00:00', '1945-01-01', '00:00:00', '{3.56, 3.45}'::real[])",
		"VALUES (22,  'Telok', 'Macar',    -9.12313, '1812-12-23 06:15:15', '1812-12-23', '06:15:15', '{-3.56, -3.45}'::real[])",
		"VALUES (999, 'Harr',  'Jordjjj',  +9.12313, '1763-05-28 12:30:30', '1763-05-28', '12:30:30', '{3.56, 3.45}'::real[])",
		"VALUES (192, 'Mart',  'Smart',    -9.12313, '0454-06-02 18:45:45', '0454-06-02', '18:45:45', '{-3.56, -3.45}'::real[])",
		"VALUES (111, 'Storm', 'Tropical', 9.12313,  '1001-11-02 23:59:59', '1001-12-31', '18:29:30', '{3.56, 3.45}'::real[])",
	}

	for _, v := range sql_vals {
		dbh.Do(sql_cols + v)
	}

	return dbh
}

func Test_(t *testing.T) {
	_01_Array_Int(t)
	_02_Array_Float(t)
	_03_Array_String(t)
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
