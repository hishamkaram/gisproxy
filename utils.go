package gisproxy

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

//GetLogger return logger
func GetLogger() (logger *logrus.Logger) {
	logger = logrus.New()
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	logger.Formatter = Formatter
	return
}

//ReadFile read file
func ReadFile(path string) (data []byte, err error) {
	logger := GetLogger()
	data, err = ioutil.ReadFile(path)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
