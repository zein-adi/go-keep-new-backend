package main

import (
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
	"strconv"
)

func RunMigration() {
	domainList := map[string]string{
		"auth": "domains/auth/repos/auth_repos_mysql/migrations",
	}

	args := flag.Args()
	if len(args) < 3 {
		panic("usage: migrate [action: up|down|version|force] [domain] [options]")
	}
	action := flag.Arg(1)
	domain := flag.Arg(2)
	version := flag.Arg(3)

	// Validation
	v := validator.New()
	data := map[string]interface{}{
		"action":  action,
		"domain":  domain,
		"version": version,
	}
	rules := map[string]interface{}{
		"action":  "required,oneof=up down version force",
		"domain":  "required,oneof=auth",
		"version": "required_if=action force,min=0",
	}
	err := v.ValidateMap(data, rules)
	helpers_error.PanicIfError(err)

	// Initialization
	m := &Migration{
		domainList: domainList,
	}
	m.init(domain)

	// Execution
	switch action {
	case "up":
		m.Up()
	case "down":
		m.Down()
	case "version":
		version, dirty, _ := m.Version()
		dirtyText := "true"
		if !dirty {
			dirtyText = "false"
		}
		fmt.Printf("Version  : %d\nDirty    : %s", version, dirtyText)
	case "force":
		versionInt, err := strconv.Atoi(version)
		helpers_error.PanicIfError(err)
		m.Force(versionInt)
	}
}

type Migration struct {
	*migrate.Migrate
	domainList map[string]string
}

func (m *Migration) Up() {
	err := m.Migrate.Up()
	fmt.Print("Migrate Up  : ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Success")
	}
	version, _, _ := m.Version()
	fmt.Printf("Version     : %d", version)
}
func (m *Migration) Down() {
	err := m.Migrate.Steps(-1)
	fmt.Print("Migrate Down : ")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Success")
	}
	version, _, err := m.Version()
	fmt.Printf("Version      : %d %s", version, err)
}
func (m *Migration) Force(version int) {
	helpers_error.PanicIfError(m.Migrate.Force(version))
	fmt.Printf("Migrate Force Success.\nVersion : %d", version)
}
func (m *Migration) Version() (version uint, dirty bool, err error) {
	version, dirty, err = m.Migrate.Version()
	if err != nil {
		return 0, false, err
	}
	return version, dirty, nil

}
func (m *Migration) Close() {
	se, de := m.Migrate.Close()
	helpers_error.PanicIfError(se)
	helpers_error.PanicIfError(de)
}
func (m *Migration) init(domain string) {
	db, _ := helpers_mysql.OpenMySqlConnection()
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	helpers_error.PanicIfError(err)
	m.Migrate, err = migrate.NewWithDatabaseInstance("file://"+m.domainList[domain], "mysql", driver)
	helpers_error.PanicIfError(err)
}
