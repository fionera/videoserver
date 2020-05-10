package logging

import "github.com/sirupsen/logrus"

type Logger struct {
	Module string
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args)
}

func (l *Logger) Fatal(err error) {
	logrus.Fatal(err)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	logrus.Printf(format, args)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	logrus.Infof(format, args)
}