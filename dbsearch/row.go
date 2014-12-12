package dbsearch

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/iostrovok/go-iutils/iutils"
	"log"
	"time"
)

type Row struct {
	dbh     *Searcher
	date    map[string]interface{}
	cols    []string
	columns map[string]bool
}

func (s *Searcher) One(mType *AllRows, sqlLine string, values ...[]interface{}) (*Row, error) {
	value := []interface{}{}
	if len(values) > 0 {
		value = values[0]
	}
	row, err := s.GetOne(mType, sqlLine, value)
	if err != nil {
		return nil, err
	}
	return s.SetRow(row, s.LastCols), nil
}

func (s *Searcher) SetRow(date map[string]interface{}, cols []string, columns ...map[string]bool) *Row {

	r := Row{s, date, cols, map[string]bool{}}

	if len(columns) > 0 {
		r.columns = columns[0]
	} else {
		for _, n := range cols {
			r.columns[n] = true
		}
	}
	return &r
}

func (r *Row) _checkColumn(name string) bool {
	if _, find := r.columns[name]; find {
		return true
	}

	if r.dbh.DieOnColsName {
		log.Fatalln(" Not found" + name)
	}

	return false
}

func (r *Row) IsEmpty() bool {
	if len(r.date) > 0 {
		return false
	}
	return true
}

func (r *Row) Cols() []string {
	return r.cols
}

func (r *Row) Interface() map[string]interface{} {
	return r.date
}

func (r *Row) Str(name string) string {
	if !r._checkColumn(name) {
		return ""
	}

	v, find := r.date[name]
	if !find {
		return ""
	}
	return iutils.AnyToString(v)
}

func (r *Row) Float(name string) float64 {
	if !r._checkColumn(name) {
		return 0.0
	}

	spew.Dump(r.date)
	return iutils.AnyToFloat64(r.date[name])
}

func (r *Row) Int(name string) int {
	if !r._checkColumn(name) {
		return 0
	}

	return iutils.AnyToInt(r.date[name])
}

func (r *Row) IntArray(name string) []int {
	if !r._checkColumn(name) {
		return []int{}
	}

	return iutils.AnyToIntArray(r.date[name])
}

func (r *Row) StrArray(name string) []string {
	if !r._checkColumn(name) {
		return []string{}
	}

	return iutils.AnyToStringArray(r.date[name])
}

func (r *Row) Date(name string) *time.Time {
	out := r.DateTime(name)

	if out == nil {
		return nil
	}

	ret := out.Truncate(24 * time.Hour)
	return &ret
}

func (r *Row) Time(name string) *time.Time {
	out := r.DateTime(name)

	if out == nil {
		return nil
	}

	return out
}

func (r *Row) DateTime(name string) *time.Time {
	if !r._checkColumn(name) {
		log.Fatalln(" Not found" + name)
	}

	switch r.date[name].(type) {
	case time.Time:
		out := r.date[name].(time.Time)
		return &out
	}

	return nil
}
