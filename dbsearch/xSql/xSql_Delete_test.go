package xSql

import (
	"testing"
)

func Test_Delete(t *testing.T) {
	_01_Test_Delete(t)

	//t.Fatal("test case")
}

func _01_Test_Delete(t *testing.T) {
	sql, values := Delete("public.mytable").
		Logic("AND").Mark("f_name", "=", "John").
		Mark("l_name", "=", "Lennon").
		Mark("age", "<", 40).
		Mark("age", ">", 0).
		Comp()

	check_result(t, sql, "DELETE FROM public.mytable WHERE (f_name = $1 AND l_name = $2 AND age < $3 AND age > $4)", values, 4)
}
