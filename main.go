package main

import (
	"flag"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_env"
	"github.com/zein-adi/go-keep-new-backend/routes"
)

/*
 * APP_ENV: production (default), development, testing
 */
func main() {
	helpers_env.Init(0)
	cliHandler()
}

func cliHandler() {
	version := ""
	flag.StringVar(&version, "v", "", "version: force to x version")
	name := ""
	flag.StringVar(&name, "n", "", "name: create migration with name: x")
	action := ""
	flag.StringVar(&action, "a", "", "action: up down version force create")
	domain := ""
	flag.StringVar(&domain, "d", "", "domain: auth keep")
	method := ""
	flag.StringVar(&method, "m", "", "method: migrate")
	username := ""
	flag.StringVar(&username, "u", "", "username: new user seed")
	password := ""
	flag.StringVar(&password, "p", "", "password: new user password")
	flag.Parse()
	if method == "migrate" {
		RunMigration(action, domain, version, name)
	} else if method == "seed" {
		RunSeed(username, password)
	} else {
		routes.StartHttpServer()
	}
}
