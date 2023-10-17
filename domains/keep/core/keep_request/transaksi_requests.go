package keep_request

type Transaksi struct {
	Id          string            `json:"id,omitempty"`
	Jenis       string            `json:"jenis,omitempty" validate:"required,oneof=pemasukan pengeluaran mutasi"`
	Jumlah      int               `json:"jumlah,omitempty" validate:"required"`
	PosAsalId   string            `json:"posAsalId,omitempty" validate:"required_if=Jenis pengeluaran,required_if=Jenis mutasi"`
	PosTujuanId string            `json:"posTujuanId,omitempty" validate:"required_if=Jenis pemasukan,required_if=Jenis mutasi"`
	Uraian      string            `json:"uraian,omitempty" validate:"required"`
	Keterangan  string            `json:"keterangan,omitempty" validate:""`
	UrlFoto     string            `json:"urlFoto,omitempty" validate:"excluded_if=Jenis mutasi"`
	Lokasi      string            `json:"lokasi,omitempty" validate:"excluded_if=Jenis pemasukan,excluded_if=Jenis mutasi,required_with=Details"`
	Details     []TransaksiDetail `json:"details,omitempty" validate:"excluded_if=Jenis pemasukan,excluded_if=Jenis mutasi"`
}

type TransaksiDetail struct {
	Uraian       string  `json:"uraian,omitempty" validate:"required"`
	Harga        int     `json:"harga,omitempty" validate:"required,number,min:1"`
	Jumlah       float32 `json:"jumlah,omitempty" validate:"required,number,min:0.1"`
	Diskon       int     `json:"diskon,omitempty" validate:"required,number,min:0"`
	Satuan       string  `json:"satuan,omitempty" validate:"required,alpha"`
	SatuanJumlah float32 `json:"satuanJumlah,omitempty" validate:"required,number,min:0.1"`
	Keterangan   string  `json:"keterangan,omitempty" validate:""`
}
