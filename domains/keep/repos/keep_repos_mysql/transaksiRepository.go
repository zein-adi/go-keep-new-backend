package keep_repos_mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_datetime"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"strconv"
	"time"
)

var transaksiEntityName = "transaksi"
var transaksiTableName = "keep_transaksi"

func NewTransaksiMySqlRepository() *TransaksiMysqlRepository {
	db, dbCleanup := helpers_mysql.OpenMySqlConnection()
	return &TransaksiMysqlRepository{
		db:        db,
		dbCleanup: dbCleanup,
	}
}

type TransaksiMysqlRepository struct {
	db        *sql.DB
	dbCleanup func()
}

func (x *TransaksiMysqlRepository) Get(ctx context.Context, request *keep_request.GetTransaksi) []*keep_entities.Transaksi {
	q := x.newQueryRequest(ctx, "aktif", request)
	return x.newEntitiesFromRows(q.Get())
}
func (x *TransaksiMysqlRepository) FindById(ctx context.Context, id string) (*keep_entities.Transaksi, error) {
	q := x.newQueryRequest(ctx, "aktif", keep_request.NewGetTransaksi())
	q.Where("id", "=", id)
	models := x.newEntitiesFromRows(q.Get())
	if len(models) == 0 {
		return nil, helpers_error.NewEntryNotFoundError(transaksiEntityName, "id", id)
	}
	return models[0], nil
}
func (x *TransaksiMysqlRepository) Insert(ctx context.Context, transaksi *keep_entities.Transaksi) (*keep_entities.Transaksi, error) {
	q := x.newQueryRequest(ctx, "aktif", keep_request.NewGetTransaksi())
	q.Where("id", "=", transaksi.Id)
	lastId, err := q.Insert(x.newMapFromEntities(transaksi))
	if err != nil {
		return nil, err
	}
	model := transaksi.Copy()
	model.Id = strconv.Itoa(lastId)
	return model, nil
}
func (x *TransaksiMysqlRepository) Update(ctx context.Context, transaksi *keep_entities.Transaksi) (affected int, err error) {
	q := x.newQueryRequest(ctx, "aktif", keep_request.NewGetTransaksi())
	q.Where("id", "=", transaksi.Id)
	affected = q.Update(x.newMapFromEntities(transaksi))
	return affected, nil
}
func (x *TransaksiMysqlRepository) SoftDeleteById(ctx context.Context, id string) (affected int, err error) {
	q := x.newQueryRequest(ctx, "aktif", keep_request.NewGetTransaksi())
	q.Where("id", "=", id)
	affected = q.Update(map[string]any{
		"status": "trashed",
	})
	if affected == 0 {
		return affected, helpers_error.NewEntryNotFoundError(transaksiEntityName, "id", id)
	}
	return affected, nil
}
func (x *TransaksiMysqlRepository) GetTrashed(ctx context.Context) []*keep_entities.Transaksi {
	q := x.newQueryRequest(ctx, "trashed", keep_request.NewGetTransaksi())
	return x.newEntitiesFromRows(q.Get())
}
func (x *TransaksiMysqlRepository) FindTrashedById(ctx context.Context, id string) (*keep_entities.Transaksi, error) {
	q := x.newQueryRequest(ctx, "trashed", keep_request.NewGetTransaksi())
	q.Where("id", "=", id)
	models := x.newEntitiesFromRows(q.Get())
	if len(models) == 0 {
		return nil, helpers_error.NewEntryNotFoundError(transaksiEntityName, "id", id)
	}
	return models[0], nil
}
func (x *TransaksiMysqlRepository) RestoreTrashedById(ctx context.Context, id string) (affected int, err error) {
	q := x.newQueryRequest(ctx, "trashed", keep_request.NewGetTransaksi())
	q.Where("id", "=", id)
	affected = q.Update(map[string]any{
		"status": "aktif",
	})
	if affected == 0 {
		return affected, helpers_error.NewEntryNotFoundError(transaksiEntityName, "id", id)
	}
	return affected, nil
}
func (x *TransaksiMysqlRepository) HardDeleteTrashedById(ctx context.Context, id string) (affected int, err error) {
	q := x.newQueryRequest(ctx, "trashed", keep_request.NewGetTransaksi())
	q.Where("id", "=", id)
	affected = q.Delete()
	if affected == 0 {
		return affected, helpers_error.NewEntryNotFoundError(transaksiEntityName, "id", id)
	}
	return affected, nil
}
func (x *TransaksiMysqlRepository) GetJumlahByPosId(ctx context.Context, posId string) (saldo int) {
	pengeluaran := 0
	pemasukan := 0

	q := x.newQueryRequest(ctx, "aktif", keep_request.NewGetTransaksi())
	q.Select("IFNULL(SUM(jumlah), 0)")
	q.Where("pos_asal_id", "=", posId)
	rows, cleanup := q.Get()
	defer cleanup()
	if rows.Next() {
		helpers_error.PanicIfError(rows.Scan(&pengeluaran))
	}

	q = x.newQueryRequest(ctx, "aktif", keep_request.NewGetTransaksi())
	q.Select("IFNULL(SUM(jumlah), 0)")
	q.Where("pos_tujuan_id", "=", posId)
	rows, cleanup2 := q.Get()
	defer cleanup2()
	if rows.Next() {
		helpers_error.PanicIfError(rows.Scan(&pemasukan))
	}

	return pemasukan - pengeluaran
}

func (x *TransaksiMysqlRepository) newQueryRequest(ctx context.Context, status string, request *keep_request.GetTransaksi) *helpers_mysql.QueryBuilder {
	q := helpers_mysql.NewQueryBuilder(ctx, x.db, transaksiTableName)
	q.Select(
		"id",
		"waktu",
		"jenis",
		"jumlah",
		"pos_asal_id",
		"pos_asal_nama",
		"pos_tujuan_id",
		"pos_tujuan_nama",
		"kantong_asal_id",
		"kantong_asal_nama",
		"kantong_tujuan_id",
		"kantong_tujuan_nama",
		"uraian",
		"keterangan",
		"lokasi",
		"url_foto",
		"created_at",
		"updated_at",
		"details",
		"status",
	)

	if status != "" {
		q.Where("status", "=", status)
	}
	if request.PosId != "" {
		sub := q.WhereSub()
		sub.OrWhere("pos_asal_id", "=", request.PosId)
		sub.OrWhere("pos_tujuan_id", "=", request.PosId)
	}
	if request.KantongId != "" {
		sub := q.WhereSub()
		sub.OrWhere("kantong_asal_id", "=", request.KantongId)
		sub.OrWhere("kantong_tujuan_id", "=", request.KantongId)
	}
	if request.JenisTanggal != "" && request.Tanggal != 0 {
		q.Where("YEAR(waktu)", "=", time.Now().UTC().Format("2006"))
		if request.JenisTanggal == "bulan" {
			q.Where("MONTH(waktu)", "=", time.Now().UTC().Format("01"))
		}
		if request.JenisTanggal == "tanggal" {
			q.Where("DAY(waktu)", "=", time.Now().UTC().Format("02"))
		}
	}
	if request.WaktuAwal != 0 {
		q.Where("waktu", ">=", time.Now().UTC().Format(time.DateTime))
	}
	if request.Jenis != "" {
		q.Where("jenis", "=", request.Jenis)
	}

	return q
}
func (x *TransaksiMysqlRepository) newEntitiesFromRows(rows *sql.Rows, cleanup func()) []*keep_entities.Transaksi {
	defer cleanup()

	models := make([]*keep_entities.Transaksi, 0)

	if rows == nil {
		return models
	}

	for rows.Next() {
		model := &keep_entities.Transaksi{}
		models = append(models, model)

		waktuString := ""
		createdAtString := ""
		updatedAtString := ""
		detailsString := ""
		var posAsalId sql.NullInt64
		var posTujuanId sql.NullInt64
		var kantongAsalId sql.NullInt64
		var kantongTujuanId sql.NullInt64
		helpers_error.PanicIfError(rows.Scan(
			&model.Id,
			&waktuString,
			&model.Jenis,
			&model.Jumlah,
			&posAsalId,
			&model.PosAsalNama,
			&posTujuanId,
			&model.PosTujuanNama,
			&kantongAsalId,
			&model.KantongAsalNama,
			&kantongTujuanId,
			&model.KantongTujuanNama,
			&model.Uraian,
			&model.Keterangan,
			&model.Lokasi,
			&model.UrlFoto,
			&createdAtString,
			&updatedAtString,
			&detailsString,
			&model.Status,
		))
		model.PosAsalId = x.convertNullInt64ToString(posAsalId)
		model.PosTujuanId = x.convertNullInt64ToString(posTujuanId)
		model.KantongAsalId = x.convertNullInt64ToString(kantongAsalId)
		model.KantongTujuanId = x.convertNullInt64ToString(kantongTujuanId)
		model.Waktu = helpers_datetime.ParseStringGmt(waktuString)
		model.CreatedAt = helpers_datetime.ParseStringGmt(createdAtString)
		model.UpdatedAt = helpers_datetime.ParseStringGmt(updatedAtString)
		helpers_error.PanicIfError(json.Unmarshal([]byte(detailsString), &model.Details))
	}

	return models
}
func (x *TransaksiMysqlRepository) newMapFromEntities(model *keep_entities.Transaksi, fields ...string) map[string]any {

	waktuString := helpers_datetime.ParseTimestampGmt(model.Waktu)
	createdAtString := helpers_datetime.ParseTimestampGmt(model.CreatedAt)
	updatedAtString := helpers_datetime.ParseTimestampGmt(model.UpdatedAt)
	detailsString, err := json.Marshal(model.Details)
	helpers_error.PanicIfError(err)

	posAsalId := x.convertStringIdToNullInt64(model.PosAsalId)
	posTujuanId := x.convertStringIdToNullInt64(model.PosTujuanId)
	kantongAsalId := x.convertStringIdToNullInt64(model.KantongAsalId)
	kantongTujuanId := x.convertStringIdToNullInt64(model.KantongTujuanId)

	tmp := map[string]any{
		"waktu":               waktuString,
		"jenis":               model.Jenis,
		"jumlah":              model.Jumlah,
		"pos_asal_id":         posAsalId,
		"pos_asal_nama":       model.PosAsalNama,
		"pos_tujuan_id":       posTujuanId,
		"pos_tujuan_nama":     model.PosTujuanNama,
		"kantong_asal_id":     kantongAsalId,
		"kantong_asal_nama":   model.KantongAsalNama,
		"kantong_tujuan_id":   kantongTujuanId,
		"kantong_tujuan_nama": model.KantongTujuanNama,
		"uraian":              model.Uraian,
		"keterangan":          model.Keterangan,
		"lokasi":              model.Lokasi,
		"url_foto":            model.UrlFoto,
		"created_at":          createdAtString,
		"updated_at":          updatedAtString,
		"details":             string(detailsString),
		"status":              model.Status,
	}
	var result map[string]any
	if len(fields) == 0 {
		result = tmp
	} else {
		for _, field := range fields {
			result[field] = tmp[field]
		}
	}
	return result
}
func (x *TransaksiMysqlRepository) convertStringIdToNullInt64(id string) sql.NullInt64 {
	var res = sql.NullInt64{
		Int64: 0,
		Valid: false,
	}
	if id == "" {
		res.Valid = false
	} else {
		idInt, err2 := strconv.ParseInt(id, 10, 64)
		helpers_error.PanicIfError(err2)
		res.Int64 = idInt
		res.Valid = true
	}
	return res
}
func (x *TransaksiMysqlRepository) convertNullInt64ToString(id sql.NullInt64) string {
	if id.Valid {
		return strconv.FormatInt(id.Int64, 10)
	}
	return ""
}

func (x *TransaksiMysqlRepository) Cleanup() {
	x.dbCleanup()
}
