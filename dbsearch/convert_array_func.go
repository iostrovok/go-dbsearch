package dbsearch

import (
	"fmt"
	"github.com/iostrovok/go-iutils/iutils"
	"regexp"
	"strconv"
	"strings"
)

// construct a regexp to extract values:
var (
	unquotedRe  = regexp.MustCompile(`([^",\\{}\s]|NULL)+,`)
	_arrayValue = fmt.Sprintf("\"(%s)+\",", `[^"\\]|\\"|\\\\`)
	quotedRe    = regexp.MustCompile(_arrayValue)

	intArrayBrace = regexp.MustCompile(`[^-0-9\.\,]+`)
	intArraySplit = regexp.MustCompile(`,`)
	intArrayTail  = regexp.MustCompile(`\.[0-9]*`)

	noNumberDots      = regexp.MustCompile(`[^-0-9\.,]+`)
	noNumberDotsSplit = regexp.MustCompile(`(,|\s+)+`)

	parseBoolArrayRe     = regexp.MustCompile(`[^FTft]+`)
	parseBoolArrayReTail = regexp.MustCompile(`^[^FTft]+|[^FTft]+$`)
)

func parseBoolArrayForBool(s interface{}) []bool {
	line := parseBoolArrayReTail.ReplaceAllString(_AnyToString(s), "")
	r := parseBoolArrayRe.Split(line, -1)
	out := make([]bool, len(r))
	for i, v := range r {
		if v == "t" || v == "T" {
			out[i] = true
		} else {
			out[i] = false
		}
	}
	return out
}

func parseBoolArrayForString(s interface{}) []bool {
	r := parseArray(_AnyToString(s))
	out := make([]bool, len(r))
	for i, v := range r {
		v := strings.TrimSpace(v)
		if v == "" {
			out[i] = false
		} else {
			out[i] = true
		}
	}
	return out
}

func parseBoolArrayForReal(s interface{}) []bool {
	r := parseFloat64Array(s)
	out := make([]bool, len(r))
	for i, v := range r {
		if v == 0.0 {
			out[i] = false
		} else {
			out[i] = true
		}
	}
	return out
}

func parseBoolArrayForNumber(s interface{}) []bool {
	r := parseInt64Array(s)
	out := make([]bool, len(r))
	for i, v := range r {
		if v == 0 {
			out[i] = false
		} else {
			out[i] = true
		}
	}
	return out
}

func parseUint64Array(s interface{}) []uint64 {
	r := parseInt64Array(s)
	out := make([]uint64, len(r))
	for i, v := range r {
		if v < 0 {
			out[i] = 0
		} else {
			out[i] = uint64(v)
		}
	}
	return out
}

func parseUint8Array(s interface{}) []uint8 {
	r := parseIntArray(s)
	out := make([]uint8, len(r))
	for i, v := range r {
		if v < 0 {
			out[i] = 0
		} else {
			out[i] = uint8(v)
		}
	}
	return out
}

func parseUintArray(s interface{}) []uint {
	r := parseIntArray(s)
	out := make([]uint, len(r))
	for i, v := range r {
		if v < 0 {
			out[i] = 0
		} else {
			out[i] = uint(v)
		}
	}
	return out
}

func parseInt64Array(s interface{}) []int64 {
	str := strings.TrimSpace(_AnyToString(s))
	str = intArrayBrace.ReplaceAllString(str, "")
	str = intArrayTail.ReplaceAllString(str, "")
	k := intArraySplit.Split(str, -1)

	out := make([]int64, len(k))

	for i, v := range k {
		v := intArrayTail.ReplaceAllString(v, "")

		if v == "" {
			out[i] = 0
			continue
		}

		j, err := strconv.Atoi(v)
		if err != nil {
			//log.Println(err)
			out[i] = 0
			continue
		}
		out[i] = int64(j)
	}
	return out
}

func parseIntArray(s interface{}) []int {
	k := parseInt64Array(s)
	out := make([]int, len(k))
	for i, v := range k {
		out[i] = int(v)
	}
	return out
}

func parseFloat64Array(s interface{}) []float64 {
	out := []float64{}

	str := strings.TrimSpace(iutils.AnyToString(s))
	str = noNumberDots.ReplaceAllString(str, "")
	list := noNumberDotsSplit.Split(str, -1)

	for _, v := range list {
		out = append(out, iutils.AnyToFloat64(v))
	}

	return out
}

func parseFloat32Array(s interface{}) []float32 {
	out := []float32{}

	str := strings.TrimSpace(iutils.AnyToString(s))
	str = noNumberDots.ReplaceAllString(str, "")
	list := noNumberDotsSplit.Split(str, -1)

	for _, v := range list {
		out = append(out, float32(iutils.AnyToFloat64(v)))
	}

	return out
}
func parseArray(line string) []string {

	out := []string{}
	if line == "{}" {
		return out
	}

	if len(line)-1 != strings.LastIndex(line, "}") || strings.Index(line, "{") != 0 {
		return out
	}

	/* Removes lead & last {} and adds "," to end of string */
	line = strings.TrimPrefix(line, "{")
	line = strings.TrimSuffix(line, "}") + ","

	for len(line) > 0 {
		s := ""
		if strings.Index(line, `""`) == 0 {
			/* Empty line */
			s = ""
			line = line[3:]
		} else if strings.Index(line, `"`) != 0 {
			s = unquotedRe.FindString(line)
			line = line[strings.Index(line, ",")+1:]
			s = strings.TrimSuffix(s, ",")

			/* counvert NULL to empty string6 however we need nil string */
			if s == "NULL" {
				s = ""
			}
		} else {
			s = quotedRe.FindString(line)
			line = strings.TrimPrefix(line, s)
			s = strings.TrimPrefix(s, "\"")
			s = strings.TrimSuffix(s, "\",")
			s = strings.Join(strings.Split(s, "\\\\"), "\\")
			s = strings.Join(strings.Split(s, "\\\""), "\"")
		}
		out = append(out, s)
	}

	return out
}
