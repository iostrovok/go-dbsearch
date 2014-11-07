package xSql

import (
	"testing"
)

func TestWhere(t *testing.T) {

	sql := ""
	values := []interface{}{}

	/* check select */
	marks := []string{
		"LIKE", "ILIKE", "=", ">=", "<=", "<>", ">", "<",
	}

	for _, m := range marks {
		sql, values = Mark("f", m, "1").Comp()
		check_result(t, sql, " f "+m+" $1 ", values, 1)
	}

	/* check IN */
	sql, values = Mark("f", "IN", 1, 2, "3").Comp()
	check_result(t, sql, " f IN ( $1, $2, $3 ) ", values, 3)

	/* check IS */
	sql, values = Mark("f", "IS", "NULL").Comp()
	check_result(t, sql, " f IS NULL", values, 0)

	sql, values = Mark("f", "IS", "NOT NULL").Comp()
	check_result(t, sql, " f IS NOT NULL", values, 0)

	/* check logic condition */
	sql, values = NLogic("OR").Append(Mark("f", "=", "1")).Append(Mark("t", "<=", "2")).Comp()
	check_result(t, sql, "( f = $1  OR  t <= $2 )", values, 2)

	sql, values = NLogic("AND").Append(Mark("f", "<>", "1")).Append(Mark("t", ">=", "2")).Comp()
	check_result(t, sql, "( f <> $1 AND t >= $2 )", values, 2)

	/* check condition with function */
	sql, values = NLogic("AND").Append(Func("startdate < now()")).Append(Func("enddate > now()")).Comp()
	check_result(t, sql, "( startdate < now() AND enddate > now() )", values, 0)

	/* check condition with function for json */
	js := map[string]interface{}{
		"a": 1,
		"b": "name",
	}
	sql, values = NLogic("AND").Append(Mark("f::json->10", "@>", 12)).
		Append(Func("f::json->10 LIKE '%cat%'")).
		Append(Mark("f::json", "<>", js).Json()).Comp()
	check_result(t, sql, "(  f::json->10 @> $1  AND f::json->10 LIKE '%cat%' AND  f::json <> $2 )", values, 2)

	/* check array */
	array_marks := []string{
		"=", "<>", "<", ">", "<=", ">=", "@>", "<@", "&&", "||",
	}
	for _, m := range array_marks {
		sql, values = Array("f", m, 1, 2, 3).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3)
	}

	for _, m := range array_marks {
		sql, values = TArray("int", "f", m, 1, 2).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2 ]::int[]", values, 2)
	}

	/* test insert */
	insert := Insert("public.mytable").Append(Mark("a", "=", 1)).Append(Array("b", "=", 1, 2, 3)).Append(Array("c", "="))
	insert.Append(TArray("text", "i", "=", "w"))
	insert.Append(Mark("*", "RET", "")).Append(Mark("e", "=", js).Json()).Append(Mark("b as d", "RET", ""))
	sql, values = insert.Comp()
	check_result(t, sql, "INSERT INTO public.mytable  ( a, b, c, i, e) VALUES "+
		" ($1 , ARRAY[ $2, $3, $4 ],  '{}' , ARRAY[ $5 ]::text[], $6 ) RETURNING *, b as d", values, 6)

	update := Update("public.mytable").Append(Mark("a", "=", 1)).Append(Array("b", "=", 1, 2, 3)).Append(Array("c", "="))
	update.Append(Mark("*", "RET", "")).Append(Mark("b as d", "RET", ""))
	update.Append(Mark("e", "=", js).Json())

	update_where := NLogic("AND").Append(Func("startdate < now()")).Append(Func("enddate > now()"))
	update_where.Append(Mark("e", "<>", js).Json())
	update.Where(update_where)
	sql, values = update.Comp()
	check_sql := "UPDATE public.mytable SET a = $2, b = ARRAY[$3, $4, $5], c = '{}', e = $6" +
		" WHERE (startdate < now() AND enddate > now() AND e <> $1) RETURNING *, b as d"
	check_result(t, sql, check_sql, values, 6)

	//t.Fatal("error insert xSql: text view")
}
