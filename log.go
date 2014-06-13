package kamino

import "github.com/Sirupsen/logrus"

var logger = logrus.New()

func init() {
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.Info // default to Info
}

/*
SetLogLevel sets the log level (from Sirupsen/logrus) on kamino's internal
logger
*/
func SetLogLevel(level logrus.Level) {
	logger.Level = level
}
