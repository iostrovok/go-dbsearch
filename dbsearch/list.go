package dbsearch

type List struct {
	dbh  *Searcher
	rows []*Row
	iter int
}

func (s *Searcher) List(mType *AllRows, sqlLine string, values ...[]interface{}) (*List, error) {

	list, err := s.Get(mType, sqlLine, values...)
	if err != nil {
		return nil, err
	}

	columns := _cols_to_hash(s.LastCols)
	out := List{s, []*Row{}, -1}
	for _, date := range list {
		r := s.SetRow(date, s.LastCols, columns)
		out.rows = append(out.rows, r)
	}
	return &out, nil
}

func (l *List) Interface() []map[string]interface{} {
	out := []map[string]interface{}{}

	for _, v := range l.rows {
		out = append(out, v.Interface())
	}

	return out
}

func (l *List) Count() int {
	return len(l.rows)
}

func (l *List) All() []*Row {
	return l.rows
}

func (l *List) Last() *Row {
	if 0 == len(l.rows) {
		return nil
	}

	return l.rows[len(l.rows)-1]
}

func (l *List) Fist() *Row {
	if 0 == len(l.rows) {
		return nil
	}

	return l.rows[0]
}

func (l *List) Row(i int) *Row {
	if i >= len(l.rows) {
		return nil
	}

	return l.rows[i]
}

func (l *List) Reset() {
	l.iter = -1
}

func (l *List) Next() *Row {

	l.iter++
	if l.iter >= len(l.rows) {
		return nil
	}

	return l.rows[l.iter]
}

func _cols_to_hash(cols []string) map[string]bool {

	r := map[string]bool{}

	for _, n := range cols {
		r[n] = true
	}
	return r
}
