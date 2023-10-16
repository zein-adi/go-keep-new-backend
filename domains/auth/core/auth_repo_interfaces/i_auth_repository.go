package auth_repo_interfaces

import (
	"context"
)

type IAuthRepository interface {
	FindBlacklistByToken(ctx context.Context, refreshToken string) (entryNotFoundError error)
	InsertBlackList(ctx context.Context, refreshToken string) error
}
