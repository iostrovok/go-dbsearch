package xSql

import (
	"testing"
)

func Test_Update(t *testing.T) {
	_01_Test_Update(t)
	//t.Fatal("test case")
}

func _01_Test_Update(t *testing.T) {
	where := Logic("AND").Mark("f_name", "=", "John").Mark("l_name", "=", "Lennon").
		Mark("age", "<", 40)

	sql, values := Update("public.mytable").
		Mark("sended", "=", 1).
		Where(where).
		Comp()
	check_result(t, sql, "UPDATE public.mytable SET sended = $4 WHERE (f_name = $1 AND l_name = $2 AND age < $3 )", values, 4)
}
