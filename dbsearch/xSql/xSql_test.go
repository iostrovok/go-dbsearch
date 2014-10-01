package xSql

import (
	"regexp"
	"strings"
	"testing"

	"log"
)

func TestWhere(t *testing.T) {

	var test_line = "((((fias_id=$1orfias_id>=$2orfias_id<$3orfias_id<=$4ornamein($5)orfias_idilike$6orfias_idlike$7)andid<$8)orfirst_canmeisnotnull)andsecond_anmeisnullandstartdate<now()andenddate>now())"
	c := Mark("fias_id", "=", "test_eq")
	c1 := Mark("fias_id", ">=", "test_eg")
	c2 := Mark("fias_id", "<", "test_lt")
	c3 := Mark("fias_id", "<=", "test_el")
	c4 := Mark("fias_id", "ILIKE", "test_el")
	c5 := Mark("fias_id", "LIKE", "test_el")
	t1 := NLogic("OR", c, c1, c2, c3).Append(Mark("name", "IN", "John"), c4, c5)
	t2 := NLogic("AND").Append(t1).Append(Mark("id", "<", 100))
	t3 := NLogic("OR").Append(t2).Append(Mark("first_canme", "IS", "NOT NULL"))
	top := NLogic("AND").Append(t3).Append(Mark("second_anme", "IS", "NULL"))
	sql, values := top.Append(Func("startdate < now()")).Append(Func("enddate > now()")).Comp()

	var N = regexp.MustCompile(`\s+`)
	st := strings.ToLower(N.ReplaceAllString(sql, ""))

	if test_line != st {
		t.Fatal("error xSql: sqlLine")
	}
	if len(values) != 8 {
		t.Fatal("error xSql: values")
	}
}
