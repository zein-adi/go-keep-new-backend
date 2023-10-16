package auth_requests

type UserUpdateRequest struct {
	Id       string   `json:"id" validate:"required"`
	Username string   `json:"username" validate:"required,alphanum,min=8,max=64"`
	Nama     string   `json:"nama" validate:"required,min=3,max=128"`
	RoleIds  []string `json:"roleIds" validate:"required,min=1"`
}
