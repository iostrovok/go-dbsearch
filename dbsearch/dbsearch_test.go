package dbsearch

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

type PlaceTest struct {
	Id         int    `db:"id" type:"int"`
	ParentId   int    `db:"parent_id" type:"int"`
	ParentsIds int    `db:"parents_ids" type:"int"`
	Name       string `db:"name" type:"text"`
}

var mType1 *AllRows = &AllRows{}
var mType2 *AllRows = &AllRows{}
var mType3 *AllRows = &AllRows{}
var mType4 *AllRows = &AllRows{}

func openTestConnConninfo(conninfo string) (*sql.DB, error) {

	params := map[string]string{
		"PGDATABASE":        "pqgotest",
		"PGSSLMODE":         "disable",
		"PGCONNECT_TIMEOUT": "20",
		"PGUSER":            "pqgotest",
		"PGPASSWORD":        "pqgotest",
	}

	for k, v := range params {
		env := os.Getenv(k)
		if env == "" {
			os.Setenv(k, v)
		}
	}

	return sql.Open("postgres", conninfo)
}

func openTestConn(t *testing.T) *sql.DB {
	conn, err := openTestConnConninfo("")
	if err != nil {
		log.Println("Connection error")
		t.Fatal(err)
	}

	log.Printf("CONNECTION: %v\n", conn)
	return conn
}

func Test_FieldName(t *testing.T) {

	tests := map[string]string{
		"pre_init_super":   "PreInitSuper",
		"_pre_init_super":  "PreInitSuper",
		"pre_init_super_":  "PreInitSuper",
		"_pre_init_super_": "PreInitSuper",
		"pre":              "Pre",
		"_pre_":            "Pre",
	}

	for k, v := range tests {
		if a := _field_name(k); v != a {
			t.Fatal("error Test_FieldName for '" + k + "', result must be '" + v + "', it is '" + a + "'")
		}
	}

}

func Test_PreInitDB(t *testing.T) {

	s := new(Searcher)
	s.db = openTestConn(t)
	s.db.SetMaxOpenConns(2)

	log.Printf("s.db: %v\n", s.db)

	mType1.PreInitDB(s, "public.test")
	log.Printf("Test_PreInitDB: %v\n", mType1)

	if !mType1.Done {
		t.Fatal("error Test_PreInitDB")
	}

	//t.Fatal("error test")
}

func Test_PreInit(t *testing.T) {
	mType2.PreInit(PlaceTest{})
	if !mType2.Done {
		t.Fatal("error PreInit")
	}
	//t.Fatal("error test")
}

func Test_Array(t *testing.T) {

	int_line := []byte("{1,23,45,6,0,2,2323, 32432423 }")
	int_list := parseIntArray(int_line)
	if len(int_list) != 8 || int_list[4] != 0 || int_list[7] != 32432423 {
		t.Fatal("error parseIntArray")
	}

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
			t.Fatal("error parseArray")
		}
	}

	//t.Fatal("error test")
}
