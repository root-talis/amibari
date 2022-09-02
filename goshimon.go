package goshimon

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/wk8/go-ordered-map/v2"
)

type join struct {
	joinType  string
	tableName string
	onClause  string
}

type QueryBuilder struct {
	sqlSelect    []string
	sqlFrom      string
	sqlJoins     []join
	sqlAndWhere  []string
	sqlGroupBy   []string
	sqlAndHaving []string
	sqlOrderBy   []string

	maxResults  uint
	firstResult uint

	args *orderedmap.OrderedMap[string, interface{}]
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		sqlSelect:    make([]string, 0, 8),
		sqlAndWhere:  make([]string, 0, 8),
		sqlJoins:     make([]join, 0, 4),
		sqlGroupBy:   make([]string, 0, 8),
		sqlAndHaving: make([]string, 0, 8),
		sqlOrderBy:   make([]string, 0, 8),
		args:         orderedmap.New[string, interface{}](),
	}
}

//
// select
//

func (qb *QueryBuilder) Select(sql string) *QueryBuilder {
	if len(qb.sqlSelect) > 0 {
		qb.sqlSelect = make([]string, 0, 8)
	}

	qb.sqlSelect = append(qb.sqlSelect, sql)
	return qb
}

func (qb *QueryBuilder) AddSelect(sql string) *QueryBuilder {
	qb.sqlSelect = append(qb.sqlSelect, sql)
	return qb
}

//
// from
//

func (qb *QueryBuilder) From(sql string) *QueryBuilder {
	qb.sqlFrom = sql
	return qb
}

//
// join
//

func (qb *QueryBuilder) Join(table string, on string) *QueryBuilder {
	qb.sqlJoins = append(qb.sqlJoins, join{
		joinType:  "",
		tableName: table,
		onClause:  on,
	})
	return qb
}

func (qb *QueryBuilder) InnerJoin(table string, on string) *QueryBuilder {
	qb.sqlJoins = append(qb.sqlJoins, join{
		joinType:  "INNER ",
		tableName: table,
		onClause:  on,
	})
	return qb
}

func (qb *QueryBuilder) LeftJoin(table string, on string) *QueryBuilder {
	qb.sqlJoins = append(qb.sqlJoins, join{
		joinType:  "LEFT ",
		tableName: table,
		onClause:  on,
	})
	return qb
}

func (qb *QueryBuilder) RightJoin(table string, on string) *QueryBuilder {
	qb.sqlJoins = append(qb.sqlJoins, join{
		joinType:  "RIGHT ",
		tableName: table,
		onClause:  on,
	})
	return qb
}

func (qb *QueryBuilder) FullJoin(table string, on string) *QueryBuilder {
	qb.sqlJoins = append(qb.sqlJoins, join{
		joinType:  "FULL ",
		tableName: table,
		onClause:  on,
	})
	return qb
}

//
// where
//

func (qb *QueryBuilder) Where(sql string) *QueryBuilder {
	if len(qb.sqlAndWhere) > 0 {
		qb.sqlAndWhere = make([]string, 0, 8)
	}

	qb.sqlAndWhere = append(qb.sqlAndWhere, "("+sql+")")
	return qb
}

func (qb *QueryBuilder) AndWhere(sql string) *QueryBuilder {
	qb.sqlAndWhere = append(qb.sqlAndWhere, "("+sql+")")
	return qb
}

//
// group by
//

func (qb *QueryBuilder) GroupBy(sql string) *QueryBuilder {
	if len(qb.sqlGroupBy) > 0 {
		qb.sqlGroupBy = make([]string, 0, 8)
	}

	qb.sqlGroupBy = append(qb.sqlGroupBy, sql)
	return qb
}

func (qb *QueryBuilder) AddGroupBy(sql string) *QueryBuilder {
	qb.sqlGroupBy = append(qb.sqlGroupBy, sql)
	return qb
}

func (qb *QueryBuilder) CleanGroupBy() *QueryBuilder {
	qb.sqlGroupBy = make([]string, 0, 8)
	return qb
}

//
// having
//

func (qb *QueryBuilder) Having(sql string) *QueryBuilder {
	if len(qb.sqlAndHaving) > 0 {
		qb.sqlAndHaving = make([]string, 0, 8)
	}

	qb.sqlAndHaving = append(qb.sqlAndHaving, "("+sql+")")
	return qb
}

func (qb *QueryBuilder) AndHaving(sql string) *QueryBuilder {
	qb.sqlAndHaving = append(qb.sqlAndHaving, "("+sql+")")
	return qb
}

//
// order by
//

func (qb *QueryBuilder) OrderBy(sql string) *QueryBuilder {
	if len(qb.sqlOrderBy) > 0 {
		qb.sqlOrderBy = make([]string, 0, 8)
	}

	qb.sqlOrderBy = append(qb.sqlOrderBy, sql)
	return qb
}

func (qb *QueryBuilder) AddOrderBy(sql string) *QueryBuilder {
	qb.sqlOrderBy = append(qb.sqlOrderBy, sql)
	return qb
}

func (qb *QueryBuilder) CleanOrderBy() *QueryBuilder {
	qb.sqlOrderBy = make([]string, 0, 8)
	return qb
}

//
// offset, limit
//

func (qb *QueryBuilder) SetMaxResults(value uint) *QueryBuilder {
	qb.maxResults = value
	return qb
}

func (qb *QueryBuilder) SetFirstResult(value uint) *QueryBuilder {
	qb.firstResult = value
	return qb
}

//
// parameters
//

func (qb *QueryBuilder) SetParameter(name string, value interface{}) *QueryBuilder {
	qb.args.Set(name, value)
	return qb
}

func (qb *QueryBuilder) UnsetParameter(name string) *QueryBuilder {
	qb.args.Delete(name)
	return qb
}

func (qb *QueryBuilder) MergeParametersFrom(other *QueryBuilder) *QueryBuilder {
	for el := other.args.Oldest(); el != nil; el = el.Next() {
		qb.SetParameter(el.Key, el.Value)
	}

	return qb
}

//
// sql output
//

func (qb *QueryBuilder) GetSQL() string {
	query := qb.GetSQLWithNamedParams()

	//
	// replace named parameters with numbered:
	//

	var paramIdx uint
	for el := qb.args.Oldest(); el != nil; el = el.Next() {
		r := reflect.ValueOf(el.Value)

		switch r.Kind() {
		case reflect.Slice, reflect.Array:
			arr := make([]string, r.Len())
			for i := 0; i < r.Len(); i++ {
				paramIdx++
				arr[i] = fmt.Sprintf("$%d", paramIdx)
			}

			query = strings.Replace(query, ":"+el.Key, strings.Join(arr, ","), -1)
		default:
			paramIdx++
			query = strings.Replace(query, ":"+el.Key, fmt.Sprintf("$%d", paramIdx), -1)
		}
	}

	return query
}

func (qb *QueryBuilder) GetSQLWithNamedParams() string {
	const clausesCount = 9
	clauses := make([]string, 0, clausesCount)

	clauses = append(clauses, "SELECT "+strings.Join(qb.sqlSelect, ", "))
	clauses = append(clauses, "FROM "+qb.sqlFrom)

	if len(qb.sqlJoins) != 0 {
		joins := make([]string, 0, len(qb.sqlJoins))
		for _, v := range qb.sqlJoins {
			joins = append(joins, v.joinType+"JOIN "+v.tableName+" ON ("+v.onClause+")")
		}
		clauses = append(clauses, strings.Join(joins, " "))
	}

	if len(qb.sqlAndWhere) != 0 {
		clauses = append(clauses, "WHERE "+strings.Join(qb.sqlAndWhere, " AND "))
	}

	if len(qb.sqlGroupBy) != 0 {
		clauses = append(clauses, "GROUP BY "+strings.Join(qb.sqlGroupBy, ", "))
	}

	if len(qb.sqlAndHaving) != 0 {
		clauses = append(clauses, "HAVING "+strings.Join(qb.sqlAndHaving, " AND "))
	}

	if len(qb.sqlOrderBy) != 0 {
		clauses = append(clauses, "ORDER BY "+strings.Join(qb.sqlOrderBy, ", "))
	}

	if qb.maxResults != 0 {
		clauses = append(clauses, fmt.Sprintf("LIMIT %d", qb.maxResults))
	}

	if qb.firstResult != 0 {
		clauses = append(clauses, fmt.Sprintf("OFFSET %d", qb.firstResult))
	}

	return strings.Join(clauses, " ")
}

func (qb *QueryBuilder) GetParameters() []interface{} {
	result := make([]interface{}, 0, qb.args.Len())

	for el := qb.args.Oldest(); el != nil; el = el.Next() {
		r := reflect.ValueOf(el.Value)

		switch r.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < r.Len(); i++ {
				result = append(result, r.Index(i).Interface())
			}
		default:
			result = append(result, el.Value)
		}
	}

	return result
}
