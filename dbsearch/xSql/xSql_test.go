package xSql

import (
	//"github.com/davecgh/go-spew/spew"
	"log"
	"regexp"
	"strings"
	"testing"
)

func check_result(t *testing.T, sql1 string, sql2 string, values []interface{}, count int) {
	log.Printf("%s\n", sql1)
	log.Printf("%s\n", sql2)
	log.Printf("%v\n", values)

	var N = regexp.MustCompile(`\s+`)
	s1 := strings.TrimSpace(strings.ToLower(N.ReplaceAllString(sql1, "")))
	s2 := strings.TrimSpace(strings.ToLower(N.ReplaceAllString(sql2, "")))

	if s1 != s2 {
		t.Fatal("error where xSql: sqlLine for " + sql2)
	}
	if len(values) != count {
		t.Fatal("error where xSql: values for " + sql2)
	}
}

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

	/* check array */
	array_marks := []string{
		"=", "<>", "<", ">", "<=", ">=", "@>", "<@", "&&", "||",
	}
	for _, m := range array_marks {
		sql, values = Array("f", m, 1, 2, 3).Comp()
		check_result(t, sql, "f "+m+" ARRAY[ $1, $2, $3 ]", values, 3)
	}

	/* test insert */
	insert := Insert("public.mytable").Append(Mark("a", "=", 1)).Append(Array("b", "=", 1, 2, 3)).Append(Array("c", "="))
	insert.Append(Mark("*", "RET", "")).Append(Mark("b as d", "RET", ""))
	sql, values = insert.Comp()
	check_result(t, sql, "INSERT INTO public.mytable ( a, b, c ) VALUES ($1, ARRAY[$2, $3, $4], '{}') RETURNING *, b as d", values, 4)

	update := Update("public.mytable").Append(Mark("a", "=", 1)).Append(Array("b", "=", 1, 2, 3)).Append(Array("c", "="))
	update.Append(Mark("*", "RET", "")).Append(Mark("b as d", "RET", ""))
	update_where := NLogic("AND").Append(Func("startdate < now()")).Append(Func("enddate > now()"))
	update.Where(update_where)
	sql, values = update.Comp()
	check_sql := "UPDATE public.mytable SET a = $1, b = ARRAY[$2, $3, $4], c = '{}' " +
		"WHERE (startdate < now() AND enddate > now()) RETURNING *, b as d"
	check_result(t, sql, check_sql, values, 4)

	//t.Fatal("error insert xSql: text view")
}
