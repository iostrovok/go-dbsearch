package dbsearch

import (
	"log"
	"testing"
)

type PlaceTest struct {
	Id         int    `db:"id" type:"int"`
	ParentId   int    `db:"parent_id" type:"int"`
	ParentsIds int    `db:"parents_ids" type:"int"`
	Name       string `db:"name" type:"text"`
}

var mType *AllRows = &AllRows{}

func Test_PreInit(t *testing.T) {
	mType.PreInit(PlaceTest{})
	if !mType.Done {
		t.Fatal("error PreInit")
	}

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
	log.Printf("%s\n", text_line)

	for i := range text_array {
		if text_array[i] != text_list[i] {
			log.Printf("Need: %s GET result: %s\n", text_array[i], text_list[i])
			log.Printf("Need: %q GET result: %q\n", text_array[i], text_list[i])
			t.Fatal("error parseArray")
		}
		log.Printf("SUCCESS Need: %s GET result: %s\n", text_array[i], text_list[i])
		log.Printf("SUCCESS Need: %q GET result: %q\n", text_array[i], text_list[i])
	}

	//t.Fatal("error test")
}
