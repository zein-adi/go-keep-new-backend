package auth_handlers_restful

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"time"
)

func NewPermissionRestful(service auth_service_interfaces.IPermissionServices) *PermissionRestful {
	return &PermissionRestful{service: service}
}

type PermissionRestful struct {
	service auth_service_interfaces.IPermissionServices
}

func (x *PermissionRestful) Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	models := x.service.Get(ctx)
	count := len(models)

	h.SendMultiResponse(w, http.StatusOK, models, count)
}
