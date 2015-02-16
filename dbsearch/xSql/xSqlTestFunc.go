package xSql

import (
	"log"
	"regexp"
	"strings"
	"testing"
)

func checkResult(t *testing.T, sql1 string, sql2 string, values []interface{}, count int, mess ...string) {

	if len(mess) > 0 {
		log.Printf("%s\n", mess[0])
	}

	log.Printf("SQL1: %s\n", sql1)
	log.Printf("SQL2: %s\n", sql2)
	log.Printf("values: %v\n", values)

	var N = regexp.MustCompile(`\s+`)
	s1 := strings.TrimSpace(strings.ToLower(N.ReplaceAllString(sql1, "")))
	s2 := strings.TrimSpace(strings.ToLower(N.ReplaceAllString(sql2, "")))

	if s1 != s2 {
		t.Fatalf("error where xSql: sqlLine --%#v-- for sql=\"%s\"\n", values, sql2)
	}
	if len(values) != count {
		t.Fatalf("error where xSql: values --%#v-- for sql=\"%s\"\n", values, sql2)
	}
}
