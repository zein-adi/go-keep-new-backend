package helpers_env

import (
	"github.com/spf13/viper"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_directory"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"os"
	"strings"
)

var isInitated = false

func Init(relativeLevelToProjectRootDirectory int) {
	if isInitated == false {
		workingDirectory, _ := os.Getwd()
		workingDirectory = helpers_directory.Dir(workingDirectory, relativeLevelToProjectRootDirectory)

		appEnvironment := strings.ToLower(os.Getenv("APP_ENV"))
		envFile := workingDirectory + "/.env"
		if appEnvironment == "testing" {
			if helpers_directory.FileExists(workingDirectory + "/.env.testing") {
				envFile = workingDirectory + "/.env.testing"
			}
		} else if appEnvironment == "development" {
			if helpers_directory.FileExists(workingDirectory + "/.env.development") {
				envFile = workingDirectory + "/.env.development"
			}
		}

		viper.SetConfigFile(envFile)
		viper.SetConfigType("env")
		err := viper.ReadInConfig()
		helpers_error.PanicIfError(err)
		isInitated = true
	}
}
