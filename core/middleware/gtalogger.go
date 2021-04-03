package middleware

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/withmandala/go-log"
)

type GtaLogger struct {
	id            int
	fl            *log.Logger
	cl            *log.Logger
	enableFile    bool
	enableConsole bool
	f             *os.File
}

func NewLogger(id int, file, console bool) *GtaLogger {

	l := &GtaLogger{
		id:            id,
		enableFile:    file,
		enableConsole: console,
	}

	if file {

		filepath := filepath.Join("logs", fmt.Sprintf("app-%04d.log", id))
		l.f, _ = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		l.fl = log.New(l.f).WithDebug()
	}

	if console {
		l.cl = log.New(os.Stdout).WithColor().WithDebug()
	}

	return l
}

func (l *GtaLogger) LogTrace(f string, v ...interface{}) {
	if l.enableConsole {
		l.cl.Tracef(f, v...)
	}
	if l.enableFile {
		l.fl.Tracef(f, v...)
	}
}

func (l *GtaLogger) LogInfo(f string, v ...interface{}) {
	if l.enableConsole {
		l.cl.Infof(f, v...)
	}
	if l.enableFile {
		l.fl.Infof(f, v...)
	}
}

func (l *GtaLogger) LogWarn(f string, v ...interface{}) {
	if l.enableConsole {
		l.cl.Warnf(f, v...)
	}
	if l.enableFile {
		l.fl.Warnf(f, v...)
	}
}

func (l *GtaLogger) LogError(f string, v ...interface{}) {
	if l.enableConsole {
		l.cl.Errorf(f, v...)
	}
	if l.enableFile {
		l.fl.Errorf(f, v...)
	}
}

func (l *GtaLogger) LogFatal(f string, v ...interface{}) {
	if l.enableConsole {
		l.cl.Fatalf(f, v...)
	}
	if l.enableFile {
		l.fl.Fatalf(f, v...)
	}
}

func (l *GtaLogger) Close() {
	if l.enableFile && l.f != nil {
		l.f.Close()
	}
}
