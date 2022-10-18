package database

import (
	"fmt"
	"strings"
)

type Query struct {
	queryString string
	args        []interface{}
}

func (q *Query) String() string {
	return q.queryString
}

func (q *Query) Args() []interface{} {
	return q.args
}

type QueryBuilder interface {
	Select(args ...string) QueryBuilder
	From(from string) QueryBuilder
	Where(arg string, symbol Symbol, value interface{}) QueryBuilder
	Join(typ JoinType, to string, condition string) QueryBuilder
	OrderBy(fields ...string) QueryBuilder
	GroupBy(fields ...string) QueryBuilder
	Limit(limit, offset string) QueryBuilder
	Build() Query
}

type queryBuilder struct {
	sel     []string
	where   []string
	join    []string
	from    string
	query   string
	orderBy []string
	groupBy []string
	limit   string
	args    []interface{}
}

func NewQueryBuilder() QueryBuilder {
	return &queryBuilder{}
}

func (qb *queryBuilder) Select(args ...string) QueryBuilder {
	qb.sel = append(qb.sel, args...)
	return qb
}

func (qb *queryBuilder) From(from string) QueryBuilder {
	qb.from = from
	return qb
}

type Symbol string

const (
	EqualThan          Symbol = "="
	DistinctThan       Symbol = "!="
	LessThan           Symbol = "<"
	GreaterThan        Symbol = ">"
	EqualOrLessThan    Symbol = "<="
	EqualOrGreaterThan Symbol = ">="
	In                 Symbol = "IN"
)

func (qb *queryBuilder) Where(arg string, symbol Symbol, value interface{}) QueryBuilder {
	switch val := value.(type) {
	case string:
		if val != "" {
			qb.where = append(qb.where, fmt.Sprintf("%v %v ?", arg, symbol))
			qb.args = append(qb.args, val)
		}
	case *string:
		if val != nil && *val != "" {
			qb.where = append(qb.where, fmt.Sprintf("%v %v ?", arg, symbol))
			qb.args = append(qb.args, val)
		}
	case *int64:
		if val != nil {
			qb.where = append(qb.where, fmt.Sprintf("%v %v ?", arg, symbol))
			qb.args = append(qb.args, val)
		}
	case *uint64:
		if val != nil {
			qb.where = append(qb.where, fmt.Sprintf("%v %v ?", arg, symbol))
			qb.args = append(qb.args, val)
		}
	case *bool:
		if val != nil {
			qb.where = append(qb.where, fmt.Sprintf("%v %v ?", arg, symbol))
			qb.args = append(qb.args, *val)
		}
	case bool:
		qb.where = append(qb.where, fmt.Sprintf("%v %v ?", arg, symbol))
		qb.args = append(qb.args, val)
	case *float64:
		if val != nil {
			qb.where = append(qb.where, fmt.Sprintf("%v %v ?", arg, symbol))
			qb.args = append(qb.args, val)
		}
	case []string:
		if len(val) > 0 && symbol == In {
			gArr := make([]interface{}, len(val))
			for i, v := range val {
				gArr[i] = v
			}
			placeHolder := "?" + strings.Repeat(",?", len(val)-1)
			qb.where = append(qb.where, fmt.Sprintf("%v %v(%v)", arg, symbol, placeHolder))
			qb.args = append(qb.args, gArr...)
		}
	case []interface{}:
		if len(val) > 0 && symbol == In {
			placeHolder := "?" + strings.Repeat(",?", len(val)-1)
			qb.where = append(qb.where, fmt.Sprintf("%v %v(%v)", arg, symbol, placeHolder))
			qb.args = append(qb.args, val...)
		}
	case [][]interface{}:
		if len(val) > 0 && symbol == In {
			singlePlaceholder := "(?" + strings.Repeat(",?", len(val[0])-1) + ")"
			placeHolder := singlePlaceholder + strings.Repeat(","+singlePlaceholder, len(val)-1)
			qb.where = append(qb.where, fmt.Sprintf("%v %v(%v)", arg, symbol, placeHolder))
			for _, a := range val {
				qb.args = append(qb.args, a...)
			}
		}
	default:
		qb.where = append(qb.where, fmt.Sprintf("%v %v ?", arg, symbol))
		qb.args = append(qb.args, val)
	}
	return qb
}

type JoinType string

const (
	Inner JoinType = "INNER JOIN"
	Left  JoinType = "LEFT JOIN"
	Right JoinType = "RIGHT JOIN"
	Full  JoinType = "FULL JOIN"
)

func (qb *queryBuilder) Join(typ JoinType, to string, condition string) QueryBuilder {
	qb.join = append(qb.join, fmt.Sprintf("%s %s ON %s", typ, to, condition))
	return qb
}

func (qb *queryBuilder) OrderBy(fields ...string) QueryBuilder {
	qb.orderBy = append(qb.orderBy, fields...)
	return qb
}

func (qb *queryBuilder) GroupBy(fields ...string) QueryBuilder {
	qb.groupBy = append(qb.groupBy, fields...)
	return qb
}

func (qb *queryBuilder) Limit(limit, offset string) QueryBuilder {
	qb.limit = strings.Join([]string{offset, limit}, ",")
	return qb
}

func (qb *queryBuilder) Build() Query {
	qb.query = "SELECT " + strings.Join(qb.sel, ", ")
	qb.query += " FROM " + qb.from
	if len(qb.join) > 0 {
		qb.query += " " + strings.Join(qb.join, " ")
	}
	if len(qb.where) > 0 {
		qb.query += " WHERE " + strings.Join(qb.where, " AND ")
	}
	if len(qb.groupBy) > 0 {
		qb.query += " GROUP BY " + strings.Join(qb.groupBy, ", ")
	}
	if len(qb.orderBy) > 0 {
		qb.query += " ORDER BY " + strings.Join(qb.orderBy, ", ")
	}
	if qb.limit != "" {
		qb.query += " LIMIT " + qb.limit
	}
	return Query{queryString: qb.query, args: qb.args}
}
