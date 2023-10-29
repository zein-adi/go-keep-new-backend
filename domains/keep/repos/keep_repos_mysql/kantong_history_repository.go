package keep_repos_mysql

import (
	"context"
	"database/sql"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"strconv"
	"time"
)

var kantongHistoryEntityName = "kantongHistory"
var kantongHistoryTable = "keep_kantong_history"

func NewKantongHistoryMysqlRepository() *KantongHistoryMysqlRepository {
	db, dbCleanup := helpers_mysql.OpenMySqlConnection()
	return &KantongHistoryMysqlRepository{
		db:        db,
		dbCleanup: dbCleanup,
	}
}

type KantongHistoryMysqlRepository struct {
	db        *sql.DB
	dbCleanup func()
}

func (x *KantongHistoryMysqlRepository) Get(ctx context.Context, request *helpers_requests.Get) []*keep_entities.KantongHistory {
	q := x.newQueryRequest(ctx)
	if request.Search != "" {
		q.Where("uraian", "LIKE", "%"+request.Search+"%")
	}
	if request.Take > 0 {
		q.Skip(request.Skip)
		q.Take(request.Take)
	}
	q.OrderBy("waktu desc")
	return x.newEntitiesFromRows(q.Get())
}

func (x *KantongHistoryMysqlRepository) FindById(ctx context.Context, id string) (*keep_entities.KantongHistory, error) {
	return x.findById(ctx, id)
}

func (x *KantongHistoryMysqlRepository) Insert(ctx context.Context, kantongHistory *keep_entities.KantongHistory) (*keep_entities.KantongHistory, error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, kantongHistoryTable)

	waktuString := time.Unix(kantongHistory.Waktu, 0).Format(time.DateTime)
	insertId, err := q.Insert(map[string]any{
		"kantong_id": kantongHistory.KantongId,
		"jumlah":     kantongHistory.Jumlah,
		"uraian":     kantongHistory.Uraian,
		"waktu":      waktuString,
	})
	if err != nil {
		return nil, err
	}
	model := kantongHistory.Copy()
	model.Id = strconv.Itoa(insertId)
	return model, nil
}

func (x *KantongHistoryMysqlRepository) Update(ctx context.Context, kantongHistory *keep_entities.KantongHistory) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, kantongHistoryTable)
	q.Where("id", "=", kantongHistory.Id)

	waktuString := time.Unix(kantongHistory.Waktu, 0).Format(time.DateTime)
	affected = q.Update(map[string]any{
		"kantong_id": kantongHistory.KantongId,
		"jumlah":     kantongHistory.Jumlah,
		"uraian":     kantongHistory.Uraian,
		"waktu":      waktuString,
	})
	return affected, nil
}

func (x *KantongHistoryMysqlRepository) DeleteById(ctx context.Context, id string) (affected int, err error) {
	q := x.newQueryRequest(ctx)
	q.Where("id", "=", id)
	return q.Delete(), nil
}

func (x *KantongHistoryMysqlRepository) newQueryRequest(ctx context.Context) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, kantongHistoryTable)
	q.Select("id,kantong_id,jumlah,uraian,waktu")
	return q
}
func (x *KantongHistoryMysqlRepository) newEntitiesFromRows(rows *sql.Rows, cleanup func()) []*keep_entities.KantongHistory {
	defer cleanup()

	models := make([]*keep_entities.KantongHistory, 0)

	if rows == nil {
		return models
	}

	for rows.Next() {
		model := &keep_entities.KantongHistory{}
		models = append(models, model)

		waktuString := ""
		helpers_error.PanicIfError(rows.Scan(
			&model.Id,
			&model.KantongId,
			&model.Jumlah,
			&model.Uraian,
			&waktuString,
		))
		t, err := time.Parse(time.RFC3339, waktuString)
		helpers_error.PanicIfError(err)
		model.Waktu = t.Unix()
	}
	return models
}
func (x *KantongHistoryMysqlRepository) findById(ctx context.Context, id string) (*keep_entities.KantongHistory, error) {
	q := x.newQueryRequest(ctx)
	q.Where("id", "=", id)
	models := x.newEntitiesFromRows(q.Get())
	if len(models) == 0 {
		return nil, helpers_error.NewEntryNotFoundError(kantongHistoryEntityName, "id", id)
	}
	return models[0], nil
}
func (x *KantongHistoryMysqlRepository) Cleanup() {
	x.dbCleanup()
}
