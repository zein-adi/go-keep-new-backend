package keep_entities

import "github.com/zein-adi/go-keep-new-backend/helpers"

type Transaksi struct {
	Id                string             `json:"id,omitempty"`
	Waktu             int64              `json:"waktu,omitempty"`
	Jenis             string             `json:"jenis,omitempty"`
	Jumlah            int                `json:"jumlah,omitempty"`
	PosAsalId         string             `json:"posAsalId,omitempty"`
	PosAsalNama       string             `json:"posAsalNama,omitempty"`
	PosTujuanId       string             `json:"posTujuanId,omitempty"`
	PosTujuanNama     string             `json:"posTujuanNama,omitempty"`
	KantongAsalId     string             `json:"kantongAsalId,omitempty"`
	KantongAsalNama   string             `json:"kantongAsalNama,omitempty"`
	KantongTujuanId   string             `json:"kantongTujuanId,omitempty"`
	KantongTujuanNama string             `json:"kantongTujuanNama,omitempty"`
	Uraian            string             `json:"uraian,omitempty"`
	Keterangan        string             `json:"keterangan,omitempty"`
	Lokasi            string             `json:"lokasi,omitempty"`
	UrlFoto           string             `json:"urlFoto,omitempty"`
	CreatedAt         int64              `json:"createdAt,omitempty"`
	UpdatedAt         int64              `json:"updatedAt,omitempty"`
	Details           []*TransaksiDetail `json:"details,omitempty"`
	Status            string             `json:"status,omitempty" validate:"oneof=aktif trashed"`
}

func (p *Transaksi) Copy() *Transaksi {
	cp := *p
	cp.Details = make([]*TransaksiDetail, 0)
	cp.Details = helpers.Map(p.Details, func(v *TransaksiDetail) *TransaksiDetail {
		return v.Copy()
	})
	return &cp
}

type TransaksiDetail struct {
	Uraian       string  `json:"uraian,omitempty"`
	Harga        int     `json:"harga,omitempty"`
	Jumlah       float64 `json:"jumlah,omitempty"`
	Diskon       int     `json:"diskon,omitempty"`
	SatuanNama   string  `json:"satuanNama,omitempty"`
	SatuanJumlah float64 `json:"satuanJumlah,omitempty"`
	SatuanHarga  float64 `json:"satuanHarga,omitempty"`
	Keterangan   string  `json:"keterangan,omitempty"`
}

func (p *TransaksiDetail) Copy() *TransaksiDetail {
	cp := *p
	return &cp
}
