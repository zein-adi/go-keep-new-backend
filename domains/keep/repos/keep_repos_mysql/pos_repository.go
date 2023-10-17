package keep_repos_mysql

import (
	"context"
	"database/sql"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"strconv"
)

var posEntityName = "pos"

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

func (x *PosMysqlRepository) Get(ctx context.Context, request *keep_request.PosGetRequest) []*keep_entities.Pos {
	q := x.newQueryRequest(ctx, "aktif", request)
	return x.newEntitiesFromRows(q.Get())
}
func (x *PosMysqlRepository) FindById(ctx context.Context, id string) (*keep_entities.Pos, error) {
	q := x.newQueryRequest(ctx, "aktif", keep_request.NewPosGetRequest())
	q.Where("id", "=", id)
	models := x.newEntitiesFromRows(q.Get())
	model := &keep_entities.Pos{}
	if len(models) == 0 {
		return model, helpers_error.NewEntryNotFoundError(posEntityName, "id", id)
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

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posEntityName)
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
func (x *PosMysqlRepository) Update(ctx context.Context, pos *keep_entities.Pos) (*keep_entities.Pos, error) {
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

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posEntityName)
	q.Where("id", "=", pos.Id)
	affected := q.Update(map[string]any{
		"nama":      pos.Nama,
		"urutan":    pos.Urutan,
		"saldo":     pos.Saldo,
		"parent_id": parentId,
		"level":     pos.Level,
		"is_show":   isShowInt,
		"is_leaf":   isLeafInt,
		"status":    pos.Status,
	})
	if affected == 0 {
		return model, helpers_error.NewEntryNotFoundError(posEntityName, "id", pos.Id)
	}
	return pos.Copy(), nil
}
func (x *PosMysqlRepository) SoftDeleteById(ctx context.Context, id string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posEntityName)
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
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posEntityName)
	q.Where("id", "=", id)
	q.Where("status", "=", "trashed")
	affected = q.Delete()
	if affected == 0 {
		return affected, helpers_error.NewEntryNotFoundError(posEntityName, "id", id)
	}
	return affected, nil
}
func (x *PosMysqlRepository) GetTrashed(ctx context.Context) []*keep_entities.Pos {
	q := x.newQueryRequest(ctx, "trashed", keep_request.NewPosGetRequest())
	return x.newEntitiesFromRows(q.Get())
}
func (x *PosMysqlRepository) FindTrashedById(ctx context.Context, id string) (*keep_entities.Pos, error) {
	q := x.newQueryRequest(ctx, "trashed", keep_request.NewPosGetRequest())
	q.Where("id", "=", id)
	models := x.newEntitiesFromRows(q.Get())
	model := &keep_entities.Pos{}
	if len(models) == 0 {
		return model, helpers_error.NewEntryNotFoundError(posEntityName, "id", id)
	}
	return models[0], nil
}
func (x *PosMysqlRepository) RestoreTrashedById(ctx context.Context, id string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posEntityName)
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

func (x *PosMysqlRepository) newQueryRequest(ctx context.Context, status string, request *keep_request.PosGetRequest) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, posEntityName)
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
