package auth_handlers_local

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_requests"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_service_interfaces"
	"time"
)

func NewRoleLocal(service auth_service_interfaces.IRoleServices) *RoleLocal {
	return &RoleLocal{service: service}
}

type RoleLocal struct {
	service auth_service_interfaces.IRoleServices
}

func (x *RoleLocal) Get() []*auth_entities.Role {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	request := auth_requests.NewGetRequest()
	request.Take = 0
	return x.service.Get(ctx, request)
}
