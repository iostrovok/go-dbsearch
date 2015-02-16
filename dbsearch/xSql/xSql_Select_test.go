package xSql

import (
	"testing"
)

func Test(t *testing.T) {
	_001_Test_Select(t)
	_002_Test_Select(t)
	_011_Test_Select(t)
	_012_Test_Select(t)
	_020_Test_Select(t)
	_030_Test_Select(t)

	//t.Fatal("test case")
}

func _001_Test_Select(t *testing.T) {
	/*  Simple condition 1 */
	sql, values := Select("public.mytable", "*").
		Mark("parents", "=", "papa").
		Comp()
	checkResult(t, sql, "SELECT * FROM public.mytable WHERE parents = $1", values, 1)
}

func _002_Test_Select(t *testing.T) {
	/*  Simple condition 2 */
	sql, values := Select("public.mytable", "*").IN("parent", []interface{}{"mama", "papa"}).Comp()
	checkResult(t, sql, "SELECT * FROM public.mytable WHERE parent IN ( $1, $2 ) ", values, 2)
}

func _011_Test_Select(t *testing.T) {
	/* Simple AND */
	sql, values := Select("public.mytable", "*").
		Logic("AND").
		Mark("a", "=", 1).
		Mark("b", "=", 2).
		Mark("c", "=", 3).
		Mark("d", "=", "cat").
		Comp()
	checkResult(t, sql, "SELECT * FROM public.mytable WHERE (a = $1 AND b = $2 AND c = $3 AND d = $4)", values, 4)
}

func _012_Test_Select(t *testing.T) {
	/* Simple OR */
	sql, values := Select("public.mytable", "*").
		Logic("OR").
		Mark("a", "=", 1).
		Mark("b", "=", 2).
		Mark("c", "=", 3).
		Mark("d", "=", "cat").
		Comp()
	checkResult(t, sql, "SELECT * FROM public.mytable WHERE (a = $1 OR b = $2 OR c = $3 OR d = $4)", values, 4)
}

func _020_Test_Select(t *testing.T) {
	/* Common test  */
	res := "SELECT * FROM public.mytable WHERE (id > $1 AND f::json->10 LIKE '%cat%' AND f::json <> $2)"

	sql, values := Select("public.mytable", "*").
		Logic("AND").
		Mark("id", ">", 12).
		Func("f::json->10 LIKE '%cat%'").
		Json("f::json", "<>", map[string]interface{}{"sss": 1}).
		Comp()
	checkResult(t, sql, res, values, 2)
}

func _030_Test_Select(t *testing.T) {
	/* Combination "AND" and "OR" */
	And := Logic("AND").Func("group ILIKE '%beatles%'")
	Or1 := Logic("OR").Mark("f_name", "=", "Paul").Mark("f_name", "=", "John")
	Or2 := Logic("OR").Mark("l_name", "=", "McCartney").Mark("l_name", "=", "Lennon")

	And.Append(Or1).Append(Or2)

	sql_where, values_1 := And.Comp()

	sql_full, values_2 := Select("public.mytable", "DOB").Append(And).Comp()

	checkResult(t, sql_where, "(group ILIKE '%beatles%' AND (f_name = $1 OR f_name = $2) AND (l_name = $3 OR l_name = $4))", values_1, 4)
	checkResult(t, sql_full, "SELECT DOB FROM public.mytable WHERE "+sql_where, values_2, 4)
}
