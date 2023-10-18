package keep_entities

type Kantong struct {
	Id             string `json:"id,omitempty"`
	Nama           string `json:"nama,omitempty"`
	Urutan         int    `json:"urutan,omitempty"`
	Saldo          int    `json:"saldo,omitempty"`
	SaldoMengendap int    `json:"saldoMengendap,omitempty"`
	PosId          string `json:"posId,omitempty"`
	IsShow         bool   `json:"isShow,omitempty"`
	Status         string `json:"status,omitempty" validate:"oneof=aktif trashed"`
}

func (x *Kantong) Copy() *Kantong {
	cp := *x
	return &cp
}
func (x *Kantong) CalculateSaldoAktif() int {
	return x.Saldo - x.SaldoMengendap
}

type KantongHistory struct {
	Id        string `json:"id,omitempty"`
	KantongId string `json:"kantongId,omitempty" validate:"required,number"`
	Jumlah    int    `json:"jumlah,omitempty"`
	Uraian    string `json:"uraian,omitempty"`
	Waktu     int64  `json:"waktu,omitempty"`
}

func (x *KantongHistory) Copy() *KantongHistory {
	cp := *x
	return &cp
}
