package auth_handlers_local

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_requests"
	"time"
)

func NewRoleLocalHandler(service auth_service_interfaces.IRoleServices) *RoleLocalHandler {
	return &RoleLocalHandler{service: service}
}

type RoleLocalHandler struct {
	service auth_service_interfaces.IRoleServices
}

func (x *RoleLocalHandler) Get() []*auth_entities.Role {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	request := helpers_requests.NewGet()
	request.Take = 0
	return x.service.Get(ctx, request)
}
