package xSql

import (
	"github.com/iostrovok/go-iutils/iutils"
	"regexp"
)

var (
	regQuote = regexp.MustCompile(`'`)
)

func Quote(val interface{}) string {

	if val == nil {
		return "NULL"
	}

	str := iutils.AnyToString(val)
	str = regQuote.ReplaceAllString(str, "''")

	return "'" + str + "'"
}
