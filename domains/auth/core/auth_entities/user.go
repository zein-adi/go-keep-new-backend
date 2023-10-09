package auth_entities

type User struct {
	Id       string   `json:"id"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Nama     string   `json:"nama"`
	RoleIds  []string `json:"role_ids"`
}

func (r *User) Copy() *User {
	val := *r
	return &val
}
