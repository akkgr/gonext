package logger

import (
	"log"
	"os"
)

// New returns a configured logger
func New() *log.Logger {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	return logger
}
