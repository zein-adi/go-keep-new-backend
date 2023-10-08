package auth_entities

type Role struct {
	Id          string   `json:"id"`
	Nama        string   `json:"nama"`
	Deskripsi   string   `json:"deskripsi"`
	Level       int      `json:"level"`
	Permissions []string `json:"permissions"`
}

func (r *Role) Copy() *Role {
	val := *r
	return &val
}
