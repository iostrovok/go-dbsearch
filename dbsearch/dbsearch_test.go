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
}
