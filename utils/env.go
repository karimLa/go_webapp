package utils

import (
	"fmt"

	"github.com/nicholasjackson/env"
)

var (
	bindAddress = env.String("BIND_ADDRESS", false, "", "Bind address for the server")
	bindPort    = env.String("BIND_PORT", false, "3000", "Bind port for the server")
	dbHost      = env.String("DB_HOST", false, "localhost", "database host, i.e. localhost")
	dbPort      = env.String("DB_PORT", false, "5432", "database port, i.e. 5432")
	dbName      = env.String("DB_NAME", false, "dev_db", "database name")
	dbUser      = env.String("DB_USER", false, "sora", "database user")
	dbPassword  = env.String("DB_PASSWORD", false, "sora_password", "database user password")
)

func GetBindAdress() string {
	return fmt.Sprintf("%s:%s", *bindAddress, *bindPort)
}

func GetDB() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", *dbHost, *dbUser, *dbPassword, *dbName, *dbPort)
}
