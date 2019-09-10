package storage

type QueryBuilder struct {
}

type Query struct {
	Childs   []Query
	Field    string
	Operator string
	Value    string
}

func (b QueryBuilder) And(filters ...Query) Query {
	return Query{
		Operator: "AND",
		Childs:   filters,
	}
}

func (b QueryBuilder) Or(filters ...Query) Query {
	return Query{
		Operator: "OR",
		Childs:   filters,
	}
}

func (b QueryBuilder) Equal(field string, value string) Query {
	return Query{
		Field:    field,
		Operator: "=",
		Value:    value,
	}
}

func (b QueryBuilder) NotEqual(field string, value string) Query {
	return Query{
		Field:    field,
		Operator: "<>",
		Value:    value,
	}
}

func (b QueryBuilder) Like(field string, value string) Query {
	return Query{
		Field:    field,
		Operator: "LIKE",
		Value:    value,
	}
}

func (b QueryBuilder) Exists(field string) Query {
	return Query{
		Field:    field,
		Operator: "EXISTS",
	}
}

func (b QueryBuilder) NotExists(field string) Query {
	return Query{
		Field:    field,
		Operator: "NOT_EXISTS",
	}
}

func (b QueryBuilder) Greater(field string, value string) Query {
	return Query{
		Field:    field,
		Operator: ">",
		Value:    value,
	}
}

func (b QueryBuilder) GreaterOrEqual(field string, value string) Query {
	return Query{
		Field:    field,
		Operator: ">=",
		Value:    value,
	}
}

func (b QueryBuilder) Less(field string, value string) Query {
	return Query{
		Field:    field,
		Operator: "<",
		Value:    value,
	}
}

func (b QueryBuilder) LessOrEqual(field string, value string) Query {
	return Query{
		Field:    field,
		Operator: "<=",
		Value:    value,
	}
}
