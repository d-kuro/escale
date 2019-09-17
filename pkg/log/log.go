package log

import (
	"log"
	"os"
)

var Logger logger

type logger interface {
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type stdLogger struct {
	stderr *log.Logger
	stdout *log.Logger
}

func NewStdLogger() {
	Logger = &stdLogger{
		stdout: log.New(os.Stdout, "", log.LstdFlags),
		stderr: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (l *stdLogger) Printf(format string, args ...interface{}) {
	l.stdout.Printf(format, args...)
}

func (l *stdLogger) Errorf(format string, args ...interface{}) {
	l.stderr.Printf(format, args...)
}

func (l *stdLogger) Fatalf(format string, args ...interface{}) {
	l.stderr.Fatalf(format, args...)
}
