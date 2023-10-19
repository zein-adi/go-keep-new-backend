package auth_handlers_restful_interfaces

import (
	"net/http"
)

type IRoleRestfulHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Insert(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	DeleteById(w http.ResponseWriter, r *http.Request)
}
