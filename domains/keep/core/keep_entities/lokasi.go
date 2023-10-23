package keep_entities

type Lokasi struct {
	Nama       string `json:"nama,omitempty"`
	LastUpdate int64  `json:"lastUpdate,omitempty"`
}

func (x *Lokasi) Copy() *Lokasi {
	cp := *x
	return &cp
}
