package auth_handlers_local_interfaces

import "github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_entities"

type IRoleLocalHandler interface {
	Get() []*auth_entities.Role
}
