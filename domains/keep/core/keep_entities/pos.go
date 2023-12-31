package keep_entities

type Pos struct {
	Id       string `json:"id,omitempty"`
	Nama     string `json:"nama,omitempty"`
	Urutan   int    `json:"urutan,omitempty"`
	Saldo    int    `json:"saldo,omitempty"`
	ParentId string `json:"parentId,omitempty"`
	IsShow   bool   `json:"isShow"`
	Status   string `json:"status,omitempty" validate:"oneof=aktif trashed"`
}

func (p *Pos) Copy() *Pos {
	cp := *p
	return &cp
}
