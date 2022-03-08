package middlewares

import (
	"fmt"
	"log"
	"net/http"
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

func LogRequestResponse(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := newLogResponseWriter(w)

		next.ServeHTTP(lw, r)

		logger.Infof("%s %s %d %s", r.Method, r.RequestURI, lw.statusCode, http.StatusText(lw.statusCode))
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{w, http.StatusOK}
}

func (lw *logResponseWriter) WriteHeader(statusCode int) {
	lw.statusCode = statusCode
	lw.ResponseWriter.WriteHeader(statusCode)
}
