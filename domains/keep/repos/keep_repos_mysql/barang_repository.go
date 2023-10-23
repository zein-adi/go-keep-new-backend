package keep_repos_mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"strings"
	"time"
)

var barangEntityName = "barang"
var barangTableName = "keep_barang"

func NewBarangMySqlRepository() *BarangMysqlRepository {
	db, dbCleanup := helpers_mysql.OpenMySqlConnection()
	return &BarangMysqlRepository{
		db:        db,
		dbCleanup: dbCleanup,
	}
}

type BarangMysqlRepository struct {
	db        *sql.DB
	dbCleanup func()
}

func (x *BarangMysqlRepository) Get(ctx context.Context, search string, lokasi string) []*keep_entities.Barang {
	lokasi = strings.ToLower(lokasi)
	q := x.newQueryRequest(ctx, search)
	models := x.newEntitiesFromRows(q.Get())
	return helpers.Map(models, func(model *keep_entities.Barang) *keep_entities.Barang {
		model.Details = helpers.Filter(model.Details, func(d *keep_entities.BarangDetail) bool {
			return strings.Contains(strings.ToLower(d.Lokasi), lokasi)
		})
		return model
	})
}
func (x *BarangMysqlRepository) Insert(ctx context.Context, barang *keep_entities.Barang) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, barangTableName)

	timezone := +7 * time.Hour
	lastUpdateString := time.Unix(barang.LastUpdate, 0).Add(-timezone).Format(time.DateTime)
	detailBytes, err := json.Marshal(barang.Details)
	helpers_error.PanicIfError(err)
	detailString := string(detailBytes)
	_, err = q.Insert(map[string]any{
		"nama":         barang.Nama,
		"harga":        barang.Harga,
		"diskon":       barang.Diskon,
		"satuanNama":   barang.SatuanNama,
		"satuanJumlah": barang.SatuanJumlah,
		"satuanHarga":  barang.SatuanHarga,
		"keterangan":   barang.Keterangan,
		"lastUpdate":   lastUpdateString,
		"detail":       detailString,
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}
func (x *BarangMysqlRepository) Update(ctx context.Context, barang *keep_entities.Barang) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, barangTableName)
	q.Where("nama", "=", barang.Nama)

	timezone := +7 * time.Hour
	lastUpdateString := time.Unix(barang.LastUpdate, 0).Add(-timezone).Format(time.DateTime)
	detailBytes, err := json.Marshal(barang.Details)
	helpers_error.PanicIfError(err)
	detailString := string(detailBytes)
	affected = q.Update(map[string]any{
		"nama":         barang.Nama,
		"harga":        barang.Harga,
		"diskon":       barang.Diskon,
		"satuanNama":   barang.SatuanNama,
		"satuanJumlah": barang.SatuanJumlah,
		"satuanHarga":  barang.SatuanHarga,
		"keterangan":   barang.Keterangan,
		"lastUpdate":   lastUpdateString,
		"detail":       detailString,
	})
	if affected == 0 {
		return 0, helpers_error.NewEntryNotFoundError(barangEntityName, "nama", barang.Nama)
	}
	return affected, nil
}
func (x *BarangMysqlRepository) DeleteByNama(ctx context.Context, nama string) (affected int, err error) {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, barangTableName)
	q.Where("nama", "=", nama)
	affected = q.Delete()
	if affected == 0 {
		return 0, helpers_error.NewEntryNotFoundError(barangEntityName, "nama", nama)
	}
	return affected, nil
}

func (x *BarangMysqlRepository) newQueryRequest(ctx context.Context, search string) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, barangTableName)
	q.Select("nama,harga,diskon,satuanNama,satuanJumlah,satuanHarga,keterangan,lastUpdate,detail")

	if search != "" {
		q.Where("nama", "LIKE", "%"+search+"%")
	}

	return q
}
func (x *BarangMysqlRepository) newEntitiesFromRows(rows *sql.Rows, cleanup func()) []*keep_entities.Barang {
	defer cleanup()

	models := make([]*keep_entities.Barang, 0)

	if rows == nil {
		return models
	}

	for rows.Next() {
		model := &keep_entities.Barang{}
		models = append(models, model)

		detailString := ""
		lastUpdateString := ""
		helpers_error.PanicIfError(rows.Scan(
			&model.Nama,
			&model.Harga,
			&model.Diskon,
			&model.SatuanNama,
			&model.SatuanJumlah,
			&model.SatuanHarga,
			&model.Keterangan,
			&lastUpdateString,
			&detailString,
		))
		helpers_error.PanicIfError(json.Unmarshal([]byte(detailString), &model.Details))
		t, err := time.Parse(time.RFC3339, lastUpdateString)
		helpers_error.PanicIfError(err)
		model.LastUpdate = t.Unix()
	}

	return models
}

func (x *BarangMysqlRepository) Cleanup() {
	x.dbCleanup()
}