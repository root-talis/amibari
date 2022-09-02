package goshimon

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//
// clauses
//

func TestAllClauses(t *testing.T) {
	t.Log("Should generate correct SQL - using all clauses.")

	qb := NewQueryBuilder().
		// select
		Select("t1.id, t1.name").
		AddSelect("t2.title as some_title").
		// tables
		From("table1 t1").
		Join("table2 t2", "t2.id = t1.t2id").
		Join("table3 t3", "t3.id = t1.t3id").
		InnerJoin("table4 t4", "t4.id = t1.t4id").
		InnerJoin("table5 t5", "t5.id = t1.t5id").
		LeftJoin("table6 t6", "t6.id = t1.t6id").
		LeftJoin("table7 t7", "t7.id = t1.t7id").
		RightJoin("table8 t8", "t8.id = t1.t8id").
		RightJoin("table9 t9", "t9.id = t1.t9id").
		FullJoin("table10 t10", "t10.id = t1.t10id").
		FullJoin("table10 t10", "t10.id = t1.t10id").
		// where
		Where("t1.name = :arg1 AND t2.name = :arg2").
		AndWhere("t3.name = :arg3 AND t4.name = :arg4").
		// group by
		GroupBy("t1.name, t2.name").
		AddGroupBy("t3.name, t4.name").
		// having
		Having("count(distinct t5.name) > :arg5").
		AndHaving("max(t6.price) < :arg6").
		// order by
		OrderBy("max(t6.price) asc, count(distinct t5.name) desc").
		AddOrderBy("t1.name").
		// limit and offset
		SetFirstResult(30).
		SetMaxResults(15).
		// parameters
		SetParameter("arg1", 1).
		SetParameter("arg2", "2").
		SetParameter("arg3", 3).
		SetParameter("arg4", "4").
		SetParameter("arg5", 5.5).
		SetParameter("arg6", 6.6)

	sql := qb.GetSQL()
	args := qb.GetParameters()

	assert.Equal(t,
		// select
		"select t1.id, t1.name, t2.title as some_title "+
			// tables
			"from table1 t1 "+
			"join table2 t2 on (t2.id = t1.t2id) "+
			"join table3 t3 on (t3.id = t1.t3id) "+
			"inner join table4 t4 on (t4.id = t1.t4id) "+
			"inner join table5 t5 on (t5.id = t1.t5id) "+
			"left join table6 t6 on (t6.id = t1.t6id) "+
			"left join table7 t7 on (t7.id = t1.t7id) "+
			"right join table8 t8 on (t8.id = t1.t8id) "+
			"right join table9 t9 on (t9.id = t1.t9id) "+
			"full join table10 t10 on (t10.id = t1.t10id) "+
			"full join table10 t10 on (t10.id = t1.t10id) "+
			// where
			"where (t1.name = $1 and t2.name = $2) "+
			"and (t3.name = $3 and t4.name = $4) "+
			// group by
			"group by t1.name, t2.name, t3.name, t4.name "+
			// having
			"having (count(distinct t5.name) > $5) "+
			"and (max(t6.price) < $6) "+
			// order by
			"order by max(t6.price) asc, count(distinct t5.name) desc, t1.name "+
			// limit and offset
			"limit 15 offset 30",
		strings.ToLower(sql),
	)

	assert.Equal(
		t,
		[]interface{}{1, "2", 3, "4", 5.5, 6.6},
		args,
	)
}

func TestAllClausesReverseOrder(t *testing.T) {
	t.Log("Should generate correct SQL when building in reverse order.")

	qb := NewQueryBuilder().
		// parameters
		SetParameter("arg1", 1).
		SetParameter("arg2", "2").
		SetParameter("arg3", 3).
		SetParameter("arg4", "4").
		SetParameter("arg5", 5.5).
		SetParameter("arg6", 6.6).
		// limit and offset
		SetMaxResults(15).
		SetFirstResult(30).
		// order by
		OrderBy("max(t6.price) asc, count(distinct t5.name) desc").
		AddOrderBy("t1.name").
		// having
		Having("count(distinct t5.name) > :arg5").
		AndHaving("max(t6.price) < :arg6").
		// group by
		GroupBy("t1.name, t2.name").
		AddGroupBy("t3.name, t4.name").
		// where
		Where("t1.name = :arg1 AND t2.name = :arg2").
		AndWhere("t3.name = :arg3 AND t4.name = :arg4").
		// tables
		From("table1 t1").
		Join("table2 t2", "t2.id = t1.t2id").
		Join("table3 t3", "t3.id = t1.t3id").
		InnerJoin("table4 t4", "t4.id = t1.t4id").
		InnerJoin("table5 t5", "t5.id = t1.t5id").
		LeftJoin("table6 t6", "t6.id = t1.t6id").
		LeftJoin("table7 t7", "t7.id = t1.t7id").
		RightJoin("table8 t8", "t8.id = t1.t8id").
		RightJoin("table9 t9", "t9.id = t1.t9id").
		FullJoin("table10 t10", "t10.id = t1.t10id").
		FullJoin("table10 t10", "t10.id = t1.t10id").
		// select
		Select("t1.id, t1.name").
		AddSelect("t2.title as some_title")

	sql := qb.GetSQL()
	args := qb.GetParameters()

	assert.Equal(t,
		// select
		"select t1.id, t1.name, t2.title as some_title "+
			// tables
			"from table1 t1 "+
			"join table2 t2 on (t2.id = t1.t2id) "+
			"join table3 t3 on (t3.id = t1.t3id) "+
			"inner join table4 t4 on (t4.id = t1.t4id) "+
			"inner join table5 t5 on (t5.id = t1.t5id) "+
			"left join table6 t6 on (t6.id = t1.t6id) "+
			"left join table7 t7 on (t7.id = t1.t7id) "+
			"right join table8 t8 on (t8.id = t1.t8id) "+
			"right join table9 t9 on (t9.id = t1.t9id) "+
			"full join table10 t10 on (t10.id = t1.t10id) "+
			"full join table10 t10 on (t10.id = t1.t10id) "+
			// where
			"where (t1.name = $1 and t2.name = $2) "+
			"and (t3.name = $3 and t4.name = $4) "+
			// group by
			"group by t1.name, t2.name, t3.name, t4.name "+
			// having
			"having (count(distinct t5.name) > $5) "+
			"and (max(t6.price) < $6) "+
			// order by
			"order by max(t6.price) asc, count(distinct t5.name) desc, t1.name "+
			// limit and offset
			"limit 15 offset 30",
		strings.ToLower(sql),
	)

	assert.Equal(
		t,
		[]interface{}{1, "2", 3, "4", 5.5, 6.6},
		args,
	)
}

func TestClausesOverwrite(t *testing.T) {
	t.Log("Should overwrite SELECT, FROM, WHERE, ORDER BY, HAVING and GROUP BY on duplicate method calls.")

	qb := NewQueryBuilder().
		// select
		Select("t1.id, t1.name").
		// tables
		From("table1 t1").
		// where
		Where("t1.name = :arg1").
		AndWhere("t1.title = :arg2").
		// group by
		GroupBy("t1.name").
		AddGroupBy("t1.title").
		// having
		Having("t1.a > 2").
		AndHaving("t1.b < 4").
		// order by
		OrderBy("t1.c asc").
		AddOrderBy("t1.d desc").
		// parameters
		SetParameter("arg1", 1).
		SetParameter("arg2", "2")

	sql := qb.GetSQL()

	assert.Equal(t,
		"select t1.id, t1.name "+
			"from table1 t1 "+
			"where (t1.name = $1) and (t1.title = $2) "+
			"group by t1.name, t1.title "+
			"having (t1.a > 2) and (t1.b < 4) "+
			"order by t1.c asc, t1.d desc",
		strings.ToLower(sql),
	)

	// now overwrite

	qb.Select("distinct t1.name").
		From("(select 'hai' as name) t1").
		Where("t1.name = 'hai'").
		GroupBy("t1.name").
		Having("len(t1.name) < 15").
		OrderBy("t1.name")

	sql = qb.GetSQL()

	assert.Equal(t,
		"select distinct t1.name "+
			"from (select 'hai' as name) t1 "+
			"where (t1.name = 'hai') "+
			"group by t1.name "+
			"having (len(t1.name) < 15) "+
			"order by t1.name",
		strings.ToLower(sql),
	)
}

func TestClausesClean(t *testing.T) {
	t.Log("Should remove ORDER BY and GROUP BY on Clean* method calls.")

	qb := NewQueryBuilder().
		// select
		Select("t1.id, t1.name").
		// tables
		From("table1 t1").
		// where
		Where("t1.name = :arg1").
		AndWhere("t1.title = :arg2").
		// group by
		GroupBy("t1.name").
		AddGroupBy("t1.title").
		// having
		Having("t1.a > 2").
		AndHaving("t1.b < 4").
		// order by
		OrderBy("t1.c asc").
		AddOrderBy("t1.d desc").
		// parameters
		SetParameter("arg1", 1).
		SetParameter("arg2", "2")

	sql := qb.GetSQL()

	assert.Equal(t,
		"select t1.id, t1.name "+
			"from table1 t1 "+
			"where (t1.name = $1) and (t1.title = $2) "+
			"group by t1.name, t1.title "+
			"having (t1.a > 2) and (t1.b < 4) "+
			"order by t1.c asc, t1.d desc",
		strings.ToLower(sql),
	)

	// now clean order by

	qb.CleanOrderBy()

	sql = qb.GetSQL()

	assert.Equal(t,
		"select t1.id, t1.name "+
			"from table1 t1 "+
			"where (t1.name = $1) and (t1.title = $2) "+
			"group by t1.name, t1.title "+
			"having (t1.a > 2) and (t1.b < 4)",
		strings.ToLower(sql),
	)

	// now clean group by

	qb.CleanGroupBy()

	sql = qb.GetSQL()

	assert.Equal(t,
		"select t1.id, t1.name "+
			"from table1 t1 "+
			"where (t1.name = $1) and (t1.title = $2) "+
			"having (t1.a > 2) and (t1.b < 4)",
		strings.ToLower(sql),
	)
}

//
// parameters
//

func TestUnsetParameter(t *testing.T) {
	t.Log("Should substitute array parameter with correct number of templates")

	qb := NewQueryBuilder().
		Select("distinct u.id").
		From("users u").
		Where("u.age > :min_age").
		SetParameter("min_age", 2)

	sql := qb.GetSQL()
	args := qb.GetParameters()

	assert.Equal(t,
		"select distinct u.id "+
			"from users u "+
			"where (u.age > $1)",
		strings.ToLower(sql),
	)

	assert.Equal(t,
		[]interface{}{2},
		args,
	)

	// now unset

	qb.UnsetParameter("min_age")

	sql = qb.GetSQL()
	args = qb.GetParameters()

	assert.Equal(t,
		"select distinct u.id "+
			"from users u "+
			"where (u.age > :min_age)",
		strings.ToLower(sql),
	)

	assert.Equal(t,
		[]interface{}{},
		args,
	)
}

func TestArrayParameter(t *testing.T) {
	t.Log("Should substitute array parameter with correct number of templates")

	qb := NewQueryBuilder().
		Select("distinct u.id").
		From("users u").
		Join("roles r", "r.user_id = u.id").
		Where("u.name like :name").
		AndWhere("r.name in (:roles)").
		AndWhere("u.is_deleted = :is_deleted").
		SetParameter("name", "%doe%").
		SetParameter("roles", []string{"admin", "manager", "script", "root"}).
		SetParameter("is_deleted", false)

	sql := qb.GetSQL()
	args := qb.GetParameters()

	assert.Equal(t,
		"select distinct u.id "+
			"from users u "+
			"join roles r on (r.user_id = u.id) "+
			"where (u.name like $1)"+
			" and (r.name in ($2,$3,$4,$5))"+
			" and (u.is_deleted = $6)",
		strings.ToLower(sql),
	)

	assert.Equal(t,
		[]interface{}{"%doe%", "admin", "manager", "script", "root", false},
		args,
	)
}

//
// merge
//

func TestMergeParameters(t *testing.T) {
	t.Log("Should correctly copy parameters from another qb")

	qb1 := NewQueryBuilder().
		Select("u.id").From("users u").Where("u.name in (:names)").
		SetParameter("names", []string{"robert", "rob", "ken"})

	qb2 := NewQueryBuilder().
		Select("a.number").From("apartments a").
		Where("a.floor = :floor").
		SetParameter("floor", 2).
		AndWhere("a.owner_id in(" + qb1.GetSQLWithNamedParams() + ")").
		MergeParametersFrom(qb1)

	assert.Equal(t,
		"select a.number from apartments a "+
			"where (a.floor = $1) "+
			"and (a.owner_id in(select u.id from users u where (u.name in ($2,$3,$4))))",
		strings.ToLower(qb2.GetSQL()),
	)

	assert.Equal(t,
		[]interface{}{2, "robert", "rob", "ken"},
		qb2.GetParameters(),
	)
}
