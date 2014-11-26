package xSql

import (
	"testing"
)

func Test_Insert(t *testing.T) {
	_01_Test_Insert(t)
	//t.Fatal("test case")
}

func _01_Test_Insert(t *testing.T) {

	sql, values := Insert("public.mytable").
		Mark("f_name", "=", "John").
		Mark("l_name", "=", "Lennon").
		Comp()
	check_result(t, sql, "UPDATE public.mytable SET sended = $4 WHERE (f_name = $1 AND l_name = $2 AND age < $3 )", values, 4)
}
