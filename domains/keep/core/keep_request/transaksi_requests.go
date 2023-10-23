package keep_request

type TransaksiInputUpdate struct {
	Id              string                        `json:"id,omitempty"`
	Waktu           int64                         `json:"waktu,omitempty" validate:"required"`
	Jenis           string                        `json:"jenis,omitempty" validate:"required,oneof=pemasukan pengeluaran mutasi"`
	Jumlah          int                           `json:"jumlah,omitempty" validate:"number,min=0"`
	PosAsalId       string                        `json:"posAsalId,omitempty" validate:"required_if=Jenis pengeluaran,required_if=Jenis mutasi,excluded_if=Jenis pemasukan"`
	PosTujuanId     string                        `json:"posTujuanId,omitempty" validate:"required_if=Jenis pemasukan,required_if=Jenis mutasi,excluded_if=Jenis pengeluaran"`
	KantongAsalId   string                        `json:"kantongAsalId,omitempty" validate:"excluded_if=Jenis pemasukan"`
	KantongTujuanId string                        `json:"kantongTujuanId,omitempty" validate:"excluded_if=Jenis pengeluaran"`
	Uraian          string                        `json:"uraian,omitempty" validate:"required"`
	Keterangan      string                        `json:"keterangan,omitempty"`
	UrlFoto         string                        `json:"urlFoto,omitempty" validate:"omitempty,url,excluded_if=Jenis mutasi"`
	Lokasi          string                        `json:"lokasi,omitempty" validate:"excluded_if=Jenis pemasukan,excluded_if=Jenis mutasi,required_with=Details"`
	Details         []*TransaksiInputUpdateDetail `json:"details,omitempty" validate:"excluded_if=Jenis pemasukan,excluded_if=Jenis mutasi,dive"`
}

type TransaksiInputUpdateDetail struct {
	Uraian       string  `json:"uraian,omitempty" validate:"required"`
	Harga        int     `json:"harga,omitempty" validate:"number,min=1"`
	Jumlah       float64 `json:"jumlah,omitempty" validate:"number,min=0.1"`
	Diskon       int     `json:"diskon,omitempty" validate:"number,min=0"`
	SatuanNama   string  `json:"satuan,omitempty" validate:"required,alpha"`
	SatuanJumlah float64 `json:"satuanJumlah,omitempty" validate:"number,min=0.1"`
	Keterangan   string  `json:"keterangan,omitempty"`
}

func NewGetTransaksi() *GetTransaksi {
	return &GetTransaksi{}
}

type GetTransaksi struct {
	PosId        string `json:"posId,omitempty"`
	KantongId    string `json:"kantongId,omitempty"`
	JenisTanggal string `json:"jenisTanggal,omitempty" validate:"oneof=tahun bulan tanggal"`
	Tanggal      int64  `json:"tanggal,omitempty"`
	WaktuAwal    int64  `json:"tanggalAwal,omitempty"`
	Jenis        string `json:"jenis,omitempty" validate:"required,oneof=pemasukan pengeluaran mutasi"`
}
