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

