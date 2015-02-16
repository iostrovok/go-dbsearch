package xSql

import (
	"testing"
)

func TestWhere(t *testing.T) {
	_01_Test_Where(t)
	_02_Test_Where(t)
	_03_Test_Where(t)
	_04_Test_Where(t)
	_05_Test_Where(t)
	_30_Test_Where(t)
	_31_Test_Where(t)
	_32_Test_Where(t)
	//t.Fatal("test case")
}

func _01_Test_Where(t *testing.T) {

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

	//t.Fatal("error insert xSql: text view")
}

func _02_Test_Where(t *testing.T) {

	/* check IN */
	list := []int{1, 2, 3}
	sql, values := Mark("t", "IN", &list).Comp()
	check_result(t, sql, " t IN ( $1, $2, $3 ) ", values, 3)
	//t.Fatal("error insert xSql: text view")
}

func _03_Test_Where(t *testing.T) {
	/* check IN */
	list := []string{"adsad", "asdasdas", "asdasdasd"}
	sql, values := Mark("t", "IN", list).Comp()
	check_result(t, sql, " t IN ( $1, $2, $3 ) ", values, 3)
	//t.Fatal("error insert xSql: text view")
}

func _04_Test_Where(t *testing.T) {
	/* check IN */
	list := &[]interface{}{"adsad", 2, "asdasdasd"}
	sql, values := Mark("t", "IN", list).Comp()
	check_result(t, sql, " t IN ( $1, $2, $3 ) ", values, 3)
	//t.Fatal("error insert xSql: text view")
}

func _05_Test_Where(t *testing.T) {
	/* check IN */
	list := []interface{}{"adsad", 2, "asdasdasd"}
	sql, values := Mark("t", "IN", list).Comp()
	check_result(t, sql, " t IN ( $1, $2, $3 ) ", values, 3)
	//t.Fatal("error insert xSql: text view")
}

func _30_Test_Where(t *testing.T) {

	/* check array */
	array_marks := []string{
		"=", "<>", "<", ">", "<=", ">=", "@>", "<@", "&&", "||",
	}
	for _, m := range array_marks {
		sql, values := Array("f", m, 1, 2, 3).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3, "_30_Test_Where. 1.")
	}

	for _, m := range array_marks {
		sql, values := TArray("int", "f", m, 1, 2, 3).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]::int[]", values, 3, "_30_Test_Where. 2.")
	}

}

func _31_Test_Where(t *testing.T) {

	list := []interface{}{"1", 2, "3"}
	sql := ""
	values := []interface{}{}

	/* check array */
	array_marks := []string{
		"=", "<>", "<", ">", "<=", ">=", "@>", "<@", "&&", "||",
	}
	for _, m := range array_marks {
		sql, values = Array("f", m, list).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3, "_31_Test_Where. 1. Array, Simple slice.")
	}

	for _, m := range array_marks {
		sql, values = TArray("int", "f", m, list).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]::int[]", values, 3, "_31_Test_Where. 2. Array, Simple slice.")
	}

	for _, m := range array_marks {
		sql, values = Array("f", m, &list).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3, "_31_Test_Where. 3. Array, Ref slice.")
	}

	for _, m := range array_marks {
		sql, values = TArray("int", "f", m, &list).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]::int[]", values, 3, "_31_Test_Where. 4. Array, Ref slice.")
	}

	//t.Fatal("error insert xSql: text view")
}

func _32_Test_Where(t *testing.T) {

	list := []int{1, 2, 3}
	sql := ""
	values := []interface{}{}

	/* check array */
	array_marks := []string{
		"=", "<>", "<", ">", "<=", ">=", "@>", "<@", "&&", "||",
	}
	for _, m := range array_marks {
		sql, values = Array("f", m, list).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3, "_32_Test_Where. 1.")
	}

	for _, m := range array_marks {
		sql, values = TArray("int", "f", m, list).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]::int[]", values, 3, "_32_Test_Where. 2.")
	}

	for _, m := range array_marks {
		sql, values = Array("f", m, &list).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3, "_32_Test_Where. 3.")
	}

	for _, m := range array_marks {
		sql, values = TArray("int", "f", m, list).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]::int[]", values, 3, "_32_Test_Where. 4.")
	}

	//t.Fatal("error insert xSql: text view")
}
