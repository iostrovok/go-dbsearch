## Simple GO interface to postgresql database and SQL templater ##

### Introduction ###

No yet

### Installing ###

	go get github.com/iostrovok/go-dbsearch/dbsearch

### How use example (not add) ###

See example for more inforamtion. They are in "./dbsearch/example/" path.

	> go get github.com/iostrovok/go-dbsearch/dbsearch

# Using #
## DB-Search ##
### Import ###

```go

import "github.com/iostrovok/go-dbsearch/dbsearch"

```

### Quick start ###

#### SQL ####
```sql

drop table if exists public.person;

create table public.person ( 
		id serial, 
		active boolean default false,
		created timestamp default now(),
		changed timestamp default now(),
		dob date,
		fname varchar(50), 
		lname varchar(50), 
		children_names varchar(100)[], 
		cv text,
		disks json,
		count_refs int[] default '{0}'
);

INSERT INTO public.person ( active, dob, fname, lname, children_names, cv, disks, count_refs  ) 
VALUES ( true, '1942-06-18', 'Paul', 'McCartney', 
'{"Stella McCartney","James McCartney","Mary McCartney","Heather McCartney","Beatrice McCartney"}', 
'Sir James Paul McCartney MBE is ... Wikipedia', 
'{"disks":[{"title":"McCartney","year":1970},{"title":"Ram","year":1971},{"title":"Wild Life","year":1971}]}'::json,
'{1,2,5,3}'),
( false, '1940-10-09', 'John', 'Lennon', 
'{"Sean Lennon", "Julian Lennon"}', 
'John Ono Lennon, MBE ... Wikipedia', 
'{"disks":[{"title":"Imagine","year":1971},{"title":"Some Time in New York City","year":1972},{"title":"Mind Games","year":1973}]}'::json,
'{1,2,5,3}');

```
#### GO ####
```go
package main

import "github.com/iostrovok/go-dbsearch/dbsearch"
import "fmt"
import "time"
import "log"

type Singer struct {
	Id     int
	Active bool
	//Created       time.Time // Skip. DBA's Inner fields %)
	//Changed       time.Time // Skip. DBA's Inner fields %)
	Dob           time.Time
	Fname         string
	Lname         string
	ChildrenNames []string
	Cv            string
	Disks         map[string]interface{}
	CountRefs     []int
}

var mSinger *dbsearch.AllRows = &dbsearch.AllRows{
	Table:  "person",
	Schema: "public",
}

func main() {
	dbh, err := dbsearch.DBI(10, "host=127.0.0.1 port=5432 user=postgres dbname=pqgotest sslmode=disable", false)
	if err != nil {
		log.Fatal(err)
	}

	sql := "SELECT * FROM public.person WHERE lname = $1 "
	values := []interface{}{"Lennon"}

	Lennon := Singer{}
	if err := dbh.GetOne(mSinger, &Lennon, sql, values); err != nil {
		log.Fatal(err)
	}

	Singers := []Singer{}
	if err := dbh.Get(mSinger, &Singers, sql, values); err != nil {
		log.Fatal(err)
	}

	/* View list of records */
	fmt.Printf("%s %s (%s)\n%s\n", Lennon.Lname, Lennon.Fname, Lennon.Dob.Format("Jan 2 2006"), Lennon.Cv)
	fmt.Printf("Children: %#v\n", Lennon.ChildrenNames)
	fmt.Printf("Discography: %#v\n", Lennon.Disks)
	fmt.Printf("Count referres: %#v\n", Lennon.CountRefs)

	/* View list of records */
	fmt.Printf("\n\nSingers:\n%#v\n", Singers)

}
```

## Dependency Types ##

## Simple Types ##

| From Postgresql | To Strustrue | Additional  |
| ------------- | :------------- | :----- |
| bigint, smallintint, integer, serial, bigserial, smallserial | int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string |  |
| bigint, smallintint, integer, serial, bigserial, smallserial | []int, []int64, []uint64, []float32, []float64, []string,  | Result is placed to first element of slice |
| real, double, numeric, decimal, money | int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string |  |
| real, double, numeric, decimal, money | []int, []int64, []float32, []float64, []string,  | Result is placed to first element of slice |
| bytea | []byte, []uint8, string |  |
| varchar, char, text | string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64,  | |
| varchar, char, text | []string, []int, []int64, []float32, []float64  |  Result is placed to first element of slice  |
| boolean, bool | bool, boolean,  | it's honest boolean value |
| boolean, bool | int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64,  | Result is 0 (false) or 1 (true)|
| boolean, bool | string, | Result is "" (false) or "1" (true) |
| boolean, bool  | []int, []int64, []uint, []uint8, []uint64, []float32, []float64, []string, []byte | Result is placed to first element of slice |

## Date & Time ##

| From Postgresql | To Strustrue | Additional  |
| ------------- | :------------- | :----- |
| date, time, timestamp | time.Time, | Package "time" is used |
| date, time, timestamp | string, | Result is string(t.Format("2006-01-02 15:04:05 -0700")) |
| date, time, timestamp | int, int64, | Result is time.Unix() from package "time" |
| date, time, timestamp | uint, uint64, | Result is time.Unix() from package "time". There is error for time.Unix() < 0 |
| date, time, timestamp | float64, | Result is "float64(t.UnixNano()) / 1000" from package "time" |
| date, time, timestamp | map[string]int, map[string]int64, | Result is map[string]int{					"year": t.Year(), "month": int(t.Month()), "day": t.Day(), "hour": t.Hour(), "minute": t.Minute(), "second": t.Second(), "nanosecond": t.Nanosecond(), "zone": int(t.Zone()), } |
| date, time, timestamp | []int, []int64, | Result is []int{ t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), int(t.Zone()), } |

## JSON & JSONB ##

| From Postgresql | To Strustrue | Additional  |
| ------------- | :------------- | :----- |
| json, jsonb | map[string]interface{}, string |  |
| varchar, char, text | map[string]interface{} | Tag \`type:"json"\` is necessary |

## ARRAYs ##

| From Postgresql | To Strustrue | Additional  |
| ------------- | :------------- | :----- |
| []boolean, []bool | []int, []int64, []uint, []byte, []uint8, []uint64 | Result is 0 (false) or 1 (true)  |
| []boolean, []bool | []string | Result is "" (false) or "1" (true)  |
| []boolean, []bool | []bool ||
| []varchar, []char, []text | []string, []int, []int64, []uint, []byte, []uint8, []uint64  ||
| []bigint, []smallintint, []integer, []real, []double, []numeric, []decimal, []money | []string, []int, []int64, []uint, []byte, []uint8, []uint64  ||


### Connect/Init ###
```go

import "github.com/iostrovok/go-dbsearch/dbsearch"

pool_size := 10	
stop_error := true // die if connect has errors
dsn := "user=pqgotest dbname=pqgotest sslmode=verify-full"

dbh, err := dbsearch.DBI(pool_size, dsn, stop_error)
dbh.StartReConnect(10) // Tries to reconnect each 10 seconds if connect was broken

```

or

```go

import (
	"github.com/iostrovok/go-dbsearch/dbsearch"
	_ "github.com/lib/pq"
	"database/sql"
)

db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
if err != nil {
	log.Fatal(err)
}
dbh := dbsearch.SetDBI( db )
```

### Read map[string]interface{} / []map[string]interface{} ###
```go

import "github.com/iostrovok/go-dbsearch/dbsearch"

type Singer struct {
	Id     int
	Active bool
	Fname         string
	Lname         string
}

var mSinger *dbsearch.AllRows = &dbsearch.AllRows{
	Table:  "person",
	Schema: "public",
}

sql := "SELECT * FROM public.singer LIMIT 1"
map, err := dbh.GetFace(mSinger, sql)
// map is map[string]interface{} likes that: { db_filed: interface{} }

sql := "SELECT * FROM public.singer"
slice, err := dbh.GetFace(mSinger, sql)
// slice is []map[string]interface{}

```

### Read structure fields ###
#### The structure definition ways  ####

1) Tags "bd" & "type"

2) Title table and special title fields

### Tags "bd" & "type" ###

db - column name

type - postresql's column type

```go

import "github.com/iostrovok/go-dbsearch/dbsearch"
import "reflect"

type Singer struct {
	Id     				int 	`db:"id"     type:"serial"`
	Active 				bool  	`db:"active" type:"bool"`
	Name_of_Singer		string	`db:"fname"  type:"text"`
	SuperName_of_Singer	string	`db:"lname"  type:"text"`
}

var mSinger *dbsearch.AllRows = &dbsearch.AllRows{
	SType: reflect.TypeOf(Singer{}), // Optional
}

...

p := []Singer{}
dbh.Get(mSinger, &p, "SELECT * FROM public.singer ORDER BY 1")

```
### Title table and special title fields ###

If we defined table name in our *Allows structure and defined structure field by our "conversion rules" we can leave out "db" & "type" tags.

##### Rules of conversion from database column titles to structure names ####

1) "\<first letter>" turns into "uppercase(\<first letter>)"

2) "\<letter 1>_<letter 2>" turns into "\<letter 1>uppercase(\<letter 2>)"

#######Examples:#######

my_long_column_title -> MyLongColumnTitle

t -> T

name -> Name

```go

import "github.com/iostrovok/go-dbsearch/dbsearch"
import "reflect"

type Singer struct {
	Id     	int 	// db: "id"     -> struct field "Id"
	Active 	bool  	// db: "active" -> struct field "Active"
	Fname	string	// db: "fname"  -> struct field "Fname"
	Lname 	string	// db: "lname"  -> struct field "Lname"
}

var mSinger *dbsearch.AllRows = &dbsearch.AllRows{
	Table: "singer",  // necessary
	Schema: "public", // optional, default "public"
	SType: reflect.TypeOf(Singer{}), // Optional
}

...

p := []Singer{}
dbh.Get(mSinger, &p, "SELECT * FROM public.singer ORDER BY 1")

```
### Skip columns ###
In cases:

1) We have to skip columns, because we don't want to use extra fields in our structure

2) We get too more columns from select

we need to use 

```go
dbh.SetDieOnColsName(false)
```

### Read rows as []map[string]interface{} ###
If we want to get row(s) as map[string]interface{} (or []map[string]interface{} for list of columns) we need to use GetFaceOne and GetFace function. 
In this case we have to define:

1) tags "db" and "type" for each structure field

or

2) use "Rules of conversion" and define dbsearch.AllRows.Table and dbsearch.AllRows.SType properties


```go

import "github.com/iostrovok/go-dbsearch/dbsearch"
import "reflect"
import "log"

type Singer struct {
	Id     	int 	// db: "id"     -> struct field "Id"
	Active 	bool  	// db: "active" -> struct field "Active"
	Fname	string	// db: "fname"  -> struct field "Fname"
	Lname 	string	// db: "lname"  -> struct field "Lname"
}

var mSinger *dbsearch.AllRows = &dbsearch.AllRows{
	Table: "singer",  // necessary
	Schema: "public", // optional, default "public"
	SType: reflect.TypeOf(Singer{}), // Optional
}

//...

sql := "SELECT * FROM public.person WHERE fname IN( $1, $2) "
values := []interface{}{"John", "Paul"}
slice, err := dbh.GetFace(mSinger, sql, values)
if err != nil {
	log.Panicln(err)
}

// or 

sql := "SELECT * FROM public.person fname = $ LIMIT 1"
values := []interface{}{"John", "Paul"}
map, err := dbh.GetFace(mSinger, sql, values)
if err != nil {
	log.Panicln(err)
}

```

### JSON and ARRAY ###

We can get json as map[string]interface{} and arrays from db.

```go

type ArrayPlace struct {
	Col1  []in			`db:"col1" type:"[]bigint"`
	Col2  []in			`db:"col2" type:"[]smallint"`
	Col3  []in			`db:"col3" type:"[]int"`
	Col4  []strin			`db:"col4" type:"[]string"`
	
       // col6 has 'text', 'varchar', 'json' or 'jsonb' type in table
       Col5  map[string]interface{}	`db:"col5" type:"json"`
}

### type Searcher ###

```go

type Searcher struct {
}s

```

Searcher contains information about connection to db, log level e.t.c

#### func DBI(int, string, stop_error ...) ####
#####func DBI(poolSize int, dsn string, stop_error ...bool) (*Searcher, error)#####

Returns *Searcher object.

#### func SetDBI(*sql.DB) ####
#####func SetDBI(db *sql.DB) (*Searcher, error)#####

It defines db connect for current if *Searcher already object has existed

#### func (*Searcher) GetOne ####
#####func (s *Searcher) GetOne(mType *AllRows, p interface{}, sqlLine string, values ...[]interface{}) error#####

Returns one rows from DB with sqlLine and values.

#### func (*Searcher) Get ####
#####func (s *Searcher) Get(mType *AllRows, p interface{}, sqlLine string, values ...[]interface{}) error#####

Returns slice rows from DB with sqlLine and values.

```go

type A struct {
	Id   int
	Name string
}

var mType *dbsearch.AllRows = &dbsearch.AllRows{
	Table:  "test",
}
p := A{}
GetOne(mType, &p, "SELECT * FROM public.test")

s := []A{}
Get(mType, &p, "SELECT * FROM public.test")

```

#### func (*Searcher) GetFaceOne ####
#####func (s *Searcher) GetFaceOne(mType *AllRows, sqlLine string,
	values ...[]interface{}) (map[string]interface{}, error)#####

Returns one rows from DB with sqlLine and values as map[string]interface{}

#### func (*Searcher) GetFace ####
#####func (s *Searcher) GetFace(mType *AllRows, sqlLine string, values ...[]interface{}) ([]map[string]interface{}, error)#####

Returns slice rows from DB with sqlLine and values as []map[string]interface{}

```go

type A struct {
	Id   int
	Name string
}

var mType *dbsearch.AllRows = &dbsearch.AllRows{
	Table:  "test",
}
slice, err := GetFace(mType, "SELECT * FROM public.test")

map, err := GetFaceOne(mType, "SELECT * FROM public.test")

```

#### func (*Searcher) SetDieOnColsName ####
#####func (s *Searcher) SetDieOnColsName(isDie ...bool)#####

Sets to die or not to die when we have wrong column name or structure fields

### type AllRows ###

```go
type AllRows struct {
	SType         reflect.Type // type of returned structure
	Table         string       // Tables name
	Schema        string       // Schema name
	...
}
```
AllRows contains information which is needed for select a single row.


## xSql ##
### Import ###
```go
	import "github.com/iostrovok/go-dbsearch/dbsearch/xSql"
```

### Quote ###
```go
	testSuite := map[interface{}]string{
		"the molecule's structure": "'the molecule''s structure'",
		" I'''am an actor.":        "' I''''''am an actor.'",
		100500:                     "'100500'",
	}

	for k, v := range testSuite {
		q := xSql.Quote(k) // v == q
		log.Printf("%v : %s => %s\n", k, v, q)
	}
```
Notece!
Don't use xSql.Quote for float, real & structure.
So xSql.Quote(100.5) is "'100.500000'" on my system.

### Simple AND / OR ###
#### Example AND ####
```go
	sql, values := xSql.Select("public.mytable", "*").
		Logic("AND").
		Mark("a", "=", 1).
		Mark("b", "=", 2).
		Mark("c", "=", 3).
		Mark("d", "=", "cat").
		Comp()
```
#### Result
sql is 
```sql
	SELECT * FROM public.mytable WHERE (a = $1 AND b = $2 AND c = $3 AND d = $4)
```	
values is 

	[ 1, 2, 3, "cat" ]

#### Example OR
```go
	sql, values := xSql.Select("public.mytable", "*").
		Logic("OR").
		Mark("a", "=", 1).
		Mark("b", "=", 2).
		Mark("c", "=", 3).
		Mark("d", "=", "cat").
		Comp()
```
#### Result
sql is
```sql
	SELECT * FROM public.mytable WHERE (a = $1 OR b = $2 OR c = $3 OR d = $4)
```	
values is 
```go
	[]interface{}{1, 2, 3, "cat"}
```

### Combination "AND" and "OR"

#### Example
```go
	And := xSql.Logic("AND").Func("group ILIKE '%beatles%'")
	Or1 := xSql.Logic("OR").Mark("f_name", "=", "Paul").Mark("f_name", "=", "John")
	Or2 := xSql.Logic("OR").Mark("l_name", "=", "McCartney").Mark("l_name", "=", "Lennon")

	And.Append(Or1).Append(Or2)

	sql_where, values_1 :=  And.Comp()

	sql_full, values_2 := xSql.Select("public.mytable", "DOB").Append(And).Comp()
```
#### Result
sql_where is
```sql
	(group ILIKE '%beatles%' AND (f_name = $1 OR f_name = $2) AND (l_name = $3 OR l_name = $4))
```
sql_full is 
```sql
	SELECT DOB FROM public.mytable
	WHERE (group ILIKE '%beatles%' AND (f_name = $1 OR f_name = $2) AND (l_name = $3 OR l_name = $4))
```

values_1, values_2 are
```go
	[]interface{}{ "Paul", "John", "McCartney", "Lennon"}
```

### Update row(s)

#### Example
```go
	where := xSql.Logic("AND").Mark("f_name", "=", "John").Mark("l_name", "=", "Lennon").
		Mark("age", "<", 40)

	sql, values := xSql.Update("public.mytable").
		Mark("sended", "=", 1).
		Where(where).
		Comp()
```

#### Result
sql is
```sql
	UPDATE public.mytable SET sended = $4 WHERE (f_name = $1 AND l_name = $2 AND age < $3 )
```

values is 
```go
	[]interface{}{"John", "Lennon", 40, 1}
```

### Delete row(s)

#### Example
```go
	sql, values := xSql.Delete("public.mytable").
		Logic("AND").Mark("f_name", "=", "John").
		Mark("l_name", "=", "Lennon").
		Mark("age", "<", 40).
		Mark("age", ">", 0).
		Comp()
```

#### Result
sql is
```sql
	DELETE FROM public.mytable WHERE (f_name = $1 AND l_name = $2 AND age < $3 AND age > $4)
```

values is 
```go
	[]interface{}{"John", "Lennon", 40, 0}
```

### Pass slice
#### Example "Array" and "TArray"

```go
	// 1
	sql, values = Array("col", " = ", 1, 2, 3).Comp()
	
	// 2
	sql, values = TArray("int", "col", " = ", 1, 2, 3).Comp()
```

#### Result
sql is
```sql
	-- 1
	col = ARRAY[ $1, $2, $3 ]
	
	-- 2
	col = ARRAY[ $1, $2 ]::int[]
```
values is
```go
	// 1, 2 
	[]interface{}{1,2,3}
```

#### Example "IN"
```go
	// 1
	list := []interface{}{"adsad", 2, "asdasdasd"}
	sql, values := Mark("t", "IN", list).Comp()
	
	// 2
	list := []interface{}{"adsad", 2, "asdasdasd"}
	sql, values := Mark("t", "IN", &list).Comp()
	
	// 3
	sql, values := Mark("t", "IN", "adsad", 2, "asdasdasd").Comp()
	
	// 4
	list := []int{1, 2, 3}
	sql, values := Mark("t", "IN", list).Comp()
	
	// 5
	list := []int{1, 2, 3}
	sql, values := Mark("t", "IN", &list).Comp()
	
	// 6
	sql, values := Mark("t", "IN", 1, 2, 3).Comp()
```

#### Result
sql is
```sql
	-- 1,2,3,4,5,6 
	t IN ( $1, $2, $3 )
```
values is 
```go
	// 1,2,3
	[]interface{}{"adsad", 2, "asdasdasd"}
	// 4,5,6
	[]interface{}{1, 2, 3}
```

