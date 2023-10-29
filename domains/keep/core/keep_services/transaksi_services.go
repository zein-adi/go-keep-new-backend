package keep_services

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_events"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
	"time"
)

func NewTransaksiServices(repo keep_repo_interfaces.ITransaksiRepository,
	posRepo keep_repo_interfaces.IPosRepository,
	kantongRepo keep_repo_interfaces.IKantongRepository,
) *TransaksiServices {
	return &TransaksiServices{
		repo:        repo,
		posRepo:     posRepo,
		kantongRepo: kantongRepo,
	}
}

type TransaksiServices struct {
	repo        keep_repo_interfaces.ITransaksiRepository
	posRepo     keep_repo_interfaces.IPosRepository
	kantongRepo keep_repo_interfaces.IKantongRepository
}

func (x *TransaksiServices) Get(ctx context.Context, request *keep_request.GetTransaksi) []*keep_entities.Transaksi {
	return x.repo.Get(ctx, request)
}
func (x *TransaksiServices) Insert(ctx context.Context, transaksiRequest *keep_request.TransaksiInputUpdate) (*keep_entities.Transaksi, error) {
	err := x.validateBasic(transaksiRequest)
	if err != nil {
		return nil, err
	}

	source := &keep_entities.Transaksi{
		CreatedAt: time.Now().Unix(),
		Status:    "aktif",
	}
	transaksi, err := x.newEntityFromRequest(ctx, transaksiRequest, source)
	if err != nil {
		return nil, err
	}

	model, err := x.repo.Insert(ctx, transaksi)
	if err != nil {
		return nil, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(
		keep_events.TransaksiCreated,
		keep_events.TransaksiCreatedEventData{
			Time: time.Now(),
			Data: x.createEvent(transaksi),
		},
	)

	return model, nil
}
func (x *TransaksiServices) Update(ctx context.Context, transaksiRequest *keep_request.TransaksiInputUpdate) (affected int, err error) {
	source, err := x.repo.FindById(ctx, transaksiRequest.Id)
	if err != nil {
		return 0, err
	}

	err = x.validateBasic(transaksiRequest)
	if err != nil {
		return 0, err
	}

	transaksi, err := x.newEntityFromRequest(ctx, transaksiRequest, source)
	if err != nil {
		return 0, err
	}

	affected, err = x.repo.Update(ctx, transaksi)
	if err != nil {
		return 0, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(
		keep_events.TransaksiUpdated,
		keep_events.TransaksiUpdatedEventData{
			Time: time.Now(),
			Old:  x.createEvent(source),
			New:  x.createEvent(transaksi),
		})

	return affected, nil
}
func (x *TransaksiServices) DeleteById(ctx context.Context, id string) (affected int, err error) {
	m, err := x.repo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	affected, err = x.repo.SoftDeleteById(ctx, id)
	if err != nil {
		return 0, err
	}

	_ = helpers_events.GetDispatcher().Dispatch(
		keep_events.TransaksiSoftDeleted,
		keep_events.TransaksiSoftDeletedEventData{
			Time: time.Now(),
			Data: x.createEvent(m),
		})

	return affected, nil
}
func (x *TransaksiServices) GetTrashed(ctx context.Context) []*keep_entities.Transaksi {
	return x.repo.GetTrashed(ctx)
}
func (x *TransaksiServices) RestoreTrashedById(ctx context.Context, id string) (affected int, err error) {
	m, err := x.repo.FindTrashedById(ctx, id)
	if err != nil {
		return 0, err
	}
	affected, err = x.repo.RestoreTrashedById(ctx, id)
	if err != nil {
		return 0, nil
	}

	_ = helpers_events.GetDispatcher().Dispatch(
		keep_events.TransaksiRestored,
		keep_events.TransaksiRestoredEventData{
			Time: time.Now(),
			Data: x.createEvent(m),
		})

	return affected, nil
}
func (x *TransaksiServices) DeleteTrashedById(ctx context.Context, id string) (affected int, err error) {
	m, err := x.repo.FindTrashedById(ctx, id)
	if err != nil {
		return 0, err
	}
	affected, err = x.repo.HardDeleteTrashedById(ctx, id)
	if err != nil {
		return 0, nil
	}

	_ = helpers_events.GetDispatcher().Dispatch(
		keep_events.TransaksiHardDeleted,
		keep_events.TransaksiHardDeletedEventData{
			Time: time.Now(),
			Data: x.createEvent(m),
		})

	return affected, nil
}

func (x *TransaksiServices) newEntityFromRequest(ctx context.Context, transaksiRequest *keep_request.TransaksiInputUpdate, source *keep_entities.Transaksi) (*keep_entities.Transaksi, error) {
	posAsalNama := ""
	posTujuanNama := ""
	kantongAsalNama := ""
	kantongTujuanNama := ""
	if transaksiRequest.PosAsalId != "" {
		posAsal, err := x.posRepo.FindById(ctx, transaksiRequest.PosAsalId)
		if err != nil {
			return nil, err
		}
		count := x.posRepo.CountChildren(ctx, transaksiRequest.PosAsalId)
		if count > 0 {
			return nil, helpers_error.NewValidationErrors("posAsalId", "invalid", "has children")
		}
		posAsalNama = posAsal.Nama
	}
	if transaksiRequest.PosTujuanId != "" {
		posTujuan, err := x.posRepo.FindById(ctx, transaksiRequest.PosTujuanId)
		if err != nil {
			return nil, err
		}
		count := x.posRepo.CountChildren(ctx, transaksiRequest.PosTujuanId)
		if count > 0 {
			return nil, helpers_error.NewValidationErrors("posTujuanId", "invalid", "has children")
		}
		posTujuanNama = posTujuan.Nama
	}
	if transaksiRequest.KantongAsalId != "" {
		kantongAsal, err := x.kantongRepo.FindById(ctx, transaksiRequest.KantongAsalId)
		if err != nil {
			return nil, err
		}
		kantongAsalNama = kantongAsal.Nama
	}
	if transaksiRequest.KantongTujuanId != "" {
		kantongTujuan, err := x.kantongRepo.FindById(ctx, transaksiRequest.KantongTujuanId)
		if err != nil {
			return nil, err
		}
		kantongTujuanNama = kantongTujuan.Nama
	}

	transaksi := &keep_entities.Transaksi{
		Id:                transaksiRequest.Id,
		Waktu:             transaksiRequest.Waktu,
		Jenis:             transaksiRequest.Jenis,
		Jumlah:            transaksiRequest.Jumlah,
		PosAsalId:         transaksiRequest.PosAsalId,
		PosAsalNama:       posAsalNama,
		PosTujuanId:       transaksiRequest.PosTujuanId,
		PosTujuanNama:     posTujuanNama,
		KantongAsalId:     transaksiRequest.KantongAsalId,
		KantongAsalNama:   kantongAsalNama,
		KantongTujuanId:   transaksiRequest.KantongTujuanId,
		KantongTujuanNama: kantongTujuanNama,
		Uraian:            transaksiRequest.Uraian,
		Keterangan:        transaksiRequest.Keterangan,
		Lokasi:            transaksiRequest.Lokasi,
		UrlFoto:           transaksiRequest.UrlFoto,
		CreatedAt:         source.CreatedAt,
		UpdatedAt:         time.Now().Unix(),
		Details:           make([]*keep_entities.TransaksiDetail, 0),
		Status:            source.Status,
	}

	if len(transaksiRequest.Details) > 0 {
		transaksi.Jumlah = 0
		for _, detailRequest := range transaksiRequest.Details {
			total := int(float64(detailRequest.Harga)*detailRequest.Jumlah - float64(detailRequest.Diskon))
			satuanHarga := float64(total) / detailRequest.SatuanJumlah
			transaksi.Jumlah += total

			d := &keep_entities.TransaksiDetail{
				Uraian:       detailRequest.Uraian,
				Harga:        detailRequest.Harga,
				Jumlah:       detailRequest.Jumlah,
				Diskon:       detailRequest.Diskon,
				SatuanNama:   detailRequest.SatuanNama,
				SatuanJumlah: detailRequest.SatuanJumlah,
				SatuanHarga:  satuanHarga,
				Keterangan:   detailRequest.Keterangan,
			}
			transaksi.Details = append(transaksi.Details, d)
		}
	}

	return transaksi, nil
}
func (x *TransaksiServices) validateBasic(request *keep_request.TransaksiInputUpdate) error {
	v := validator.New()
	err := v.ValidateStruct(request)
	if err != nil {
		return err
	}
	return nil
}
func (x *TransaksiServices) createEvent(transaksi *keep_entities.Transaksi) keep_events.TransaksiEventData {
	detailsCopy := helpers.Map(transaksi.Details, func(d *keep_entities.TransaksiDetail) *keep_entities.TransaksiDetail {
		return d.Copy()
	})
	return keep_events.TransaksiEventData{
		Id:              transaksi.Id,
		PosAsalId:       transaksi.PosAsalId,
		PosTujuanId:     transaksi.PosTujuanId,
		KantongAsalId:   transaksi.KantongAsalId,
		KantongTujuanId: transaksi.KantongTujuanId,
		Jumlah:          transaksi.Jumlah,
		Lokasi:          transaksi.Lokasi,
		Details:         detailsCopy,
		Uraian:          transaksi.Uraian,
	}
}
