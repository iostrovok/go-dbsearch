package xSql

import (
	"testing"
)

func Test_Update(t *testing.T) {
	_01_Test_Update(t)
	_11_Test_Update(t)
	_21_Test_Update(t)
	//t.Fatal("test case")
}

func _01_Test_Update(t *testing.T) {
	where := Logic("AND").Mark("f_name", "=", "John").Mark("l_name", "=", "Lennon").
		Mark("age", "<", 40)

	sql, values := Update("public.mytable").
		Mark("sended", "=", 1).
		Mark("*", "RET", "").
		Mark("f_name as n, l_name as l", "RET", "").
		Where(where).
		Comp()
	check_result(t, sql, "UPDATE public.mytable SET sended = $4 WHERE (f_name = $1 AND l_name = $2 AND age < $3 )  RETURNING *, f_name as n, l_name as l", values, 4)
}

func _11_Test_Update(t *testing.T) {

	js := map[string]interface{}{"a": 1, "b": "name"}

	update := Update("public.mytable").Append(Mark("a", "=", 1)).Append(Array("b", "=", 1, 2, 3)).Append(Array("c", "="))
	update.Append(Mark("*", "RET", "")).Append(Mark("b as d", "RET", ""))
	update.Append(Mark("e", "=", js).Json())

	update_where := NLogic("AND").Append(Func("startdate < now()")).Append(Func("enddate > now()"))
	update_where.Append(Mark("e", "<>", js).Json())
	update.Where(update_where)
	sql, values := update.Comp()
	check_sql := "UPDATE public.mytable SET a = $2, b = ARRAY[$3, $4, $5], c = '{}', e = $6" +
		" WHERE (startdate < now() AND enddate > now() AND e <> $1) RETURNING *, b as d"
	check_result(t, sql, check_sql, values, 6)
}

func _21_Test_Update(t *testing.T) {

	update := Update("public.mytable").Append(Func("a = a + 1")).Func("b = b + 1")

	update_where := NLogic("AND").Mark("a", "=", 1)
	update.Where(update_where)
	sql, values := update.Comp()
	check_sql := " UPDATE public.mytable SET a = a + 1, b = b + 1 WHERE (  a = $1 ) "
	check_result(t, sql, check_sql, values, 1)
}
