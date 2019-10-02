package storage

import (
	"fmt"
	"strconv"
	"strings"
)

type Filter struct{}

func (i *Filter) Test(q Query, row Message) (bool, error) {
	if q.Operator == "" {
		return true, nil
	}

	if q.Operator == "OR" {
		if len(q.Childs) > 0 {
			for _, childQuery := range q.Childs {
				ok, err := i.Test(childQuery, row)

				if err != nil {
					return false, err
				}

				if ok {
					return true, nil
				}
			}
			return false, nil
		}
	}

	if q.Operator == "AND" {
		if len(q.Childs) > 0 {
			for _, childQuery := range q.Childs {
				ok, err := i.Test(childQuery, row)

				if err != nil {
					return false, err
				}

				if !ok {
					return false, nil
				}
			}
			return true, nil
		}
	}

	if q.Operator == "EXISTS" {
		_, exists := row[q.Field]
		return exists, nil
	}

	if q.Operator == "NOT_EXISTS" {
		_, exists := row[q.Field]
		return !exists, nil
	}

	if q.Operator == "LIKE" {
		if val, exists := row[q.Field]; exists {
			return strings.Contains(fmt.Sprintf("%v", val), q.Value), nil
		} else {
			return false, nil
		}
	}

	if q.Operator == "=" {
		if val, exists := row[q.Field]; exists {
			return fmt.Sprintf("%v", val) == q.Value, nil
		} else {
			return false, nil
		}
	}

	if q.Operator == "<>" {
		if val, exists := row[q.Field]; exists {
			return fmt.Sprintf("%v", val) != q.Value, nil
		} else {
			return true, nil
		}
	}

	if q.Operator == ">=" || q.Operator == ">" || q.Operator == "<=" || q.Operator == "<" {
		if val, exists := row[q.Field]; exists {
			// 1. оба значения (и проверяемое и ограничивающее) целые
			{
				inp, fil, err := castToInt64(val, q.Value)

				if err == nil {
					switch q.Operator {
					case ">=":
						return inp >= fil, nil
					case ">":
						return inp > fil, nil
					case "<=":
						return inp <= fil, nil
					case "<":
						return inp < fil, nil
					}
				}
			}

			// 2. пробуем привести к float
			{
				inp, fil, err := castToFloat64(val, q.Value)

				if err == nil {
					switch q.Operator {
					case ">=":
						return inp >= fil, nil
					case ">":
						return inp > fil, nil
					case "<=":
						return inp <= fil, nil
					case "<":
						return inp < fil, nil
					}
				}
			}
		}
		return false, nil
	}

	return false, fmt.Errorf("unknown operator %s", q.Operator)
}

func castToInt64(input interface{}, filter string) (int64, int64, error) {
	inp, errInput := interToInt64(input)
	fil, errFilter := strconv.ParseInt(filter, 10, 64)

	if errInput != nil {
		return 0, 0, errInput
	}

	if errFilter != nil {
		return 0, 0, errFilter
	}

	return inp, fil, nil
}

func castToFloat64(input interface{}, filter string) (float64, float64, error) {
	inp, errInput := interToFloat64(input)
	fil, errFilter := strconv.ParseFloat(filter, 64)

	if errInput != nil {
		return 0, 0, errInput
	}

	if errFilter != nil {
		return 0, 0, errFilter
	}

	return inp, fil, nil
}

func interToInt64(inp interface{}) (int64, error) {
	switch inp.(type) {
	case int:
		return int64(inp.(int)), nil
	case int8:
		return int64(inp.(int8)), nil
	case int16:
		return int64(inp.(int16)), nil
	case int32:
		return int64(inp.(int32)), nil
	case int64:
		return inp.(int64), nil
	case uint:
		return int64(inp.(uint)), nil
	case uint8:
		return int64(inp.(uint8)), nil
	case uint16:
		return int64(inp.(uint16)), nil
	case uint32:
		return int64(inp.(uint32)), nil
	case uint64:
		return int64(inp.(uint64)), nil
	case string:
		return strconv.ParseInt(inp.(string), 10, 64)
	default:
		return 0, fmt.Errorf("unsupported type %T", inp)
	}
}

func interToFloat64(inp interface{}) (float64, error) {
	switch inp.(type) {
	case float32:
		return float64(inp.(float32)), nil
	case float64:
		return inp.(float64), nil
	case string:
		return strconv.ParseFloat(inp.(string), 64)
	default:
		int, err := interToInt64(inp)
		if err == nil {
			return float64(int), nil
		} else {
			return 0, fmt.Errorf("unsupported type %T", inp)
		}
	}
}
