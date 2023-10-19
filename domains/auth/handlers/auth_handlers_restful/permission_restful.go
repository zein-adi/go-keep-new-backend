package auth_handlers_restful

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"time"
)

func NewPermissionRestfulHandler(service auth_service_interfaces.IPermissionServices) *PermissionRestfulHandler {
	return &PermissionRestfulHandler{service: service}
}

type PermissionRestfulHandler struct {
	service auth_service_interfaces.IPermissionServices
}

func (x *PermissionRestfulHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	accessToken, _ := middlewares.GetAuthorizationToken(r)
	accessClaim, _ := middlewares.GetJwtClaims(accessToken)
	models := x.service.Get(ctx, accessClaim.RoleIds)
	count := len(models)

	h.SendMultiResponse(w, http.StatusOK, models, count)
}
