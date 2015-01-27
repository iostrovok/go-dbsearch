package dbsearch

import (
	"database/sql"
)

// An action transitions stochastically to a resulting score.
type ElemConvertFunc func() interface{}

type GetRowResultStr struct {
	Rows      *sql.Rows
	Cols      []string
	DestFunL  []ElemConvertFunc
	Dest      []interface{}
	RawResult []interface{}
	SkipList  map[int]bool
	UseFork   bool
}

type EnvelopeRowResult struct {
	Err       error
	N         int
	R         *GetRowResultStr
	RawResult []interface{}
	aRows     *AllRows
	IsLast    bool
	Point     int
	Res       interface{}
	ResM      map[string]interface{}
}

type EnvelopeFull struct {
	Err  error
	Res  map[int]interface{}
	ResM map[int]map[string]interface{}
}

func (R *GetRowResultStr) PrepareDestFun() ([]interface{}, []interface{}) {
	Dest := make([]interface{}, len(R.DestFunL))
	RawResult := make([]interface{}, len(R.DestFunL))

	for i, v := range R.DestFunL {
		RawResult[i] = v()
		Dest[i] = &RawResult[i]
	}

	return Dest, RawResult
}
