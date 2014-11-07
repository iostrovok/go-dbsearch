#GP := $(GOPATH) 
GP := $(shell dirname $(realpath $(lastword $(GOPATH))))
ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
XSQL := ${ROOT}/dbsearch/xSql
BIN  := ${ROOT}/bin
#export GOBIN := ${ROOT}/bin
export GOPATH := ${GP}:${GP}/go:${ROOT}:${XSQL}

#.PHONY: all test build index import run

test-xsql:
	go test -cover -coverprofile ./tmp.out ./dbsearch/xSql/*.go 
	sed 's/command-line-arguments/\.\/dbsearch\/xSql/' < ./tmp.out > ./tmp_fix.out
	go tool cover -html=./tmp_fix.out -o xSql.html
	rm ./tmp_fix.out ./tmp.out

clean:
	rm ./tmp_fix.out ./tmp.out ./xSql.html

