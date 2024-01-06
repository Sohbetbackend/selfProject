package utils

import "github.com/sirupsen/logrus"

var Logger *logrus.Entry

func LoggerDesc(desc string) *logrus.Entry {
	Logger = Logger.WithField("desc", desc)
	return Logger
}

func GetLoggerDesc() string {
	desc := ""
	if v, ok := Logger.Data["desc"].(string); ok {
		desc = v
	}
	return desc
}
