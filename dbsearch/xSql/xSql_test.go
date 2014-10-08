package xSql

import (
	//"github.com/davecgh/go-spew/spew"
	"log"
	"regexp"
	"strings"
	"testing"
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
		t.Fatal("error where xSql: sqlLine")
	}
	if len(values) != 8 {
		t.Fatal("error where xSql: values")
	}

	var insert_line = "insertintomytable(first_f,enddate,second_f,nextdate)values($1,now(),$2,interval'1day'+now())returning*,nextdateased"
	in := Insert("mytable")
	in.Append(Mark("first_f", "=", 10)).Append(Mark("enddate", "SQL", "now()")).Append(Mark("second_f", "=", "SUPER"))
	in.Append(Mark("nextdate", "SQL", "interval '1 day' + now()"))
	in.Append(Mark("*", "RET", ""))
	in.Append(Mark("nextdate as ED", "RET", ""))

	sql, values = in.Comp()

	log.Println(sql)
	//spew.Dump(values)

	st = strings.ToLower(N.ReplaceAllString(sql, ""))
	log.Println(st)
	log.Println(insert_line)
	if insert_line != st {
		t.Fatal("error insert xSql: sqlLine")
	}
	if len(values) != 2 {
		t.Fatal("error insert xSql: values")
	}
}
