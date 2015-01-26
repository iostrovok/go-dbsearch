#GP := $(GOPATH) 
GP := $(shell dirname $(realpath $(lastword $(GOPATH))))
ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
XSQL := ${ROOT}/dbsearch/xSql
MPST := ${ROOT}/dbsearch/mapstructure
DBS := ${ROOT}/dbsearch
BIN  := ${ROOT}/bin
#export GOBIN := ${ROOT}/bin
export GOPATH := ${MPST}:${GOPATH}:${ROOT}:${XSQL}:${DBS}
export GOPATH := /kmsearch/go_path_tmp/:${ROOT}
export PG_USER := postgres
export PG_PASSWD := 
export PG_HOST := 127.0.0.1
export PG_PORT := 5432
export DBNAME := pqgotest
export SSLMODE := 

#.PHONY: all test build index import run

test:
	echo ${GOPATH}
	go test ./dbsearch/

test-xsql:
	go test ./dbsearch/xSql/*.go 

test-xsql-utils:
	go test ./dbsearch/xSql/xSql_Utils_test.go ./dbsearch/xSql/xSqlTestFunc.go  ./dbsearch/xSql/xSqlUtils.go

test-xsql-insert:
	go test ./dbsearch/xSql/xSql_Insert_test.go ./dbsearch/xSql/xSqlTestFunc.go ./dbsearch/xSql/xSql.go ./dbsearch/xSql/xSqlJson.go ./dbsearch/xSql/xSqlArray.go

test-xsql-delete:
	go test ./dbsearch/xSql/xSql_Delete_test.go ./dbsearch/xSql/xSqlTestFunc.go ./dbsearch/xSql/xSql.go ./dbsearch/xSql/xSqlJson.go ./dbsearch/xSql/xSqlArray.go

test-xsql-select:
	go test ./dbsearch/xSql/xSql_Select_test.go ./dbsearch/xSql/xSqlTestFunc.go ./dbsearch/xSql/xSql.go ./dbsearch/xSql/xSqlJson.go ./dbsearch/xSql/xSqlArray.go

test-xsql-update:
	go test ./dbsearch/xSql/xSql_Update_test.go ./dbsearch/xSql/xSqlTestFunc.go ./dbsearch/xSql/xSql.go ./dbsearch/xSql/xSqlJson.go ./dbsearch/xSql/xSqlArray.go

test-xsql-where:
	go test ./dbsearch/xSql/xSql_Where_test.go ./dbsearch/xSql/xSqlTestFunc.go ./dbsearch/xSql/xSql.go ./dbsearch/xSql/xSqlJson.go ./dbsearch/xSql/xSqlArray.go

test-xsql-cover:
	go test -cover -coverprofile ./tmp.out ./dbsearch/xSql/*.go 
	sed 's/command-line-arguments/\.\/dbsearch\/xSql/' < ./tmp.out > ./tmp_fix.out
	go tool cover -html=./tmp_fix.out -o xSql.html
	rm ./tmp_fix.out ./tmp.out	


test-speed:
	go test  -bench=".*" ./dbsearch/result.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/dbsearch_test.go ./dbsearch/speed_test.go 

test-o:
	go test ./dbsearch/result.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/dbsearch_test.go

test-f:
	go test ./dbsearch/result.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/main_test.go ./dbsearch/dbsearch_test.go

test-s:
	go test ./dbsearch/result.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/slice_test.go ./dbsearch/dbsearch_test.go

test-d:
	go test ./dbsearch/result.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/date_test.go ./dbsearch/dbsearch_test.go

test-a:
	go test ./dbsearch/result.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/array_test.go ./dbsearch/dbsearch_test.go

test-l:
	go test ./dbsearch/result.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/autoload_test.go ./dbsearch/dbsearch_test.go

test-c:
	go test ./dbsearch/result.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/exception_test.go ./dbsearch/dbsearch_test.go

test-e:
	go test ./dbsearch/result.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/empty_columns_test.go ./dbsearch/dbsearch_test.go

test-i:
	go test ./dbsearch/result.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/interface_test.go ./dbsearch/dbsearch_test.go

test-m: test-o test-l test-a test-s test-f test-d test-e test-i test-c

test-full:
	go test ./dbsearch/array_test.go ./dbsearch/autoload_test.go ./dbsearch/convert_func.go ./dbsearch/date_test.go ./dbsearch/dbsearch.go ./dbsearch/dbsearch_test.go ./dbsearch/empty_columns_test.go ./dbsearch/field.go ./dbsearch/interface_test.go ./dbsearch/main_test.go ./dbsearch/mapstructure.go ./dbsearch/slice_test.go

test-cover:
	echo "Start recover"
	go test -cover -coverprofile ./tmp.out ./dbsearch/autoload_test.go ./dbsearch/empty_columns_test.go ./dbsearch/interface_test.go ./dbsearch/array_test.go ./dbsearch/main_test.go ./dbsearch/date_test.go ./dbsearch/slice_test.go ./dbsearch/field.go ./dbsearch/convert_func.go ./dbsearch/mapstructure.go ./dbsearch/dbsearch.go ./dbsearch/dbsearch_test.go
	echo "Start rendering"
	sed 's/command-line-arguments/\.\/dbsearch/' < ./tmp.out > ./tmp_fix.out
	go tool cover -html=./tmp_fix.out -o dbSerch.html
	rm ./tmp_fix.out ./tmp.out
	echo "Finish recover"

clean:
	rm ./tmp_fix.out ./tmp.out ./xSql.html ./dbSerch.html

