package auth_repos_redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_redis"
	"time"
)

var authEntityName = "auth"

func NewAuthRedisRepository() *AuthMysqlRepository {
	db, dbCleanup := helpers_redis.OpenRedisConnection()
	return &AuthMysqlRepository{
		db:        db,
		dbCleanup: dbCleanup,
	}
}

type AuthMysqlRepository struct {
	db        *redis.Client
	data      []string
	dbCleanup func()
}

func (r *AuthMysqlRepository) FindBlacklistByToken(ctx context.Context, refreshToken string) (entryNotFoundError error) {
	result, err := r.db.Get(ctx, refreshToken).Result()
	if err != nil {
		return helpers_error.NewEntryNotFoundError(authEntityName, "token", refreshToken)
	}
	println(result)
	return nil
}
func (r *AuthMysqlRepository) InsertBlackList(ctx context.Context, refreshToken string) error {
	err := r.db.Set(ctx, refreshToken, "blacklisted", time.Hour*24*7).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthMysqlRepository) Cleanup() {
	r.dbCleanup()
}
