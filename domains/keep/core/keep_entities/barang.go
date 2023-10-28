package keep_entities

import "github.com/zein-adi/go-keep-new-backend/helpers"

type Barang struct {
	Nama         string          `json:"nama,omitempty"`
	Harga        float64         `json:"harga,omitempty"`
	Diskon       float64         `json:"diskon,omitempty"`
	SatuanNama   string          `json:"satuanNama,omitempty"`
	SatuanJumlah float64         `json:"satuanJumlah,omitempty"`
	SatuanHarga  float64         `json:"satuanHarga,omitempty"`
	Keterangan   string          `json:"keterangan,omitempty"`
	LastUpdate   int64           `json:"lastUpdate,omitempty"`
	Details      []*BarangDetail `json:"details,omitempty"`
}

func (x *Barang) Copy() *Barang {
	cp := *x
	x.Details = helpers.Map(x.Details, func(v *BarangDetail) *BarangDetail {
		return v.Copy()
	})
	return &cp
}

type BarangDetail struct {
	Lokasi      string  `json:"lokasi,omitempty"`
	Harga       float64 `json:"harga,omitempty"`
	Diskon      float64 `json:"diskon,omitempty"`
	SatuanHarga float64 `json:"satuanHarga,omitempty"`
	Keterangan  string  `json:"keterangan,omitempty"`
}

func (x *BarangDetail) Copy() *BarangDetail {
	cp := *x
	return &cp
}
