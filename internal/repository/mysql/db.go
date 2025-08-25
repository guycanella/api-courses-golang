package mysql

import (
	"fmt"
	"os"

	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func dsnFromEnv(mode string) string {
	host := getenv("DB_HOST", "localhost")

	var name, port string
	if mode == "test" {
		name = getenv("DB_NAME_TEST", "go_api_db_test")
		port = getenv("DB_PORT_TEST", "3307")
		host = getenv("DB_HOST_TEST", host)
	} else {
		name = getenv("DB_NAME", "go_api_db")
		port = getenv("DB_PORT", "3306")
	}

	user := getenv("DB_USER", "api_user")
	pass := getenv("DB_PASSWORD", "mysql")
	params := getenv("DB_PARAMS", "charset=utf8mb4&parseTime=true&loc=Local")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, pass, host, port, name, params)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func open(mode string) (*gorm.DB, error) {
	return gorm.Open(gormmysql.Open(dsnFromEnv(mode)), &gorm.Config{})
}

type Env struct{ E *string }

func OpenDatabase(env *Env) (*gorm.DB, error) {
	mode := "development"

	if env != nil && env.E != nil && *env.E != "" {
		mode = *env.E
	}

	return open(mode)
}
