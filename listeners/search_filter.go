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
		s = f.Field + f.Operator + f.Value
	}
	return s
}
