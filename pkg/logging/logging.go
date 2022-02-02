// Package logging for work with logrus logger
package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// hook struct
type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

// Fire - func for writing logs
func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

// GetLogger - func for getting logger to another packages
func GetLogger() *Logger {
	return &Logger{e}
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}

// init - func for init logger
func init() {
	// create logrus logger
	l := logrus.New()
	l.SetReportCaller(true)
	// set log`s message format
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", fileName, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	// creating dir for save log file
	err := os.MkdirAll("../../logs", 0644)
	if err != nil {
		log.Fatal(err)
	}

	// check log file and create if not exist
	allFile, err := os.OpenFile("../../logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		log.Fatal(err)
	}

	l.SetOutput(io.Discard)

	l.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
