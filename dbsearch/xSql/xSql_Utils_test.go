package xSql

import (
	"log"
	"testing"
)

func Test(t *testing.T) {
	_001_Test_Quote(t)
	t.Fatal("test case")
}

func _001_Test_Quote(t *testing.T) {
	testSuite := map[string]interface{}{
		"the molecule's structure": "'the molecule''s structure'",
		" I'''am an actor.":        "' I''''''am an actor.'",
	}

	for k, v := range testSuite {
		q := Quote(k)
		log.Printf("%s : %s => %s\n", k, v, q)
		if v != q {
			t.Fatal("Error func Quote(val interface{}) string")
		}
	}
}
