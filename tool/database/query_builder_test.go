package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder(t *testing.T) {
	qb := NewQueryBuilder()
	query := qb.Select("col1", "col2", "col3", "col4").
		From("test").
		Join(Inner, "test2", "test.test2_id = test2.id").
		Where("test.status", EqualThan, "active").
		Where("test2.id", EqualThan, 1).
		Build()

	assert.Equal(t, "SELECT col1, col2, col3, col4 FROM test INNER JOIN test2 ON test.test2_id = test2.id WHERE test.status = ? AND test2.id = ?", query.String())
	assert.Equal(t, []interface{}{"active", 1}, query.Args())
}

func TestQueryBuilderWhereSymbol(t *testing.T) {
	qb := NewQueryBuilder()
	query := qb.Select("col1", "col2", "col3", "col4").
		From("test").
		Join(Inner, "test2", "test.test2_id = test2.id").
		Where("test.status", EqualThan, "active").
		Where("test2.id", LessThan, 1).
		Where("test.status1", EqualOrLessThan, "active").
		Where("test3.id", GreaterThan, 1).
		Where("test.status2", EqualOrGreaterThan, "active").
		Build()

	assert.Equal(t, "SELECT col1, col2, col3, col4 FROM test INNER JOIN test2 ON test.test2_id = test2.id WHERE test.status = ? AND test2.id < ? AND test.status1 <= ? AND test3.id > ? AND test.status2 >= ?", query.String())
	assert.Equal(t, []interface{}{"active", 1, "active", 1, "active"}, query.Args())
}

func TestQueryBuilderWithEmptyArgs(t *testing.T) {
	qb := NewQueryBuilder()
	query := qb.Select("col1", "col2").
		From("test").
		Join(Left, "test2", "test.test2_id = test2.id").
		Where("test.status", EqualThan, "active").
		Where("empty", EqualThan, "").
		Build()

	assert.Equal(t, "SELECT col1, col2 FROM test LEFT JOIN test2 ON test.test2_id = test2.id WHERE test.status = ?", query.String())
	assert.Equal(t, []interface{}{"active"}, query.Args())
}

func TestQueryBuilderWithPointerArgs(t *testing.T) {
	arg := uint64(1)
	arg2 := int64(1)
	arg3 := float64(1)
	var arg4 *int64
	qb := NewQueryBuilder()
	query := qb.Select("col1", "col2").
		From("test").
		Join(Inner, "test2", "test.test2_id = test2.id").
		Where("test.status", EqualThan, "active").
		Where("pointer", EqualThan, &arg).
		Where("pointer2", EqualThan, &arg2).
		Where("pointer3", EqualThan, &arg3).
		Where("pointer4", EqualThan, arg4). //NIL POINTER
		Build()

	assert.Equal(t, "SELECT col1, col2 FROM test INNER JOIN test2 ON test.test2_id = test2.id WHERE test.status = ? AND pointer = ? AND pointer2 = ? AND pointer3 = ?", query.String())
	assert.Equal(t, []interface{}{"active", &arg, &arg2, &arg3}, query.Args())
}

func TestQueryBuilderWithLimit(t *testing.T) {
	qb := NewQueryBuilder()
	query := qb.Select("col1", "col2").
		From("test").
		Join(Right, "test2", "test.test2_id = test2.id").
		Where("test.status", EqualThan, "active").
		Limit("10", "2").
		Build()

	assert.Equal(t, "SELECT col1, col2 FROM test RIGHT JOIN test2 ON test.test2_id = test2.id WHERE test.status = ? LIMIT 2,10", query.String())
	assert.Equal(t, []interface{}{"active"}, query.Args())
}

func TestQueryBuilderWithGroupBy(t *testing.T) {
	qb := NewQueryBuilder()
	query := qb.Select("test.id", "col1", "col2").
		From("test").
		Join(Right, "test2", "test.test2_id = test2.id").
		Where("test.status", EqualThan, "active").
		GroupBy("test.id").
		Build()

	assert.Equal(t, "SELECT test.id, col1, col2 FROM test RIGHT JOIN test2 ON test.test2_id = test2.id WHERE test.status = ? GROUP BY test.id", query.String())
	assert.Equal(t, []interface{}{"active"}, query.Args())
}

func TestQueryBuilderWithOrderBy(t *testing.T) {
	qb := NewQueryBuilder()
	query := qb.Select("col1", "col2").
		From("test").
		Join(Full, "test2", "test.test2_id = test2.id").
		Where("test.status", EqualThan, "active").
		OrderBy("col1 ASC", "col2 DESC").
		Build()

	assert.Equal(t, "SELECT col1, col2 FROM test FULL JOIN test2 ON test.test2_id = test2.id WHERE test.status = ? ORDER BY col1 ASC, col2 DESC", query.String())
	assert.Equal(t, []interface{}{"active"}, query.Args())
}

func TestQueryBuilderWithInOperatorOfString(t *testing.T) {
	ids := []string{"id1", "id2", "id3"}
	qb := NewQueryBuilder()
	query := qb.Select("col1", "col2").
		From("test").
		Where("test.id", In, ids).
		Build()

	assert.Equal(t, "SELECT col1, col2 FROM test WHERE test.id IN(?,?,?)", query.String())
	assert.Equal(t, []interface{}{"id1", "id2", "id3"}, query.Args())
}

func TestQueryBuilderWithInOperatorOfInterface(t *testing.T) {
	ids := []interface{}{"1", "2", "3", "4"}
	qb := NewQueryBuilder()
	query := qb.Select("col1", "col2").
		From("test").
		Where("test.id", In, ids).
		Build()

	assert.Equal(t, "SELECT col1, col2 FROM test WHERE test.id IN(?,?,?,?)", query.String())
	assert.Equal(t, []interface{}{"1", "2", "3", "4"}, query.Args())
}

func TestQueryBuilderWithBooleanParam(t *testing.T) {
	qb := NewQueryBuilder()
	vFalse := false
	query := qb.Select("col1", "col2").
		From("test").
		Where("test.is", EqualThan, true).
		Where("test.is_not", EqualThan, &vFalse).
		Build()

	assert.Equal(t, "SELECT col1, col2 FROM test WHERE test.is = ? AND test.is_not = ?", query.String())
	assert.Equal(t, []interface{}{true, false}, query.Args())
}

func TestQueryBuilderWithCompositeIn(t *testing.T) {
	qb := NewQueryBuilder()
	query := qb.Select("col1", "col2").
		From("test").
		Where("(test.col1,test.col2)", In, [][]interface{}{{"1", "2"}, {"2", "3"}}).
		Build()

	assert.Equal(t, "SELECT col1, col2 FROM test WHERE (test.col1,test.col2) IN((?,?),(?,?))", query.String())
	assert.Equal(t, []interface{}{"1", "2", "2", "3"}, query.Args())
}
