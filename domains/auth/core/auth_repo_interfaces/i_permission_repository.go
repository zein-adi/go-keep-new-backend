package auth_repo_interfaces

import (
	"context"
)

type IPermissionRepository interface {
	Get(ctx context.Context) []string
}
