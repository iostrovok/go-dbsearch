package dbsearch

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/iostrovok/go-dbsearch/dbsearch/sqler"
	_ "github.com/lib/pq"
	"log"
	"reflect"
	"runtime"
	"sync"
	"time"
)

var m sync.Mutex

type Searcher struct {
	db            *sql.DB
	poolSize      int
	dsn           string
	log           bool
	logFull       bool
	DieOnColsName bool
	LastCols      []string
}

func (s *Searcher) Close() error {
	return s.db.Close()
}

func (s *Searcher) DBH() *sql.DB {
	return s.db
}

func (s *Searcher) SetDieOnColsName(is_debug ...bool) {
	if len(is_debug) > 0 {
		s.DieOnColsName = is_debug[0]
	} else {
		s.DieOnColsName = true
	}
}

func (s *Searcher) SetDebug(is_debug ...bool) {
	if len(is_debug) > 0 {
		s.log = is_debug[0]
	} else {
		s.log = true
	}
}

func (s *Searcher) SetStrongDebug(is_debug ...bool) {
	if len(is_debug) > 0 {
		s.logFull = is_debug[0]
	} else {
		s.logFull = true
	}
}

func (s *Searcher) Ping() error {
	return s.db.Ping()
}

func SetDBI(db *sql.DB) (*Searcher, error) {

	s := new(Searcher)
	s.db = db

	return s, nil
}

func (s *Searcher) StartReConnect(rto_in ...int) {

	rto := 5
	if len(rto_in) > 0 && rto_in[0] > 0 {
		rto = rto_in[0]
	}

	go func() {
		// Pings our DB each rto_in seconds
		for {
			time.Sleep(time.Duration(rto) * time.Second)
			if err := s.Ping(); err != nil {
				s.Close()

				if db, err := sql.Open("postgres", s.dsn); err != nil {
					log.Println(err)
				} else {
					s.db = db
					s.db.SetMaxOpenConns(s.poolSize)
				}
			}
		}
	}()
}

func DBI(poolSize int, dsn string, stop_error ...bool) (*Searcher, error) {

	s := new(Searcher)

	db, _ := sql.Open("postgres", dsn)

	if err := db.Ping(); err != nil {
		if len(stop_error) > 0 && stop_error[0] {
			log.Fatalf("DB Error: %s\n", err)
		}
		return nil, errors.New(fmt.Sprintf("DB Error: %s\n", err))
	}

	s.db = db
	s.db.SetMaxOpenConns(poolSize)

	s.poolSize = poolSize
	s.dsn = dsn

	return s, nil
}

func (s *Searcher) GetCount(sqlLine string, values []interface{}) (int, error) {
	var count int
	err := s.db.QueryRow(sqlLine, values...).Scan(&count)
	return count, err
}

func (s *Searcher) GetOne(mType *AllRows, p interface{}, sqlLine string, values ...[]interface{}) error {

	R, err := s._initGet(mType, p, sqlLine, values...)
	if err != nil {
		return err
	}
	defer R.Rows.Close()
	for R.Rows.Next() {
		resultStr := mType.GetRowResult(R)
		reflect.Indirect(reflect.ValueOf(p)).Set(reflect.Indirect(reflect.ValueOf(resultStr)))
		break
	}

	mCheckError(R.Rows.Err())
	return nil
}

func (s *Searcher) Get(mType *AllRows, p interface{}, sqlLine string, values ...[]interface{}) error {

	R, err := s._initGet(mType, p, sqlLine, values...)
	if err != nil {
		return err
	}

	if R.UseFork && runtime.GOMAXPROCS(0) > 1 {
		return s._GetFork(mType, p, R)
	}
	return s._GetNoFork(mType, p, R)
}

func (s *Searcher) GetNoFork(mType *AllRows, p interface{}, sqlLine string, values ...[]interface{}) error {

	R, err := s._initGet(mType, p, sqlLine, values...)
	if err != nil {
		return err
	}

	return s._GetNoFork(mType, p, R)
}

func (s *Searcher) _GetNoFork(mType *AllRows, p interface{}, R *GetRowResultStr) error {
	defer R.Rows.Close()

	var sliceValue = reflect.Indirect(reflect.ValueOf(p))
	for R.Rows.Next() {
		resultStr := mType.GetRowResult(R)
		sliceValue.Set(reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(resultStr))))
	}

	mCheckError(R.Rows.Err())
	return nil
}

/* Force  */
func (s *Searcher) GetFork(mType *AllRows, p interface{}, sqlLine string, values ...[]interface{}) error {

	R, err := s._initGet(mType, p, sqlLine, values...)
	if err != nil {
		return err
	}
	return s._GetFork(mType, p, R)
}

func (s *Searcher) _GetFork(mType *AllRows, p interface{}, R *GetRowResultStr) error {
	defer R.Rows.Close()

	CountFork := 4
	var wg sync.WaitGroup

	resCh := make(chan *EnvelopeRowResult, 2*CountFork)
	sendCh := make([]chan *EnvelopeRowResult, CountFork)
	for i := 0; i < CountFork; i++ {
		sendCh[i] = make(chan *EnvelopeRowResult, 1)
		go GetRowResultRoutine(i, sendCh[i], resCh)
	}

	checkValue := &EnvelopeFull{
		Res: map[int]interface{}{},
	}

	wg.Add(1)

	CheckCountFork := CountFork
	go func() {
		defer wg.Done()
		for {
			select {
			case res := <-resCh:
				if res.IsLast {
					CheckCountFork--
					if CheckCountFork == 0 {
						return
					}
				} else {
					checkValue.Res[res.N] = res.Res
				}
			}
		}
	}()

	N := -1
	cycle := true
	for cycle {

		for i := 0; i < CountFork; i++ {

			if !R.Rows.Next() {
				cycle = false
				break
			}

			N++

			Dest, RawResult := R.PrepareDestFun()
			mCheckError(R.Rows.Scan(Dest...))

			E := EnvelopeRowResult{
				aRows:     mType,
				N:         N,
				R:         R,
				RawResult: RawResult,
				IsLast:    false,
			}

			select {
			case sendCh[0] <- &E:
			case sendCh[1] <- &E:
			case sendCh[2] <- &E:
			case sendCh[3] <- &E:
			}
		}
	}

	var sliceValue = reflect.Indirect(reflect.ValueOf(p))
	for i := 0; i < CountFork; i++ {
		close(sendCh[i])
	}
	wg.Wait()
	close(resCh)

	for i := 0; i < len(checkValue.Res); i++ {
		s := checkValue.Res[i]
		sliceValue.Set(reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(s))))
	}

	mCheckError(R.Rows.Err())
	return nil
}

func (s *Searcher) GetFace(mType *AllRows, sqlLine string,
	values ...[]interface{}) ([]map[string]interface{}, error) {

	out := []map[string]interface{}{}

	R, err := s._initGet(mType, nil, sqlLine, values...)
	if err != nil {
		return out, err
	}
	defer R.Rows.Close()

	for R.Rows.Next() {
		resultStr, err := mType.GetRowResultFace(R)
		if err != nil {
			return nil, err
		}
		out = append(out, resultStr)
	}

	mCheckError(R.Rows.Err())
	return out, nil
}

func (s *Searcher) GetFaceOne(mType *AllRows, sqlLine string,
	values ...[]interface{}) (map[string]interface{}, error) {

	out := map[string]interface{}{}

	R, err := s._initGet(mType, nil, sqlLine, values...)
	if err != nil {
		return out, err
	}
	defer R.Rows.Close()

	for R.Rows.Next() {
		out, err = mType.GetRowResultFace(R)
		if err != nil {
			return nil, err
		}
		break
	}

	mCheckError(R.Rows.Err())
	return out, nil
}

func mCheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (self *Searcher) GetRowsCount(table string) (int, error) {
	return self.GetCount(fmt.Sprintf("SELECT count(*) FROM %s", table), make([]interface{}, 0))
}

func (s *Searcher) Do(sql string, values ...interface{}) {
	_, err := s.db.Exec(sql, values...)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Searcher) Insert(table string, data map[string]interface{}) {
	sql, values := sqler.InsertLine(table, data)
	s.DoCommit(sql, values)
}

func (s *Searcher) Delete(table string, data_where map[string]interface{}) {
	sql, values := sqler.DeleteLine(table, data_where)
	s.DoCommit(sql, values)
}

func (s *Searcher) Update(table string, data_where, data_update map[string]interface{}) {
	sql, values := sqler.UpdateLine(table, data_update, data_where)
	s.DoCommit(sql, values)
}

func (s *Searcher) DoCommit(sql string, values_in ...[]interface{}) {

	values := []interface{}{}
	if len(values_in) > 0 {
		values = append(values, values_in[0])
	}

	if s.log {
		log.Printf("DoCommit: %s\n", sql)
	}

	txn, err := s.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = txn.Exec(sql, values...)
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
