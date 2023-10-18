package helpers_mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"time"
)

func OpenMySqlConnection() (*sql.DB, func()) {
	viper.SetDefault("MYSQL_HOSTNAME", "127.0.0.1")
	viper.SetDefault("MYSQL_PORT", 3306)

	username := viper.GetString("MYSQL_USERNAME")
	password := viper.GetString("MYSQL_PASSWORD")
	hostname := viper.GetString("MYSQL_HOSTNAME")
	port := viper.GetInt("MYSQL_PORT")
	database := viper.GetString("MYSQL_DBNAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", username, password, hostname, port, database)
	db, err := sql.Open("mysql", dsn)
	helpers_error.PanicIfError(err)
	cleanup := func() {
		helpers_error.PanicIfError(db.Close())
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	err = db.Ping()
	helpers_error.PanicIfError(err)

	return db, cleanup
}
