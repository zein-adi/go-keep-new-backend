package keep_request

type KantongInsert struct {
	Nama           string `json:"nama,omitempty" validate:"required"`
	Urutan         int    `json:"urutan,omitempty" validate:"min=0"`
	Saldo          int    `json:"saldo,omitempty" validate:"number,min=0"`
	SaldoMengendap int    `json:"saldoMengedap,omitempty" validate:"number,min=0"`
	PosId          string `json:"posId,omitempty" validate:"required,number"`
}

type KantongUpdate struct {
	Id             string `json:"id,omitempty" validate:"required,number"`
	Nama           string `json:"nama,omitempty" validate:"required"`
	Urutan         int    `json:"urutan,omitempty" validate:"min=0"`
	Saldo          int    `json:"saldo,omitempty" validate:"number,min=0"`
	SaldoMengendap int    `json:"saldoMengnedap,omitempty" validate:"number,min=0"`
	PosId          string `json:"posId,omitempty" validate:"required,number"`
	IsShow         bool   `json:"isShow,omitempty" validate:"boolean"`
}
