package basic_entities

type Changelog struct {
	Id        string `json:"id,omitempty"`
	Version   string `json:"version,omitempty" validate:"required"`
	Timestamp int64  `json:"timestamp,omitempty" validate:"required,number"`
	Body      string `json:"body,omitempty" validate:"required"`
}

func (x *Changelog) Copy() *Changelog {
	cp := *x
	return &cp
}
