package xSql

import (
	"testing"
)

func Test_Insert(t *testing.T) {
	_01_Test_Insert(t)
	_11_Test_Insert(t)
	//t.Fatal("test case")
}

func _01_Test_Insert(t *testing.T) {
	/* Old syntax */
	sql, values := Insert("public.mytable").
		Mark("f_name", "=", "John").
		Mark("l_name", "=", "Lennon").
		Mark("*", "RET", "").
		Mark("f_name as n, l_name as l", "RET", "").
		Comp()
	check_result(t, sql, "INSERT INTO public.mytable ( f_name, l_name ) VALUES ( $1, $2 ) RETURNING *, f_name as n, l_name as l", values, 2)
}

func _11_Test_Insert(t *testing.T) {
	/* Old syntax */
	js := map[string]interface{}{"a": 1, "b": "name"}
	insert := Insert("public.mytable").Append(Mark("a", "=", 1)).Append(Array("b", "=", 1, 2, 3)).Append(Array("c", "="))
	insert.Append(TArray("text", "i", "=", "w"))
	insert.Append(Mark("*", "RET", "")).Append(Mark("e", "=", js).Json()).Append(Mark("b as d", "RET", ""))
	sql, values := insert.Comp()
	check_result(t, sql, "INSERT INTO public.mytable  ( a, b, c, i, e) VALUES "+
		" ($1 , ARRAY[ $2, $3, $4 ],  '{}' , ARRAY[ $5 ]::text[], $6 ) RETURNING *, b as d", values, 6)
}
