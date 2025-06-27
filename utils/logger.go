package utils

import (
	"log"
)

type LogLevel int

const (
	INFO LogLevel = iota
	DEBUG
)

var logger = log.New(log.Writer(), "logger: ", log.Ldate|log.Ltime|log.Lshortfile)

func Log(level LogLevel, message string) {
	switch level {
	case INFO:
		logger.Println(message)
	case DEBUG:
		logger.Println(message)
	}
}

func Info(message string) {
	Log(INFO, message)
}

func Debug(message string) {
	Log(DEBUG, message)
}