# Goshimon - Golang SQL query builder.


Goshimon allows you to fluently build SQL queries
using an expressive DSL which stays out of the way.

Bear in mind that this is WIP. API may change.

## Show me the code.

```go
qb := NewQueryBuilder().
    // select:
    Select("t1.id, t1.name").
    AddSelect("t2.title as some_title").
    // from, joins:
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
    // where:
    Where("t1.name = :arg1 AND t2.name = :arg2").
    AndWhere("t3.name = :arg3 AND t4.name in (:arg4)").
    // group by:
    GroupBy("t1.name, t2.name").
    AddGroupBy("t3.name, t4.name").
    // having:
    Having("count(distinct t5.name) > :arg5").
    AndHaving("max(t6.price) < :arg6").
    // order by:
    OrderBy("max(t6.price) asc, count(distinct t5.name) desc").
    AddOrderBy("t1.name").
    // limit and offset:
    SetFirstResult(30).
    SetMaxResults(15).
    // parameters:
    SetParameter("arg1", 1).
    SetParameter("arg2", "2").
    SetParameter("arg3", 3).
    SetParameter("arg4", []string{"4", "5"}). // yep, goshimon supports slices for IN parameters :)
    SetParameter("arg5", 5.5).
    SetParameter("arg6", 6.6)

sql := qb.GetSQL() // all named parameters are replaced with '$1'-style placeholders
args := qb.GetParameters()

// conn *pgx.Conn
rows, err := conn.Query(ctx, sql, args...)
```

### But will it blend?
Yep, subqueries are supported:
```go
qb1 := NewQueryBuilder().
    Select("u.id").
    From("users u").
    Where("u.name in (:names)").
    SetParameter("names", []string{"robert", "rob", "ken"})

qb2 := NewQueryBuilder().
    Select("a.number").
    From("apartments a").
    Where("a.floor = :floor").
    SetParameter("floor", 15).
    AndWhere("a.owner_id in(" + qb1.GetSQLWithNamedParams() + ")"). // nesty stuff
    MergeParametersFrom(qb1) // let's grab all the parameters too
```

### What else you got?

You may want to keep in mind we have these methods:
 - `CleanGroupBy()`
 - `CleanOrderBy()`

## Restrictions, plans?

1) Currently supporting only SELECT statements. We need basic INSERTs, UPDATEs, DELETEs.
2) No CTE's. Gotta figure out an elegant DSL for them.
3) `GetSQL()` substitutes named parameters with positional parameters in PostgreSQL format.
Gotta support other formats too.

## What was the name again?
"Goshimon" (「ご諮問」). Stands for "Your inquiry (sir/ma'am)".

### Inspiration?
Doctrine DBAL QueryBuilder.
