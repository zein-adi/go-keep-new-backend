package keep_repos_mysql

import (
	"context"
	"database/sql"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"strconv"
)

var kantongEntityName = "kantong"
var kantongTableName = "keep_kantong"

func NewKantongMysqlRepository() *KantongMysqlRepository {
	db, dbCleanup := helpers_mysql.OpenMySqlConnection()
	return &KantongMysqlRepository{
		db:        db,
		dbCleanup: dbCleanup,
	}
}

type KantongMysqlRepository struct {
	db        *sql.DB
	dbCleanup func()
}

func (x *KantongMysqlRepository) Get(ctx context.Context, request *helpers_requests.Get) []*keep_entities.Kantong {
	q := x.newQueryRequest(ctx, "aktif", request)
	q.OrderBy("urutan")

	if request.Take > 0 {
		q.Skip(request.Skip)
		q.Take(request.Take)
	}

	return x.newEntitiesFromRows(q.Get())
}

func (x *KantongMysqlRepository) FindById(ctx context.Context, id string) (*keep_entities.Kantong, error) {
	return x.findById(ctx, "aktif", id)
}

func (x *KantongMysqlRepository) Insert(ctx context.Context, kantong *keep_entities.Kantong) (*keep_entities.Kantong, error) {
	fields := map[string]any{
		"nama":            kantong.Nama,
		"urutan":          kantong.Urutan,
		"saldo":           kantong.Saldo,
		"saldo_mengendap": kantong.SaldoMengendap,
		"pos_id":          kantong.PosId,
		"is_show":         0,
		"status":          kantong.Status,
	}
	if kantong.IsShow {
		fields["is_show"] = 1
	}

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, kantongTableName)
	insertId, err := q.Insert(fields)
	if err != nil {
		return nil, err
	}
	model := kantong.Copy()
	model.Id = strconv.Itoa(insertId)
	return model, nil
}

func (x *KantongMysqlRepository) Update(ctx context.Context, kantong *keep_entities.Kantong) (affected int, err error) {
	fields := map[string]any{
		"nama":            kantong.Nama,
		"urutan":          kantong.Urutan,
		"saldo":           kantong.Saldo,
		"saldo_mengendap": kantong.SaldoMengendap,
		"pos_id":          kantong.PosId,
		"is_show":         0,
		"status":          kantong.Status,
	}
	if kantong.IsShow {
		fields["is_show"] = 1
	}

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, kantongTableName)
	q.Where("id", "=", kantong.Id)
	affected = q.Update(fields)
	return affected, nil
}
func (x *KantongMysqlRepository) UpdateSaldo(ctx context.Context, id string, saldo int) (affected int, err error) {
	q := x.newQueryRequest(ctx, "aktif", helpers_requests.NewGet())
	q.Where("id", "=", id)
	affected = q.Update(map[string]any{
		"saldo": saldo,
	})
	return affected, nil
}
func (x *KantongMysqlRepository) UpdateUrutan(ctx context.Context, id string, urutan int, posId string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, kantongTableName)
	q.Where("id", "=", id)
	affected = q.Update(map[string]any{
		"urutan": urutan,
		"pos_id": posId,
	})
	return affected, nil
}
func (x *KantongMysqlRepository) UpdateVisibility(ctx context.Context, id string, isShow bool) (affected int, err error) {
	isShowInt := 0
	if isShow {
		isShowInt = 1
	}

	q := helpers_mysql.NewQueryBuilder(ctx, x.db, kantongTableName)
	q.Where("id", "=", id)
	affected = q.Update(map[string]any{
		"is_show": isShowInt,
	})
	return affected, nil
}

func (x *KantongMysqlRepository) SoftDeleteById(ctx context.Context, id string) (affected int, err error) {
	q := x.newQueryRequest(ctx, "aktif", helpers_requests.NewGet())
	q.Where("id", "=", id)
	affected = q.Update(map[string]any{
		"status": "trashed",
	})
	return affected, nil
}
func (x *KantongMysqlRepository) GetTrashed(ctx context.Context, request *helpers_requests.Get) []*keep_entities.Kantong {
	q := x.newQueryRequest(ctx, "trashed", request)
	if request.Take > 0 {
		q.Skip(request.Skip)
		q.Take(request.Take)
	}
	return x.newEntitiesFromRows(q.Get())
}
func (x *KantongMysqlRepository) FindTrashedById(ctx context.Context, id string) (*keep_entities.Kantong, error) {
	return x.findById(ctx, "trashed", id)
}
func (x *KantongMysqlRepository) RestoreTrashedById(ctx context.Context, id string) (affected int, err error) {
	q := x.newQueryRequest(ctx, "trashed", helpers_requests.NewGet())
	q.Where("id", "=", id)
	affected = q.Update(map[string]any{
		"status": "aktif",
	})
	return affected, nil
}
func (x *KantongMysqlRepository) HardDeleteTrashedById(ctx context.Context, id string) (affected int, err error) {
	q := x.newQueryRequest(ctx, "trashed", helpers_requests.NewGet())
	q.Where("id", "=", id)
	return q.Delete(), nil
}

func (x *KantongMysqlRepository) newQueryRequest(ctx context.Context, status string, request *helpers_requests.Get) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, kantongTableName)
	q.Select("id,nama,urutan,saldo,saldo_mengendap,pos_id,is_show,status")

	if status != "" {
		q.Where("status", "=", status)
	}
	if request.Search != "" {
		q.Where("nama", "LIKE", "%"+request.Search+"%")
	}

	return q
}
func (x *KantongMysqlRepository) findById(ctx context.Context, status string, id string) (*keep_entities.Kantong, error) {
	q := x.newQueryRequest(ctx, status, helpers_requests.NewGet())
	q.Where("id", "=", id)
	models := x.newEntitiesFromRows(q.Get())
	if len(models) == 0 {
		return nil, helpers_error.NewEntryNotFoundError(kantongEntityName, "id", id)
	}
	return models[0], nil
}
func (x *KantongMysqlRepository) newEntitiesFromRows(rows *sql.Rows, cleanup func()) []*keep_entities.Kantong {
	defer cleanup()

	models := make([]*keep_entities.Kantong, 0)

	if rows == nil {
		return models
	}

	for rows.Next() {
		model := &keep_entities.Kantong{}
		models = append(models, model)

		isShowInt := 0
		helpers_error.PanicIfError(rows.Scan(
			&model.Id,
			&model.Nama,
			&model.Urutan,
			&model.Saldo,
			&model.SaldoMengendap,
			&model.PosId,
			&isShowInt,
			&model.Status,
		))
		model.IsShow = isShowInt == 1
	}
	return models
}
func (x *KantongMysqlRepository) Cleanup() {
	x.dbCleanup()
}
