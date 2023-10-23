package auth_requests

type UserInputRequest struct {
	Username             string   `json:"username" validate:"required,alphanum,min=8,max=64"`
	Password             string   `json:"password" validate:"required,min=8,max=72,valid_password"`
	PasswordConfirmation string   `json:"passwordConfirmation" validate:"eqfield=Password"`
	Nama                 string   `json:"nama" validate:"required,min=3,max=128"`
	RoleIds              []string `json:"roleIds" validate:"required,min=1"`
}

type UserUpdateRequest struct {
	Id       string   `json:"id" validate:"required"`
	Username string   `json:"username" validate:"required,alphanum,min=8,max=64"`
	Nama     string   `json:"nama" validate:"required,min=3,max=128"`
	RoleIds  []string `json:"roleIds" validate:"required,min=1"`
}

type UserUpdatePasswordRequest struct {
	Id                   string `json:"id" validate:"required"`
	Password             string `json:"password" validate:"required,min=8,max=72,valid_password"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
}
