package xSql

import (
	"log"
	"regexp"
	"strings"
	"testing"
)

func check_result(t *testing.T, sql1 string, sql2 string, values []interface{}, count int) {
	log.Printf("SQL1: %s\n", sql1)
	log.Printf("SQL2: %s\n", sql2)
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
