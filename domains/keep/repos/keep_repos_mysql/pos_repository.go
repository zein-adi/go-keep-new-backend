package keep_repos_mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"strconv"
)

var posEntityName = "pos"
var posTableName = "keep_pos"

func NewPosMySqlRepository() *PosMysqlRepository {
	db, dbCleanup := helpers_mysql.OpenMySqlConnection()
	return &PosMysqlRepository{
		db:        db,
		dbCleanup: dbCleanup,
	}
}

type PosMysqlRepository struct {
	db        *sql.DB
	dbCleanup func()
}

func (x *PosMysqlRepository) Get(ctx context.Context, request *keep_request.GetPos) []*keep_entities.Pos {
	q := x.newQueryRequest(ctx, "aktif", request)
	return x.newEntitiesFromRows(q.Get())
}
func (x *PosMysqlRepository) FindById(ctx context.Context, id string) (*keep_entities.Pos, error) {
	q := x.newQueryRequest(ctx, "aktif", keep_request.NewGetPos())
	q.Where("id", "=", id)
	models := x.newEntitiesFromRows(q.Get())
	if len(models) == 0 {
		return nil, helpers_error.NewEntryNotFoundError(posEntityName, "id", id)
	}
	return models[0], nil
}
func (x *PosMysqlRepository) Insert(ctx context.Context, pos *keep_entities.Pos) (*keep_entities.Pos, error) {
	model := &keep_entities.Pos{}

	isShowInt := 0
	if pos.IsShow {
		isShowInt = 1
	}
	isLeafInt := 0
	if pos.IsLeaf {
		isLeafInt = 1
	}
	var parentId *string = nil
	if pos.ParentId != "" {
		parentId = &pos.ParentId
	}

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	insertId, err := q.Insert(map[string]any{
		"nama":      pos.Nama,
		"urutan":    pos.Urutan,
		"saldo":     pos.Saldo,
		"parent_id": parentId,
		"level":     pos.Level,
		"is_show":   isShowInt,
		"is_leaf":   isLeafInt,
		"status":    pos.Status,
	})
	if err != nil {
		return model, err
	}
	model = pos.Copy()
	model.Id = strconv.Itoa(insertId)
	return model, nil
}
func (x *PosMysqlRepository) Update(ctx context.Context, pos *keep_entities.Pos) (affected int, err error) {
	isShowInt := 0
	if pos.IsShow {
		isShowInt = 1
	}
	isLeafInt := 0
	if pos.IsLeaf {
		isLeafInt = 1
	}
	var parentId *string = nil
	if pos.ParentId != "" {
		parentId = &pos.ParentId
	}

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Where("id", "=", pos.Id)
	affected = q.Update(map[string]any{
		"nama":      pos.Nama,
		"urutan":    pos.Urutan,
		"saldo":     pos.Saldo,
		"parent_id": parentId,
		"level":     pos.Level,
		"is_show":   isShowInt,
		"is_leaf":   isLeafInt,
		"status":    pos.Status,
	})
	return affected, nil
}
func (x *PosMysqlRepository) SoftDeleteById(ctx context.Context, id string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Where("id", "=", id)
	q.Where("status", "=", "aktif")
	affected = q.Update(map[string]any{
		"status": "trashed",
	})
	if affected == 0 {
		return affected, helpers_error.NewEntryNotFoundError(posEntityName, "id", id)
	}
	return affected, nil
}
func (x *PosMysqlRepository) DeleteById(ctx context.Context, id string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Where("id", "=", id)
	q.Where("status", "=", "trashed")
	affected = q.Delete()
	if affected == 0 {
		return affected, helpers_error.NewEntryNotFoundError(posEntityName, "id", id)
	}
	return affected, nil
}
func (x *PosMysqlRepository) GetTrashed(ctx context.Context) []*keep_entities.Pos {
	q := x.newQueryRequest(ctx, "trashed", keep_request.NewGetPos())
	return x.newEntitiesFromRows(q.Get())
}
func (x *PosMysqlRepository) FindTrashedById(ctx context.Context, id string) (*keep_entities.Pos, error) {
	q := x.newQueryRequest(ctx, "trashed", keep_request.NewGetPos())
	q.Where("id", "=", id)
	models := x.newEntitiesFromRows(q.Get())
	model := &keep_entities.Pos{}
	if len(models) == 0 {
		return model, helpers_error.NewEntryNotFoundError(posEntityName, "id", id)
	}
	return models[0], nil
}
func (x *PosMysqlRepository) RestoreTrashedById(ctx context.Context, id string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Where("id", "=", id)
	q.Where("status", "=", "trashed")
	affected = q.Update(map[string]any{
		"status": "aktif",
	})
	if affected == 0 {
		return affected, helpers_error.NewEntryNotFoundError(posEntityName, "id", id)
	}
	return affected, nil
}
func (x *PosMysqlRepository) UpdateSaldo(ctx context.Context, id string, saldo int) (affected int) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Where("id", "=", id)
	affected = q.Update(map[string]any{
		"saldo": saldo,
	})
	return affected
}

func (x *PosMysqlRepository) GetJumlahById(ctx context.Context, id string) (saldo int) {
	q := x.newQueryRequest(ctx, "aktif", keep_request.NewGetPos())
	q.Select("IFNULL(SUM(saldo), 0)")
	q.Where("parent_id", "=", id)
	rows, cleanup := q.Get()
	defer cleanup()

	if !rows.Next() {
		panic(fmt.Errorf("next failed"))
	}

	saldo = 0
	err := rows.Scan(&saldo)
	helpers_error.PanicIfError(err)

	return saldo
}
func (x *PosMysqlRepository) CountChildren(ctx context.Context, id string) (count int) {
	q := x.newQueryRequest(ctx, "aktif", keep_request.NewGetPos())
	q.Select("COUNT(0)")
	q.Where("parent_id", "=", id)
	rows, cleanup := q.Get()
	defer cleanup()

	if !rows.Next() {
		panic("failed to next")
	}
	err := rows.Scan(&count)
	helpers_error.PanicIfError(err)

	return count
}
func (x *PosMysqlRepository) UpdateLeaf(ctx context.Context, id string, leaf bool) (affected int, err error) {
	isLeafInt := 0
	if leaf {
		isLeafInt = 1
	}

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Where("id", "=", id)
	affected = q.Update(map[string]any{
		"is_leaf": isLeafInt,
	})
	return affected, nil
}

func (x *PosMysqlRepository) newQueryRequest(ctx context.Context, status string, request *keep_request.GetPos) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Select("id, nama, urutan, saldo, parent_id, level, is_show, is_leaf, status")

	if status != "" {
		q.Where("status", "=", status)
	}
	if request.IsLeafOnly {
		q.Where("is_leaf", "=", "1")
	}

	return q
}
func (x *PosMysqlRepository) newEntitiesFromRows(rows *sql.Rows, cleanup func()) []*keep_entities.Pos {
	defer cleanup()

	models := make([]*keep_entities.Pos, 0)

	if rows == nil {
		return models
	}

	for rows.Next() {
		model := &keep_entities.Pos{}
		models = append(models, model)

		isShowInt := 0
		isLeafInt := 0
		var parentId *string
		helpers_error.PanicIfError(rows.Scan(
			&model.Id,
			&model.Nama,
			&model.Urutan,
			&model.Saldo,
			&parentId,
			&model.Level,
			&isShowInt,
			&isLeafInt,
			&model.Status,
		))
		if parentId != nil {
			model.ParentId = *parentId
		}
		model.IsLeaf = isLeafInt == 1
		model.IsShow = isShowInt == 1
	}

	return models
}

func (x *PosMysqlRepository) Cleanup() {
	x.dbCleanup()
}
