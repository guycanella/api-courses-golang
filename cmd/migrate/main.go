package main

import (
	"flag"
	"log"

	"github.com/guycanella/api-courses-golang/internal/domain"
	mysqlrepo "github.com/guycanella/api-courses-golang/internal/repository/mysql"
)

var (
	flagTest = flag.Bool("test", false, "run migrate in test environment")
)

func main() {
	flag.Parse()

	isTest := *flagTest
	env := "development"

	if isTest {
		env = "test"
	}

	db, err := mysqlrepo.OpenDatabase(&mysqlrepo.Env{E: &env})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Course{},
		&domain.Enrollment{},
	); err != nil {
		log.Fatal(err)
	}

	if isTest {
		log.Println("✅ AutoMigrate test concluded.")
		return
	}

	log.Println("✅ AutoMigrate concluded.")
}
