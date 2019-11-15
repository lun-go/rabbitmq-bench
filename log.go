package main

import (
	"github.com/sirupsen/logrus"
)

var (
	logmq = new(mylog)
)

// mylog level
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelEnd
)

type Fields map[string]interface{}

type mylog struct {
	level int
}

func (l *mylog) SetLevel(lvl string) {
	switch lvl {
	case "debug":
		l.level = LevelDebug
	case "info":
		l.level = LevelInfo
	case "warn":
		l.level = LevelWarn
	case "error":
		l.level = LevelEnd
	default:
		l.level = LevelEnd
	}
	logrus.SetLevel(logrus.DebugLevel)
}

func (l *mylog) Debug(msg string, fields Fields) {
	if LevelDebug >= l.level {
		logrus.WithFields(logrus.Fields(fields)).Debug(msg)
	}
}

func (l *mylog) Info(msg string, fields Fields) {
	if LevelInfo >= l.level {
		logrus.WithFields(logrus.Fields(fields)).Info(msg)
	}
}

func (l *mylog) Warn(msg string, fields Fields) {
	if LevelWarn >= l.level {
		logrus.WithFields(logrus.Fields(fields)).Warn(msg)
	}
}

func (l *mylog) End(msg string, fields Fields) {
	if LevelEnd >= l.level {
		logrus.WithFields(logrus.Fields(fields)).Error(msg)
	}
}

func (l *mylog) Printf(level int, format string, v ...interface{}) {
	if level >= l.level {
		logrus.Printf(format, v...)
	}
}

func (l *mylog) Fields(level int, msg string, fields Fields) {
	if level >= l.level {
		logrus.WithFields(logrus.Fields(fields)).Info(msg)
	}
}
