package lib

import (
	"log"
	"os"
)

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "webapp ", log.LstdFlags)
}
