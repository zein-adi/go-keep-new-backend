package auth_service_interfaces

import (
	"context"
)

type IPermissionServices interface {
	Get(ctx context.Context) []string
}
