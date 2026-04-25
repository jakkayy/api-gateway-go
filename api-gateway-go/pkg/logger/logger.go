package logger

import "log"

type Logger struct{}

func New() *Logger {
	return &Logger{}
}

func (i *Logger) Info(v ...interface{}) {
	log.Println(v...)
}
