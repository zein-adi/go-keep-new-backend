package helpers_mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"strings"
)

func NewQueryBuilder(ctx context.Context, db *sql.DB, table string) *QueryBuilder {
	return &QueryBuilder{
		table: table,
		skip:  0,
		take:  0,
		db:    db,
		ctx:   ctx,
	}
}

type QueryBuilder struct {
	table   string
	fields  []string
	wheres  []Where
	joins   []Join
	groupBy []string
	orderBy []string
	skip    int
	take    int
	db      *sql.DB
	ctx     context.Context
}
type Where struct {
	whereType string
	field     string
	operator  string
	argument  any
	isRaw     bool
	wheres    []Where
}
type Join struct {
	table    string
	joinType string
	wheres   []Where
}

/*
 * Closing Methods
 */

func (x *QueryBuilder) Get() (rows *sql.Rows, cleanup func()) {
	where, arguments := x.renderWhere()
	query := fmt.Sprintf("SELECT %s FROM %s %s %s %s %s %s",
		strings.Join(x.fields, ","),
		x.table,
		x.renderJoin(),
		where,
		x.renderGroupBy(),
		x.renderOrderBy(),
		x.renderLimit(),
	)
	rows, err := x.db.QueryContext(x.ctx, query, arguments...)
	helpers_error.PanicIfError(err)
	cleanup = func() {
		helpers_error.PanicIfError(rows.Close())
	}
	return rows, cleanup
}
func (x *QueryBuilder) Count() (count int) {
	where, arguments := x.renderWhere()
	query := fmt.Sprintf("SELECT COUNT(0) FROM %s %s %s %s %s",
		x.table,
		x.renderJoin(),
		where,
		x.renderGroupBy(),
		x.renderLimit(),
	)
	rows, err := x.db.QueryContext(x.ctx, query, arguments...)
	defer func(rows *sql.Rows) {
		helpers_error.PanicIfError(rows.Close())
	}(rows)

	helpers_error.PanicIfError(err)
	if !rows.Next() {
		panic(errors.New("no row available"))
	}
	helpers_error.PanicIfError(rows.Scan(&count))
	return count
}
func (x *QueryBuilder) Insert(model map[string]any) (lastId int, err error) {
	execContext, err := x.insertBatch(model)
	if err != nil {
		return 0, err
	}

	id, err := execContext.LastInsertId()
	helpers_error.PanicIfError(err)
	return int(id), err
}
func (x *QueryBuilder) InsertBatch(models ...map[string]any) (affected int, err error) {
	execContext, err := x.insertBatch(models...)
	if err != nil {
		return 0, err
	}

	aff, err := execContext.RowsAffected()
	helpers_error.PanicIfError(err)
	return int(aff), err
}
func (x *QueryBuilder) insertBatch(models ...map[string]any) (sql.Result, error) {
	model := models[0]
	var fields []string
	for field := range model {
		fields = append(fields, field)
	}

	var values []any
	var placeholders []string
	for _, m := range models {
		var placeholder []string
		for _, field := range fields {
			value := m[field]
			values = append(values, value)
			placeholder = append(placeholder, "?")
		}
		placeholders = append(placeholders, strings.Join(placeholder, ","))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		x.table,
		strings.Join(fields, ","),
		strings.Join(placeholders, "),("),
	)
	return x.db.ExecContext(x.ctx, query, values...)
}
func (x *QueryBuilder) Update(model map[string]any) (affected int) {
	var arguments []any

	var fields []string
	for field := range model {
		fields = append(fields, field)
	}

	var values []any
	var placeholder []string
	for _, field := range fields {
		value := model[field]
		values = append(values, value)
		placeholder = append(placeholder, field+" = ?")
	}
	arguments = append(arguments, values...)

	where, whereArguments := x.renderWhere()
	arguments = append(arguments, whereArguments...)

	query := fmt.Sprintf("UPDATE %s %s SET %s %s",
		x.table,
		x.renderJoin(),
		strings.Join(placeholder, ","),
		where,
	)
	execContext, err := x.db.ExecContext(x.ctx, query, arguments...)
	helpers_error.PanicIfError(err)
	rowsAffected, err := execContext.RowsAffected()
	helpers_error.PanicIfError(err)
	return int(rowsAffected)
}
func (x *QueryBuilder) Delete() (affected int) {
	where, arguments := x.renderWhere()
	query := fmt.Sprintf("DELETE %s FROM %s %s %s",
		x.table,
		x.table,
		x.renderJoin(),
		where,
	)
	execContext, err := x.db.ExecContext(x.ctx, query, arguments...)
	helpers_error.PanicIfError(err)
	rowsAffected, err := execContext.RowsAffected()
	helpers_error.PanicIfError(err)
	return int(rowsAffected)
}

/*
 * Renders
 */

func (x *QueryBuilder) renderWhere() (where string, arguments []any) {
	if len(x.wheres) == 0 {
		return "", nil
	}
	where, arguments = x.renderSubWhere(x.wheres)
	return "WHERE " + where, arguments
}
func (x *QueryBuilder) renderSubWhere(wheres []Where) (where string, arguments []any) {
	if len(wheres) == 0 {
		return "", nil
	}

	arguments = []any{}
	where = ""
	for i, w := range wheres {
		whereType := w.whereType
		if i == 0 {
			whereType = ""
		}

		if len(w.wheres) > 0 {
			subWhere, subArguments := x.renderSubWhere(w.wheres)
			where += fmt.Sprintf(" %s (%s)", whereType, subWhere)
			arguments = append(arguments, subArguments...)
			continue
		}

		where += fmt.Sprintf(" %s %s %s", whereType, w.field, w.operator)
		if w.operator == "in" {
			inArguments := w.argument.([]any)
			where += " (" + strings.Repeat("?", len(inArguments)) + ")"
			arguments = append(arguments, inArguments...)
		} else {
			if w.isRaw {
				where += fmt.Sprintf(" %s", w.argument)
			} else {
				where += " ?"
				arguments = append(arguments, w.argument)
			}
		}
	}
	return where, arguments
}
func (x *QueryBuilder) renderJoin() string {
	if len(x.joins) == 0 {
		return ""
	}
	join := ""
	for _, j := range x.joins {
		join += fmt.Sprintf(" %s JOIN %s ON", j.joinType, j.table)
		for _, w := range j.wheres {
			join += fmt.Sprintf(" %s %s %s %s", w.whereType, w.field, w.operator, w.argument)
		}
	}
	return join
}
func (x *QueryBuilder) renderGroupBy() string {
	if len(x.groupBy) == 0 {
		return ""
	}
	return " GROUP BY " + strings.Join(x.groupBy, ",")
}
func (x *QueryBuilder) renderOrderBy() string {
	if len(x.orderBy) == 0 {
		return ""
	}
	return " ORDER BY " + strings.Join(x.orderBy, ",")
}
func (x *QueryBuilder) renderLimit() string {
	if x.skip == 0 && x.take == 0 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d OFFSET %d", x.take, x.skip)
}

/*
 * Chaining Methods
 */

func (x *QueryBuilder) Select(fields ...string) {
	for _, f := range fields {
		x.fields = append(x.groupBy, f)
	}
}
func (x *QueryBuilder) GroupBy(groupBy ...string) {
	for _, g := range groupBy {
		x.groupBy = append(x.groupBy, g)
	}
}
func (x *QueryBuilder) OrderBy(orderBy ...string) {
	for _, o := range orderBy {
		x.orderBy = append(x.groupBy, o)
	}
}
func (x *QueryBuilder) Skip(skip int) {
	x.skip = skip
}
func (x *QueryBuilder) Take(take int) {
	x.take = take
}
func (x *QueryBuilder) Where(field string, operator string, argument any) {
	x.wheres = append(x.wheres, Where{
		whereType: "AND",
		field:     field,
		operator:  operator,
		argument:  argument,
		isRaw:     false,
		wheres:    nil,
	})
}
func (x *QueryBuilder) OrWhere(field string, operator string, argument any) {
	x.wheres = append(x.wheres, Where{
		whereType: "OR",
		field:     field,
		operator:  operator,
		argument:  argument,
		isRaw:     false,
		wheres:    nil,
	})
}
func (x *QueryBuilder) WhereRaw(field string, operator string, argument any) {
	x.wheres = append(x.wheres, Where{
		whereType: "AND",
		field:     field,
		operator:  operator,
		argument:  argument,
		isRaw:     true,
		wheres:    nil,
	})
}
func (x *QueryBuilder) OrWhereRaw(field string, operator string, argument any) {
	x.wheres = append(x.wheres, Where{
		whereType: "OR",
		field:     field,
		operator:  operator,
		argument:  argument,
		isRaw:     true,
		wheres:    nil,
	})
}
func (x *QueryBuilder) WhereIn(field string, arguments []string) {
	var inArguments []any
	for _, inArg := range arguments {
		inArguments = append(inArguments, inArg)
	}

	x.wheres = append(x.wheres, Where{
		whereType: "AND",
		field:     field,
		operator:  "in",
		argument:  inArguments,
		isRaw:     false,
		wheres:    nil,
	})
}
