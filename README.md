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

### Get single row ###
```go
	import "github.com/iostrovok/go-dbsearch/dbsearch"

	sql := "SELECT * FROM my_table WHERE id = $1 "
	values := []interface{10}

	row, err := dbh.One( sql, values... )
	if err != nil {
		log.Fatal(err)
	}

	if row.IsEmpty() {
		log.Println("no found")
	}

	name := row.Str("name") // str
	fl := row.Float("avarege") // float64
	id := row.Int("id") // int
	mdate := row.Date("date_created") // time
	mtime := row.Time("time_created") // time
	mdt := row.DateTime("date_time_created") // time
	inter := row.Interface() // map[string]interface{}
	col_names := row.Cols()  // map[]string
```

### Get list of rows ###
```go
	import "github.com/iostrovok/go-dbsearch/dbsearch"

	sql := "SELECT * FROM my_table WHERE id > $1 AND name <> $2"
	values := []interface{10, "murka"}

	list, err := dbh.List( sql, values... )
	if err != nil {
		log.Fatal(err)
	}

	if list.IsEmpty() {
		log.Println("no found")
	}

	count := list.Count()
	log.Printf("Total %d records\n", count)

	for i, row := list.Range(); row != nil && i != -1 {
		name := row.Str("name") // str
		fl := row.Float("avarege") // float64
		id := row.Int("id") // int
		mdate := row.Date("date_created") // time
		mtime := row.Time("time_created") // time
		mdt := row.DateTime("date_time_created") // time
		inter := row.Interface() // map[string]interface{}
		col_names := row.Cols()  // map[]string
	}


	// Data for json:
	json := list.Interface() // []map[string]interface {}

```

## xSql ##
### Import ###
```go
	import "github.com/iostrovok/go-dbsearch/dbsearch/xSql"
```
### Simple AND / OR ###
#### Example AND ####
```go
	sql, values := Select("public.mytable", "*").
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
	sql, values := Select("public.mytable", "*").
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
	where := Logic("AND").Mark("f_name", "=", "John").Mark("l_name", "=", "Lennon").
		Mark("age", "<", 40)

	sql, values := Update("public.mytable").
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
	sql, values := Delete("public.mytable").
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

