package keep_repos_memory

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"strconv"
	"time"
)

var transaksiEntityName = "trasaksi"

func NewTransaksiMemoryRepository() *TransaksiMemoryRepository {
	return &TransaksiMemoryRepository{}
}

type TransaksiMemoryRepository struct {
	Data []*keep_entities.Transaksi
}

func (x *TransaksiMemoryRepository) Get(_ context.Context, request *keep_request.GetTransaksi) []*keep_entities.Transaksi {
	models := x.newQueryRequest("aktif", request)
	return helpers.Map(models, func(d *keep_entities.Transaksi) *keep_entities.Transaksi {
		return d.Copy()
	})
}
func (x *TransaksiMemoryRepository) FindById(_ context.Context, id string) (*keep_entities.Transaksi, error) {
	index, err := x.findIndexById(id, "aktif")
	if err != nil {
		return nil, err
	}
	return x.Data[index].Copy(), err
}
func (x *TransaksiMemoryRepository) Insert(_ context.Context, transaksi *keep_entities.Transaksi) (*keep_entities.Transaksi, error) {
	lastId := helpers.Reduce(x.Data, 0, func(accumulator int, pos *keep_entities.Transaksi) int {
		datumId, _ := strconv.Atoi(pos.Id)
		return max(accumulator, datumId)
	})

	model := transaksi.Copy()
	model.Id = strconv.Itoa(lastId + 1)
	x.Data = append(x.Data, model)
	return model, nil
}
func (x *TransaksiMemoryRepository) Update(_ context.Context, transaksi *keep_entities.Transaksi) (affected int, err error) {
	index, err := x.findIndexById(transaksi.Id, "aktif")
	if err != nil {
		return 0, err
	}

	model := transaksi.Copy()
	x.Data[index] = model
	return 1, nil
}
func (x *TransaksiMemoryRepository) SoftDeleteById(_ context.Context, id string) (affected int, err error) {
	index, err := x.findIndexById(id, "aktif")
	if err != nil {
		return 0, err
	}

	x.Data[index].Status = "trashed"
	return 1, nil
}
func (x *TransaksiMemoryRepository) GetTrashed(_ context.Context) []*keep_entities.Transaksi {
	models := x.newQueryRequest("trashed", keep_request.NewGetTransaksi())
	return helpers.Map(models, func(d *keep_entities.Transaksi) *keep_entities.Transaksi {
		return d.Copy()
	})
}
func (x *TransaksiMemoryRepository) FindTrashedById(_ context.Context, id string) (*keep_entities.Transaksi, error) {
	index, err := x.findIndexById(id, "trashed")
	if err != nil {
		return nil, err
	}
	return x.Data[index].Copy(), err
}
func (x *TransaksiMemoryRepository) RestoreTrashedById(_ context.Context, id string) (affected int, err error) {
	index, err := x.findIndexById(id, "trashed")
	if err != nil {
		return 0, err
	}

	x.Data[index].Status = "aktif"
	return 1, nil
}
func (x *TransaksiMemoryRepository) HardDeleteTrashedById(_ context.Context, id string) (affected int, err error) {
	index, err := x.findIndexById(id, "trashed")
	if err != nil {
		return 0, err
	}

	x.Data = append(x.Data[0:index], x.Data[index+1:]...)
	return 1, nil
}

func (x *TransaksiMemoryRepository) GetJumlahByPosId(_ context.Context, posId string) (saldo int) {
	for _, t := range x.Data {
		if t.Status != "aktif" {
			continue
		}
		if t.PosTujuanId == posId {
			saldo += t.Jumlah
		}
		if t.PosAsalId == posId {
			saldo -= t.Jumlah
		}
	}
	return saldo
}

func (x *TransaksiMemoryRepository) newQueryRequest(status string, request *keep_request.GetTransaksi) []*keep_entities.Transaksi {
	return helpers.Filter(x.Data, func(pos *keep_entities.Transaksi) bool {
		res := true
		if status != "" {
			res = res && pos.Status == status
		}
		if request.PosId != "" {
			res = res && (pos.PosAsalId == request.PosId || pos.PosTujuanId == request.PosId)
		}
		if request.KantongId != "" {
			res = res && (pos.KantongAsalId == request.KantongId || pos.KantongTujuanId == request.KantongId)
		}
		if request.JenisTanggal != "" && request.Tanggal != 0 {
			waktuTime := time.Unix(pos.Waktu, 0)
			requestTanggal := time.Unix(request.Tanggal, 0)
			format := time.DateTime
			switch request.JenisTanggal {
			case "tahun":
				format = "2006"
			case "bulan":
				format = "2006-01"
			case "tanggal":
				format = "2006-01-02"
			}
			res = res && waktuTime.Format(format) == requestTanggal.Format(format)
		}
		return res
	})
}
func (x *TransaksiMemoryRepository) findIndexById(id string, status string) (index int, err error) {
	index, err = helpers.FindIndex(x.Data, func(pos *keep_entities.Transaksi) bool {
		res := pos.Id == id
		if status != "" {
			res = res && pos.Status == status
		}
		return res
	})
	if err != nil {
		return -1, helpers_error.NewEntryNotFoundError(transaksiEntityName, "id", "id")
	}
	return index, err
}
