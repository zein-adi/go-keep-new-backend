package keep_repos_mysql

import (
	"context"
	"database/sql"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"time"
)

var lokasiEntityName = "lokasi"
var lokasiTableName = "keep_lokasi"

func NewLokasiMySqlRepository() *LokasiMysqlRepository {
	db, dbCleanup := helpers_mysql.OpenMySqlConnection()
	return &LokasiMysqlRepository{
		db:        db,
		dbCleanup: dbCleanup,
	}
}

type LokasiMysqlRepository struct {
	db        *sql.DB
	dbCleanup func()
}

func (x *LokasiMysqlRepository) Get(ctx context.Context, search string) []*keep_entities.Lokasi {
	q := x.newQueryRequest(ctx, search)
	return x.newEntitiesFromRows(q.Get())
}
func (x *LokasiMysqlRepository) Insert(ctx context.Context, lokasi *keep_entities.Lokasi) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, lokasiTableName)

	timezone := +7 * time.Hour
	lastUpdateString := time.Unix(lokasi.LastUpdate, 0).Add(-timezone).Format(time.DateTime)
	_, err = q.Insert(map[string]any{
		"nama":       lokasi.Nama,
		"lastUpdate": lastUpdateString,
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}
func (x *LokasiMysqlRepository) Update(ctx context.Context, lokasi *keep_entities.Lokasi) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, lokasiTableName)
	q.Where("nama", "=", lokasi.Nama)

	timezone := +7 * time.Hour
	lastUpdateString := time.Unix(lokasi.LastUpdate, 0).Add(-timezone).Format(time.DateTime)
	affected = q.Update(map[string]any{
		"nama":       lokasi.Nama,
		"lastUpdate": lastUpdateString,
	})
	if affected == 0 {
		return 0, helpers_error.NewEntryNotFoundError(lokasiEntityName, "nama", lokasi.Nama)
	}
	return affected, nil
}
func (x *LokasiMysqlRepository) DeleteByNama(ctx context.Context, nama string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, lokasiTableName)
	q.Where("nama", "=", nama)
	affected = q.Delete()
	if affected == 0 {
		return 0, helpers_error.NewEntryNotFoundError(lokasiEntityName, "nama", nama)
	}
	return affected, nil
}

func (x *LokasiMysqlRepository) newQueryRequest(ctx context.Context, search string) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, lokasiTableName)
	q.Select("nama,lastUpdate")

	if search != "" {
		q.Where("nama", "LIKE", "%"+search+"%")
	}

	return q
}
func (x *LokasiMysqlRepository) newEntitiesFromRows(rows *sql.Rows, cleanup func()) []*keep_entities.Lokasi {
	defer cleanup()

	models := make([]*keep_entities.Lokasi, 0)

	if rows == nil {
		return models
	}

	for rows.Next() {
		model := &keep_entities.Lokasi{}
		models = append(models, model)

		lastUpdateString := ""
		helpers_error.PanicIfError(rows.Scan(
			&model.Nama,
			&lastUpdateString,
		))
		t, err := time.Parse(time.RFC3339, lastUpdateString)
		helpers_error.PanicIfError(err)
		model.LastUpdate = t.Unix()
	}

	return models
}

func (x *LokasiMysqlRepository) Cleanup() {
	x.dbCleanup()
}
