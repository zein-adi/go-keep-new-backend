package keep_entities

import "github.com/zein-adi/go-keep-new-backend/helpers"

type Barang struct {
	Nama         string          `json:"nama,omitempty"`
	Harga        int             `json:"harga,omitempty"`
	Diskon       int             `json:"diskon,omitempty"`
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
	Lokasi      string
	Harga       int
	Diskon      int
	SatuanHarga float64
	Keterangan  string
}

func (x *BarangDetail) Copy() *BarangDetail {
	cp := *x
	return &cp
}
