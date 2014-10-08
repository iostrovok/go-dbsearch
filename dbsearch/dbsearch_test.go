package dbsearch

import (
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

	test_line := []byte("{1,23,45,6,0,2,2323, 32432423 }")
	test_list := parseIntArray(test_line)
	if len(test_list) != 8 || test_list[4] != 0 || test_list[7] != 32432423 {
		t.Fatal("error parseIntArray")
	}
}
