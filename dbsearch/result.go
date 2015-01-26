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
	CountC    int
	resCountC int
}

type EnvelopeRowResult struct {
	N         int
	R         *GetRowResultStr
	RawResult []interface{}
	aRows     *AllRows
	IsLast    bool
	Point     int
	Res       interface{}
}

func NewGetRowResult(i int) *GetRowResultStr {
	out := GetRowResultStr{
		resCountC: 0,
	}
	return &out
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
