package auth_repos_memory

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/helpers"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
)

var authEntityName = "auth"

func NewAuthMemoryRepository() *AuthRepository {
	return &AuthRepository{}
}

type AuthRepository struct {
	data []string
}

func (r *AuthRepository) FindBlacklistByToken(_ context.Context, refreshToken string) (entryNotFoundError error) {
	_, err := helpers.FindIndex(r.data, func(s string) bool {
		return s == refreshToken
	})
	if err != nil {
		return helpers_error.NewEntryNotFoundError(authEntityName, "token", refreshToken)
	}
	return nil
}

func (r *AuthRepository) InsertBlackList(_ context.Context, refreshToken string) error {
	r.data = append(r.data, refreshToken)
	return nil
}
