package auth_entities

type Role struct {
	Id          string   `json:"id"`
	Nama        string   `json:"nama" validate:"required"`
	Deskripsi   string   `json:"deskripsi"`
	Level       int      `json:"level" validate:"min=1,max=65535"`
	Permissions []string `json:"permissions"`
}

func (r *Role) Copy() *Role {
	val := *r
	return &val
}
