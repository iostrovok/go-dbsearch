# Simple GO interface to postgresql database and SQL templater

	NOTICE! It is developer version.

### Introduction

Please see more inforamtion in OSM wiki: http://wiki.openstreetmap.org/wiki/API_v0.6

### Installing 

	go get github.com/iostrovok/go-dbsearch/dbsearch

### How use example (not add).

See example for more inforamtion. They are in "./dbsearch/example/" path.

	> go get github.com/iostrovok/go-dbsearch/dbsearch

# Using.

## xSql

### Simple AND / OR

#### Example AND
	sql, values := Select("public.mytable", "*").
		Logic("AND").
		Mark("a", "=", 1).
		Mark("b", "=", 2).
		Mark("c", "=", 3).
		Mark("d", "=", "cat").
		Comp()

#### Result
	sql is 
	"SELECT * FROM public.mytable WHERE (a = $1 AND b = $2 AND c = $3 AND d = $4)"
	
	values is 
	[ 1, 2, 3, "cat" ]

#### Example OR
	sql, values := Select("public.mytable", "*").
		Logic("OR").
		Mark("a", "=", 1).
		Mark("b", "=", 2).
		Mark("c", "=", 3).
		Mark("d", "=", "cat").
		Comp()

#### Result
	sql is 
	"SELECT * FROM public.mytable WHERE (a = $1 OR b = $2 OR c = $3 OR d = $4)"
	
	values is 
	[ 1, 2, 3, "cat" ]

### Combination "AND" and "OR"

#### Example

	And := xSql.Logic("AND").Func("group ILIKE '%beatles%'")
	Or1 := xSql.Logic("OR").Mark("f_name", "=", "Paul").Mark("f_name", "=", "John")
	Or2 := xSql.Logic("OR").Mark("l_name", "=", "McCartney").Mark("l_name", "=", "Lennon")

	And.Append(Or1).Append(Or2)

	sql_where, values_1 :=  And.Comp()

	sql_full, values_2 := xSql.Select("public.mytable", "DOB").Append(And).Comp()

#### Result
	sql_where is 
	"(group ILIKE '%beatles%' AND (f_name = $1 OR f_name = $2) AND (l_name = $3 OR l_name = $4))"

	sql_full is 
	"SELECT DOB FROM public.mytable WHERE (group ILIKE '%beatles%' AND (f_name = $1 OR f_name = $2) AND (l_name = $3 OR l_name = $4))"

	values_1, values_2 are
	[ "Paul", "John", "McCartney", "Lennon" ]
