package auth_responses

type UserResponse struct {
	Id       string   `json:"id"`
	Username string   `json:"username"`
	Nama     string   `json:"nama"`
	RoleIds  []string `json:"roleIds"`
}
