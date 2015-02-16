package xSql

import (
	"testing"
)

func TestWhere(t *testing.T) {
	_01TestWhere(t)
	_02TestWhere(t)
	_03TestWhere(t)
	_04TestWhere(t)
	_05TestWhere(t)
	_30TestWhere(t)
	_31TestWhere(t)
	_32TestWhere(t)
	//t.Fatal("test case")
}

func _01TestWhere(t *testing.T) {

	sql := ""
	values := []interface{}{}

	/* check select */
	marks := []string{
		"LIKE", "ILIKE", "=", ">=", "<=", "<>", ">", "<",
	}

	for _, m := range marks {
		sql, values = Mark("f", m, "1").Comp()
		checkResult(t, sql, " f "+m+" $1 ", values, 1)
	}

	/* check IN */
	sql, values = Mark("f", "IN", 1, 2, "3").Comp()
	checkResult(t, sql, " f IN ( $1, $2, $3 ) ", values, 3)

	/* check IS */
	sql, values = Mark("f", "IS", "NULL").Comp()
	checkResult(t, sql, " f IS NULL", values, 0)

	sql, values = Mark("f", "IS", "NOT NULL").Comp()
	checkResult(t, sql, " f IS NOT NULL", values, 0)

	/* check logic condition */
	sql, values = NLogic("OR").Append(Mark("f", "=", "1")).Append(Mark("t", "<=", "2")).Comp()
	checkResult(t, sql, "( f = $1  OR  t <= $2 )", values, 2)

	sql, values = NLogic("AND").Append(Mark("f", "<>", "1")).Append(Mark("t", ">=", "2")).Comp()
	checkResult(t, sql, "( f <> $1 AND t >= $2 )", values, 2)

	/* check condition with function */
	sql, values = NLogic("AND").Append(Func("startdate < now()")).Append(Func("enddate > now()")).Comp()
	checkResult(t, sql, "( startdate < now() AND enddate > now() )", values, 0)

	/* check condition with function for json */
	js := map[string]interface{}{
		"a": 1,
		"b": "name",
	}
	sql, values = NLogic("AND").Append(Mark("f::json->10", "@>", 12)).
		Append(Func("f::json->10 LIKE '%cat%'")).
		Append(Mark("f::json", "<>", js).Json()).Comp()
	checkResult(t, sql, "(  f::json->10 @> $1  AND f::json->10 LIKE '%cat%' AND  f::json <> $2 )", values, 2)

	//t.Fatal("error insert xSql: text view")
}

func _02TestWhere(t *testing.T) {

	/* check IN */
	list := []int{1, 2, 3}
	sql, values := Mark("t", "IN", &list).Comp()
	checkResult(t, sql, " t IN ( $1, $2, $3 ) ", values, 3)
	//t.Fatal("error insert xSql: text view")
}

func _03TestWhere(t *testing.T) {
	/* check IN */
	list := []string{"adsad", "asdasdas", "asdasdasd"}
	sql, values := Mark("t", "IN", list).Comp()
	checkResult(t, sql, " t IN ( $1, $2, $3 ) ", values, 3)
	//t.Fatal("error insert xSql: text view")
}

func _04TestWhere(t *testing.T) {
	/* check IN */
	list := &[]interface{}{"adsad", 2, "asdasdasd"}
	sql, values := Mark("t", "IN", list).Comp()
	checkResult(t, sql, " t IN ( $1, $2, $3 ) ", values, 3)
	//t.Fatal("error insert xSql: text view")
}

func _05TestWhere(t *testing.T) {
	/* check IN */
	list := []interface{}{"adsad", 2, "asdasdasd"}
	sql, values := Mark("t", "IN", list).Comp()
	checkResult(t, sql, " t IN ( $1, $2, $3 ) ", values, 3)
	//t.Fatal("error insert xSql: text view")
}

func _30TestWhere(t *testing.T) {

	/* check array */
	array_marks := []string{
		"=", "<>", "<", ">", "<=", ">=", "@>", "<@", "&&", "||",
	}
	for _, m := range array_marks {
		sql, values := Array("f", m, 1, 2, 3).Comp()
		checkResult(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3, "_30TestWhere. 1.")
	}

	for _, m := range array_marks {
		sql, values := TArray("int", "f", m, 1, 2, 3).Comp()
		checkResult(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]::int[]", values, 3, "_30TestWhere. 2.")
	}

}

func _31TestWhere(t *testing.T) {

	list := []interface{}{"1", 2, "3"}
	sql := ""
	values := []interface{}{}

	/* check array */
	array_marks := []string{
		"=", "<>", "<", ">", "<=", ">=", "@>", "<@", "&&", "||",
	}
	for _, m := range array_marks {
		sql, values = Array("f", m, list).Comp()
		checkResult(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3, "_31TestWhere. 1. Array, Simple slice.")
	}

	for _, m := range array_marks {
		sql, values = TArray("int", "f", m, list).Comp()
		checkResult(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]::int[]", values, 3, "_31TestWhere. 2. Array, Simple slice.")
	}

	for _, m := range array_marks {
		sql, values = Array("f", m, &list).Comp()
		checkResult(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3, "_31TestWhere. 3. Array, Ref slice.")
	}

	for _, m := range array_marks {
		sql, values = TArray("int", "f", m, &list).Comp()
		checkResult(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]::int[]", values, 3, "_31TestWhere. 4. Array, Ref slice.")
	}

	//t.Fatal("error insert xSql: text view")
}

func _32TestWhere(t *testing.T) {

	list := []int{1, 2, 3}
	sql := ""
	values := []interface{}{}

	/* check array */
	array_marks := []string{
		"=", "<>", "<", ">", "<=", ">=", "@>", "<@", "&&", "||",
	}
	for _, m := range array_marks {
		sql, values = Array("f", m, list).Comp()
		checkResult(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3, "_32TestWhere. 1.")
	}

	for _, m := range array_marks {
		sql, values = TArray("int", "f", m, list).Comp()
		checkResult(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]::int[]", values, 3, "_32TestWhere. 2.")
	}

	for _, m := range array_marks {
		sql, values = Array("f", m, &list).Comp()
		checkResult(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3, "_32TestWhere. 3.")
	}

	for _, m := range array_marks {
		sql, values = TArray("int", "f", m, list).Comp()
		checkResult(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]::int[]", values, 3, "_32TestWhere. 4.")
	}

	//t.Fatal("error insert xSql: text view")
}
