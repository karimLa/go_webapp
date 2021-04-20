package lib

import (
	"log"
	"os"
)

func InitLog() *log.Logger {
	return log.New(os.Stdout, "webapp ", log.LstdFlags)
}
