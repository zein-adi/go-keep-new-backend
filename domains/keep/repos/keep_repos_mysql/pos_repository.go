package keep_repos_mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
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

func (x *PosMysqlRepository) Get(ctx context.Context) []*keep_entities.Pos {
	q := x.newQueryRequest(ctx, "aktif")
	q.OrderBy("urutan")
	return x.newEntitiesFromRows(q.Get())
}
func (x *PosMysqlRepository) GetJumlahById(ctx context.Context, id string) (saldo int) {
	q := x.newQueryRequest(ctx, "aktif")
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

func (x *PosMysqlRepository) GetChildrenById(ctx context.Context, id string) []*keep_entities.Pos {
	q := x.newQueryRequest(ctx, "aktif")
	q.Where("parent_id", "=", id)
	return x.newEntitiesFromRows(q.Get())
}

func (x *PosMysqlRepository) FindById(ctx context.Context, id string) (*keep_entities.Pos, error) {
	q := x.newQueryRequest(ctx, "aktif")
	q.Where("id", "=", id)
	models := x.newEntitiesFromRows(q.Get())
	if len(models) == 0 {
		return nil, helpers_error.NewEntryNotFoundError(posEntityName, "id", id)
	}
	return models[0], nil
}

func (x *PosMysqlRepository) Insert(ctx context.Context, pos *keep_entities.Pos) (*keep_entities.Pos, error) {
	model := &keep_entities.Pos{}

	fields := map[string]any{
		"nama":      pos.Nama,
		"urutan":    pos.Urutan,
		"saldo":     pos.Saldo,
		"parent_id": nil,
		"is_show":   0,
		"status":    pos.Status,
	}

	if pos.IsShow {
		fields["is_show"] = 1
	}
	if pos.ParentId != "" {
		fields["parent_id"] = pos.ParentId
	}

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	insertId, err := q.Insert(fields)
	if err != nil {
		return model, err
	}
	model = pos.Copy()
	model.Id = strconv.Itoa(insertId)
	return model, nil
}

func (x *PosMysqlRepository) Update(ctx context.Context, pos *keep_entities.Pos) (affected int, err error) {
	fields := map[string]any{
		"nama":      pos.Nama,
		"urutan":    pos.Urutan,
		"saldo":     pos.Saldo,
		"parent_id": nil,
		"is_show":   0,
		"status":    pos.Status,
	}

	if pos.IsShow {
		fields["is_show"] = 1
	}
	if pos.ParentId != "" {
		fields["parent_id"] = pos.ParentId
	}

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Where("id", "=", pos.Id)
	affected = q.Update(fields)
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
func (x *PosMysqlRepository) UpdateUrutan(ctx context.Context, id string, urutan int, parentId string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Where("id", "=", id)

	fields := map[string]any{
		"urutan":    urutan,
		"parent_id": nil,
	}
	if parentId != "" {
		fields["parent_id"] = parentId
	}

	affected = q.Update(fields)
	return affected, nil
}
func (x *PosMysqlRepository) UpdateVisibility(ctx context.Context, id string, isShow bool) (affected int, err error) {
	isShowInt := 0
	if isShow {
		isShowInt = 1
	}

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Where("id", "=", id)
	affected = q.Update(map[string]any{
		"is_show": isShowInt,
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
func (x *PosMysqlRepository) GetTrashed(ctx context.Context) []*keep_entities.Pos {
	q := x.newQueryRequest(ctx, "trashed")
	return x.newEntitiesFromRows(q.Get())
}
func (x *PosMysqlRepository) FindTrashedById(ctx context.Context, id string) (*keep_entities.Pos, error) {
	q := x.newQueryRequest(ctx, "trashed")
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
func (x *PosMysqlRepository) HardDeleteTrashedById(ctx context.Context, id string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Where("id", "=", id)
	q.Where("status", "=", "trashed")
	affected = q.Delete()
	if affected == 0 {
		return affected, helpers_error.NewEntryNotFoundError(posEntityName, "id", id)
	}
	return affected, nil
}

func (x *PosMysqlRepository) newQueryRequest(ctx context.Context, status string) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posTableName)
	q.Select("id, nama, urutan, saldo, parent_id, is_show, status")

	if status != "" {
		q.Where("status", "=", status)
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
		var parentId *string
		helpers_error.PanicIfError(rows.Scan(
			&model.Id,
			&model.Nama,
			&model.Urutan,
			&model.Saldo,
			&parentId,
			&isShowInt,
			&model.Status,
		))
		if parentId != nil {
			model.ParentId = *parentId
		}
		model.IsShow = isShowInt == 1
	}

	return models
}

func (x *PosMysqlRepository) Cleanup() {
	x.dbCleanup()
}
