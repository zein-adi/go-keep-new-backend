package commands

import (
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_mysql"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
	"os/exec"
	"strconv"
)

func RunMigration(action, domain, version, name string) {
	// Validation
	v := validator.New()
	data := map[string]interface{}{
		"action":  action,
		"domain":  domain,
		"version": version,
		"name":    name,
	}
	rules := map[string]interface{}{
		"action":  "required,oneof=up down version force create",
		"version": "required_if=action force,min=0",
		"domain":  "required_if=action create",
		"name":    "required_if=action create",
	}
	err := v.ValidateMap(data, rules)
	helpers_error.PanicIfError(err)

	// Initialization
	m := &Migration{}
	m.init()

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
	case "create":
		fmt.Printf("Creating Migration: %s %s \n", domain, name)
		resultString, err := m.Create(domain, name)
		if err != nil {
			fmt.Print("error: ")
		}
		fmt.Println(resultString)
	}
}

type Migration struct {
	*migrate.Migrate
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
	fmt.Printf("Version      : %d", version)
	if err != nil {
		fmt.Printf(" %s", err)
	}
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
func (m *Migration) init() {
	db, _ := helpers_mysql.OpenMySqlConnection()
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	helpers_error.PanicIfError(err)
	m.Migrate, err = migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	helpers_error.PanicIfError(err)
}
func (m *Migration) Create(domain, name string) (string, error) {
	cmd := exec.Command("migrate", "create", "-ext", "sql", "-dir", "migrations", fmt.Sprintf("%s__%s", domain, name))
	res, err := cmd.CombinedOutput()
	return string(res), err
}
