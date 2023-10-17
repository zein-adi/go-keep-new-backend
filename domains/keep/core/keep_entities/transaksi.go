package keep_entities

type Transaksi struct {
	Id            string            `json:"id,omitempty"`
	Waktu         int64             `json:"waktu,omitempty"`
	Jenis         string            `json:"jenis,omitempty"`
	Jumlah        int               `json:"jumlah,omitempty"`
	PosAsalId     string            `json:"posAsalId,omitempty"`
	PosAsalNama   string            `json:"posAsalNama,omitempty"`
	PosTujuanId   string            `json:"posTujuanId,omitempty"`
	PosTujuanNama string            `json:"posTujuanNama,omitempty"`
	Uraian        string            `json:"uraian,omitempty"`
	Keterangan    string            `json:"keterangan,omitempty"`
	Lokasi        string            `json:"lokasi,omitempty"`
	UrlFoto       string            `json:"urlFoto,omitempty"`
	CreatedAt     int64             `json:"createdAt,omitempty"`
	UpdateAt      int64             `json:"updateAt,omitempty"`
	Details       []TransaksiDetail `json:"details,omitempty"`
	Status        string            `json:"status,omitempty" validate:"oneof=aktif trashed"`
}

type TransaksiDetail struct {
	Uraian       string  `json:"uraian,omitempty"`
	Harga        int     `json:"harga,omitempty"`
	Jumlah       float32 `json:"jumlah,omitempty"`
	Diskon       int     `json:"diskon,omitempty"`
	Satuan       string  `json:"satuan,omitempty"`
	SatuanJumlah float32 `json:"satuanJumlah,omitempty"`
	SatuanHarga  float32 `json:"satuanHarga,omitempty"`
	Keterangan   string  `json:"keterangan,omitempty"`
}
