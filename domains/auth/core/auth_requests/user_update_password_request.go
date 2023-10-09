package auth_requests

type UserUpdatePasswordRequest struct {
	Id                   string `json:"id" validate:"required"`
	Password             string `json:"password" validate:"required,min=8,max=72,valid_password"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
}
