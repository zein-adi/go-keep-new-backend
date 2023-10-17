package tests

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_repo_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"testing"
	"time"
)

func TestTransaksi(t *testing.T) {
	x := NewTransaksiServicesTest()

	t.Run("GetSuccess", func(t *testing.T) {
		x.setUpAndPopulate()

	})
	t.Run("InsertPemasukanSuccess", func(t *testing.T) {
		// ToDo implement
	})
	t.Run("InsertPengeluaranSuccess", func(t *testing.T) {
		// ToDo implement
	})
	t.Run("InsertMutasiSuccess", func(t *testing.T) {
		// ToDo implement
	})
	t.Run("UpdatePemasukanSuccess", func(t *testing.T) {
		// ToDo implement
	})
	t.Run("UpdatePengeluaranSuccess", func(t *testing.T) {
		// ToDo implement
	})
	t.Run("UpdateMutasiSuccess", func(t *testing.T) {
		// ToDo implement
	})
	t.Run("DeleteSuccess", func(t *testing.T) {
		// ToDo implement
	})
	t.Run("RestoreSuccess", func(t *testing.T) {
		// ToDo implement
	})
	t.Run("ForceDeleteSuccess", func(t *testing.T) {
		// ToDo implement
	})
}

func NewTransaksiServicesTest() *TransaksiServicesTest {
	return &TransaksiServicesTest{}
}

type TransaksiServicesTest struct {
	repo     keep_repo_interfaces.ITransaksiRepository
	posRepo  keep_repo_interfaces.IPosRepository
	services keep_service_interfaces.ITransaksiServices
}

func (x *TransaksiServicesTest) setUp() {

}
func (x *TransaksiServicesTest) setUpMemoryRepository() {

}
func (x *TransaksiServicesTest) setUpMysqlRepository() {

}
func (x *TransaksiServicesTest) setUpAndPopulate() ([]*keep_entities.Transaksi, []*keep_entities.Pos) {
	posInput := []*keep_entities.Pos{
		{
			Nama:     "Pemasukan",
			Urutan:   0,
			Saldo:    0,
			ParentId: "",
			Level:    0,
			IsShow:   false,
		},
		{
			Nama:     "Main",
			Urutan:   0,
			Saldo:    0,
			ParentId: "",
			Level:    0,
			IsShow:   false,
		},
	}
	poses := make([]*keep_entities.Pos, 0)
	for _, pos := range posInput {
		m, _ := x.posRepo.Insert(context.Background(), pos)
		poses = append(poses, m)
	}

	transaksiInput := []*keep_entities.Transaksi{
		{
			Waktu:         time.Now().Unix(),
			Jenis:         "pemasukan",
			Jumlah:        100000,
			PosTujuanId:   "1",
			PosTujuanNama: "Main",
			Uraian:        "Gajian",
			CreatedAt:     time.Now().Unix(),
			UpdateAt:      time.Now().Unix(),
			Status:        "aktif",
		},
		{
			Waktu:       time.Now().Unix(),
			Jenis:       "pengeluaran",
			Jumlah:      10000,
			PosAsalId:   "1",
			PosAsalNama: "Main",
			Uraian:      "Cimory",
			CreatedAt:   time.Now().Unix(),
			UpdateAt:    time.Now().Unix(),
			Status:      "aktif",
		},
		{
			Waktu:       time.Now().Unix(),
			Jenis:       "pengeluaran",
			Jumlah:      10000,
			PosAsalId:   "1",
			PosAsalNama: "Main",
			Uraian:      "Cimory",
			CreatedAt:   time.Now().Unix(),
			UpdateAt:    time.Now().Unix(),
			Status:      "aktif",
			Details: []keep_entities.TransaksiDetail{
				{
					Uraian:       "Cimory Yoghurt Squeeze Peach 120ml",
					Harga:        10000,
					Jumlah:       1,
					Diskon:       0,
					SatuanJumlah: 120,
					Satuan:       "ml",
					SatuanHarga:  83.3333,
					Keterangan:   "",
				},
			},
		},
		{
			Waktu:         time.Now().Unix(),
			Jenis:         "mutasi",
			Jumlah:        10000,
			PosAsalId:     "2",
			PosAsalNama:   "Makan",
			PosTujuanId:   "1",
			PosTujuanNama: "Main",
			Uraian:        "Reimburse",
			CreatedAt:     time.Now().Unix(),
			UpdateAt:      time.Now().Unix(),
			Status:        "aktif",
		},
	}
	transaksis := make([]*keep_entities.Transaksi, 0)
	for _, v := range transaksiInput {
		m, _ := x.repo.Insert(context.Background(), v)
		transaksis = append(transaksis, m)
	}
	return transaksis, poses
}
