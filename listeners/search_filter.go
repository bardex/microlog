package listeners

type SearchBuilder struct {
}

type SearchFilter struct {
	ChildsUnion string
	Childs      []SearchFilter
	Field       string
	Operator    string
	Value       string
}

func (b SearchBuilder) And(filters ...SearchFilter) SearchFilter {
	return SearchFilter{
		ChildsUnion: "AND",
		Childs:      filters,
	}
}

func (b SearchBuilder) Or(filters ...SearchFilter) SearchFilter {
	return SearchFilter{
		ChildsUnion: "OR",
		Childs:      filters,
	}
}

func (b SearchBuilder) Equal(field string, value string) SearchFilter {
	return SearchFilter{
		Field:    field,
		Operator: "=",
		Value:    value,
	}
}

func (b SearchBuilder) NotEqual(field string, value string) SearchFilter {
	return SearchFilter{
		Field:    field,
		Operator: "<>",
		Value:    value,
	}
}

func (b SearchBuilder) Like(field string, value string) SearchFilter {
	return SearchFilter{
		Field:    field,
		Operator: "LIKE",
		Value:    value,
	}
}

func (b SearchBuilder) Exists(field string) SearchFilter {
	return SearchFilter{
		Field:    field,
		Operator: "EXISTS",
	}
}

func (b SearchBuilder) NotExists(field string) SearchFilter {
	return SearchFilter{
		Field:    field,
		Operator: "NOT_EXISTS",
	}
}

func (b SearchBuilder) Greater(field string, value string) SearchFilter {
	return SearchFilter{
		Field:    field,
		Operator: ">",
		Value:    value,
	}
}

func (b SearchBuilder) GreaterOrEqual(field string, value string) SearchFilter {
	return SearchFilter{
		Field:    field,
		Operator: ">=",
		Value:    value,
	}
}

func (b SearchBuilder) Less(field string, value string) SearchFilter {
	return SearchFilter{
		Field:    field,
		Operator: "<",
		Value:    value,
	}
}

func (b SearchBuilder) LessOrEqual(field string, value string) SearchFilter {
	return SearchFilter{
		Field:    field,
		Operator: "<=",
		Value:    value,
	}
}

func (f SearchFilter) String() string {
	s := ""
	if f.ChildsUnion != "" {
		if len(f.Childs) > 0 {
			s = "("
			sep := ""
			for _, c := range f.Childs {
				s = s + sep + c.String()
				sep = " " + f.ChildsUnion + " "
			}
			s = s + ")"
		}
	} else {
		if f.Value != "" {
			s = f.Field + " " + f.Operator + " " + f.Value
		} else {
			s = f.Operator + " " + f.Field
		}
	}
	return s
}
