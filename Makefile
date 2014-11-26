#GP := $(GOPATH) 
GP := $(shell dirname $(realpath $(lastword $(GOPATH))))
ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
XSQL := ${ROOT}/dbsearch/xSql
BIN  := ${ROOT}/bin
#export GOBIN := ${ROOT}/bin
export GOPATH := ${GOPATH}:${GP}:${GP}/go:${ROOT}:${XSQL}
export PG_USER := postgres
export PG_PASSWD := 
export PG_HOST := 127.0.0.1
export PG_PORT := 5432
export DBNAME := pqgotest
export SSLMODE := 

#.PHONY: all test build index import run

test:
	go test ./dbsearch/

test-xsql:
	go test -cover -coverprofile ./tmp.out ./dbsearch/xSql/*.go 
	sed 's/command-line-arguments/\.\/dbsearch\/xSql/' < ./tmp.out > ./tmp_fix.out
	go tool cover -html=./tmp_fix.out -o xSql.html
	rm ./tmp_fix.out ./tmp.out

clean:
	rm ./tmp_fix.out ./tmp.out ./xSql.html

