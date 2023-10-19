package auth_handlers_restful_interfaces

import (
	"net/http"
)

type IAuthRestfulHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
}
