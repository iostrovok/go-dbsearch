package dbsearch

import (
	"database/sql"
)

type GetRowResultStr struct {
	Rows       *sql.Rows
	Cols       []string
	DestL      [][]interface{}
	RawResultL [][]interface{}
	SkipList   map[int]bool
	CountC     int
	resCountC  int
}

func NewGetRowResult(i int) *GetRowResultStr {
	out := GetRowResultStr{
		resCountC: 0,
	}
	out.SetCountC(i)
	return &out
}

func (R *GetRowResultStr) ResetCountC() {
	R.resCountC = R.CountC - 1
}

func (R *GetRowResultStr) SetCountC(i int) {
	R.CountC = i
	R.DestL = make([][]interface{}, R.CountC)
	R.RawResultL = make([][]interface{}, R.CountC)
}

func (R *GetRowResultStr) AppendRawResult(s interface{}) bool {
	if R.resCountC > -1 {
		R.RawResultL[R.resCountC] = append(R.RawResultL[R.resCountC], s)
	}
	R.resCountC--
	return R.resCountC >= 0
}

func (R *GetRowResultStr) PassRawResult() {
	for i := range R.RawResultL {
		for j := range R.RawResultL[i] {
			R.DestL[i] = append(R.DestL[i], &R.RawResultL[i][j])
		}
	}
}
