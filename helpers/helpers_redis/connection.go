package helpers_redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"strconv"
)

func OpenRedisConnection() (*redis.Client, func()) {
	// TODO: get from env
	username := ""
	password := ""
	hostname := "127.0.0.1"
	port := 6379
	database := 0

	db := redis.NewClient(&redis.Options{
		Addr:     hostname + ":" + strconv.Itoa(port),
		Username: username,
		Password: password,
		DB:       database,
	})
	cleanup := func() {
		helpers_error.PanicIfError(db.Close())
	}

	helpers_error.PanicIfError(db.Ping(context.Background()).Err())
	return db, cleanup
}
