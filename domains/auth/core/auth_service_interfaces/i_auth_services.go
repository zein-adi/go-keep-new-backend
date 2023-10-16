package auth_service_interfaces

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/core/auth_responses"
)

type IAuthServices interface {
	Login(ctx context.Context, username, rawPassword string, rememberMe bool) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, refreshToken string) (accessToken, updatedRefreshToken string, err error)
	Logout(ctx context.Context, refreshToken string) error
	Profile(ctx context.Context, accessToken string) (*auth_responses.ProfileResponse, error)
}
