package server

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	log *log.Logger
}

type iLogger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
}

func NewLogger() (l *Logger) {
	l = &Logger{log: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)}
	return
}

func (l *Logger) Debug(args ...interface{}) {
	l.log.SetPrefix("DEBUG ")
	l.log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	l.log.Output(2, fmt.Sprint(args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.SetPrefix("DEBUG ")
	l.log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	l.log.Output(2, fmt.Sprintf(format, args...))
}

func (l *Logger) Info(args ...interface{}) {
	l.log.SetPrefix("INFO ")
	l.log.SetFlags(log.LstdFlags)
	l.log.Println(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.SetPrefix("INFO ")
	l.log.SetFlags(log.LstdFlags)
	l.log.Printf(format, args...)
}
