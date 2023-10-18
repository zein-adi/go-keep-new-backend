package helpers_redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"strconv"
)

func OpenRedisConnection() (*redis.Client, func()) {
	viper.SetDefault("REDIS_HOSTNAME", "127.0.0.1")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DBNAME", 0)

	username := viper.GetString("REDIS_USERNAME")
	password := viper.GetString("REDIS_PASSWORD")
	hostname := viper.GetString("REDIS_HOSTNAME")
	port := viper.GetInt("REDIS_PORT")
	database := viper.GetInt("REDIS_DBNAME")

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
