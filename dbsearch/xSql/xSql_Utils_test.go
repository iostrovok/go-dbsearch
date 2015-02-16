package xSql

import (
	"log"
	"testing"
)

func Test_Quote(t *testing.T) {
	_001TestQuote(t)
	//t.Fatal("test case")
}

func _001TestQuote(t *testing.T) {
	testSuite := map[interface{}]string{
		"the molecule's structure": "'the molecule''s structure'",
		" I'''am an actor.":        "' I''''''am an actor.'",
		100500:                     "'100500'",
	}

	for k, v := range testSuite {
		q := Quote(k)
		log.Printf("%v : %s => %s\n", k, v, q)
		if v != q {
			t.Fatal("Error func Quote(val interface{}) string")
		}
	}
}
